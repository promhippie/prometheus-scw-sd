package main

import (
	"errors"

	"github.com/go-kit/kit/log/level"
	"github.com/promhippie/prometheus-scw-sd/pkg/action"
	"github.com/promhippie/prometheus-scw-sd/pkg/config"
	"gopkg.in/urfave/cli.v2"
)

var (
	// ErrMissingOutputFile defines the error if output.file is empty.
	ErrMissingOutputFile = errors.New("Missing path for output.file")

	// ErrMissingScwToken defines the error if scw.token is empty.
	ErrMissingScwToken = errors.New("Missing required scw.token")

	// ErrMissingScwOrg defines the error if scw.org is empty.
	ErrMissingScwOrg = errors.New("Missing required scw.org")

	// ErrMissingAnyCredentials defines the error if no credentials are provided.
	ErrMissingAnyCredentials = errors.New("Missing any credentials")
)

// Server provides the sub-command to start the server.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start integrated server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "web.address",
				Value:       "0.0.0.0:9000",
				Usage:       "Address to bind the metrics server",
				EnvVars:     []string{"PROMETHEUS_SCW_WEB_ADDRESS"},
				Destination: &cfg.Server.Addr,
			},
			&cli.StringFlag{
				Name:        "web.path",
				Value:       "/metrics",
				Usage:       "Path to bind the metrics server",
				EnvVars:     []string{"PROMETHEUS_SCW_WEB_PATH"},
				Destination: &cfg.Server.Path,
			},
			&cli.StringFlag{
				Name:        "output.file",
				Value:       "/etc/prometheus/scw.json",
				Usage:       "Path to write the file_sd config",
				EnvVars:     []string{"PROMETHEUS_SCW_OUTPUT_FILE"},
				Destination: &cfg.Target.File,
			},
			&cli.IntFlag{
				Name:        "output.refresh",
				Value:       30,
				Usage:       "Discovery refresh interval in seconds",
				EnvVars:     []string{"PROMETHEUS_SCW_OUTPUT_REFRESH"},
				Destination: &cfg.Target.Refresh,
			},
			&cli.StringFlag{
				Name:    "scw.token",
				Value:   "",
				Usage:   "Access token for the Scaleway API",
				EnvVars: []string{"PROMETHEUS_SCW_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "scw.org",
				Value:   "",
				Usage:   "Organization for the Scaleway API",
				EnvVars: []string{"PROMETHEUS_SCW_ORG"},
			},
			&cli.StringFlag{
				Name:    "scw.region",
				Value:   "",
				Usage:   "Region for the Scaleway API",
				EnvVars: []string{"PROMETHEUS_SCW_REGION"},
			},
			&cli.StringFlag{
				Name:    "scw.config",
				Value:   "",
				Usage:   "Path to Scaleway configuration file",
				EnvVars: []string{"PROMETHEUS_SCW_CONFIG"},
			},
		},
		Action: func(c *cli.Context) error {
			logger := setupLogger(cfg)

			if c.IsSet("scw.config") {
				if err := readConfig(c.String("scw.config"), cfg); err != nil {
					level.Error(logger).Log(
						"msg", "Failed to read config",
						"err", err,
					)

					return err
				}
			}

			if cfg.Target.File == "" {
				level.Error(logger).Log(
					"msg", ErrMissingOutputFile,
				)

				return ErrMissingOutputFile
			}

			if c.IsSet("scw.token") && c.IsSet("scw.org") && c.IsSet("scw.region") {
				credentials := config.Credential{
					Project: "default",
					Token:   c.String("scw.token"),
					Org:     c.String("scw.org"),
					Region:  c.String("scw.region"),
				}

				cfg.Target.Credentials = append(
					cfg.Target.Credentials,
					credentials,
				)

				if credentials.Token == "" {
					level.Error(logger).Log(
						"msg", ErrMissingScwToken,
					)

					return ErrMissingScwToken
				}

				if credentials.Org == "" {
					level.Error(logger).Log(
						"msg", ErrMissingScwOrg,
					)

					return ErrMissingScwOrg
				}
			}

			if len(cfg.Target.Credentials) == 0 {
				level.Error(logger).Log(
					"msg", ErrMissingAnyCredentials,
				)

				return ErrMissingAnyCredentials
			}

			return action.Server(cfg, logger)
		},
	}
}
