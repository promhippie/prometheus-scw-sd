package action

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/promhippie/prometheus-scw-sd/pkg/adapter"
	"github.com/promhippie/prometheus-scw-sd/pkg/config"
	"github.com/promhippie/prometheus-scw-sd/pkg/middleware"
	"github.com/promhippie/prometheus-scw-sd/pkg/version"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Server handles the server sub-command.
func Server(cfg *config.Config, logger *slog.Logger) error {
	logger.Info("Launching Prometheus Scaleway SD",
		"version", version.String,
		"revision", version.Revision,
		"date", version.Date,
		"go", version.Go,
		"engine", cfg.Target.Engine,
	)

	var gr run.Group

	{
		ctx := context.Background()
		clients := make(map[string]*scw.Client, len(cfg.Target.Credentials))

		for _, credential := range cfg.Target.Credentials {
			accessKey, err := config.Value(credential.AccessKey)

			if err != nil {
				logger.Error("Failed to read access key secret",
					"project", credential.Project,
					"err", err,
				)

				return fmt.Errorf("failed to read access key secret for %s", credential.Project)
			}

			secretKey, err := config.Value(credential.SecretKey)

			if err != nil {
				logger.Error("Failed to read secret key secret",
					"project", credential.Project,
					"err", err,
				)

				return fmt.Errorf("failed to read secret key secret for %s", credential.Project)
			}

			opts := make([]scw.ClientOption, 0)

			opts = append(opts, scw.WithUserAgent(
				fmt.Sprintf(
					"prometheus-scw-sd/%s (%s; %s; %s)",
					version.String,
					runtime.Version(),
					runtime.GOOS,
					runtime.GOARCH,
				),
			))

			opts = append(opts, scw.WithAuth(
				accessKey,
				secretKey,
			))

			if credential.Org != "" {
				opts = append(opts, scw.WithDefaultOrganizationID(
					credential.Org,
				))
			}

			if credential.Zone != "" {
				zone, err := scw.ParseZone(credential.Zone)

				if err != nil {
					logger.Error("Invalid zone provided",
						"project", credential.Project,
						"zone", credential.Zone,
					)

					return fmt.Errorf("invalid zone provided")
				}

				opts = append(opts, scw.WithDefaultZone(
					zone,
				))
			}

			client, err := scw.NewClient(opts...)

			if err != nil {
				logger.Error("Failed to initialize client",
					"project", credential.Project,
				)

				return fmt.Errorf("failed to initialize client")
			}

			clients[credential.Project] = client
		}

		disc := Discoverer{
			clients:        clients,
			logger:         logger,
			refresh:        cfg.Target.Refresh,
			checkInstance:  cfg.Target.CheckInstance,
			instanceZones:  cfg.Zones.Instance,
			checkBaremetal: cfg.Target.CheckBaremetal,
			baremetalZones: cfg.Zones.Baremetal,
			separator:      ",",
			lasts:          make(map[string]struct{}),
		}

		a := adapter.NewAdapter(ctx, cfg.Target.File, "scaleway-sd", disc, logger)
		a.Run()
	}

	{
		server := &http.Server{
			Addr:         cfg.Server.Addr,
			Handler:      handler(cfg, logger),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		gr.Add(func() error {
			logger.Info("Starting metrics server",
				"address", cfg.Server.Addr,
			)

			return web.ListenAndServe(
				server,
				&web.FlagConfig{
					WebListenAddresses: sliceP([]string{cfg.Server.Addr}),
					WebSystemdSocket:   boolP(false),
					WebConfigFile:      stringP(cfg.Server.Web),
				},
				logger,
			)
		}, func(reason error) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				logger.Error("Failed to shutdown metrics gracefully",
					"err", err,
				)

				return
			}

			logger.Info("Metrics shutdown gracefully",
				"reason", reason,
			)
		})
	}

	{
		stop := make(chan os.Signal, 1)

		gr.Add(func() error {
			signal.Notify(stop, os.Interrupt)

			<-stop

			return nil
		}, func(_ error) {
			close(stop)
		})
	}

	return gr.Run()
}

func handler(cfg *config.Config, logger *slog.Logger) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer(logger))
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Timeout)
	mux.Use(middleware.Cache)

	prom := promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{
			ErrorLog: promLogger{logger},
		},
	)

	mux.Route("/", func(root chi.Router) {
		root.Get(cfg.Server.Path, func(w http.ResponseWriter, r *http.Request) {
			prom.ServeHTTP(w, r)
		})

		root.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)

			io.WriteString(w, http.StatusText(http.StatusOK))
		})

		root.Get("/readyz", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)

			io.WriteString(w, http.StatusText(http.StatusOK))
		})

		if cfg.Target.Engine == "http" {
			root.Get("/sd", func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")

				content, err := os.ReadFile(cfg.Target.File)

				if err != nil {
					logger.Error("Failed to read service discovery data",
						"err", err,
					)

					http.Error(
						w,
						"Failed to read service discovery data",
						http.StatusInternalServerError,
					)

					return
				}

				w.WriteHeader(http.StatusOK)
				w.Write(content)
			})
		}
	})

	return mux
}

func boolP(i bool) *bool {
	return &i
}

func stringP(i string) *string {
	return &i
}

func sliceP(i []string) *[]string {
	return &i
}
