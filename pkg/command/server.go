package command

import (
	"errors"

	"github.com/go-kit/kit/log/level"
	"github.com/promhippie/prometheus-scw-sd/pkg/action"
	"github.com/promhippie/prometheus-scw-sd/pkg/config"
	"github.com/urfave/cli/v2"
)

// Server provides the sub-command to start the server.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start integrated server",
		Flags: ServerFlags(cfg),
		Before: func(c *cli.Context) error {
			cfg.Zones.Instance = c.StringSlice("scw.instance_zone")
			cfg.Zones.Baremetal = c.StringSlice("scw.baremetal_zone")

			return nil
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
					"msg", "Missing path for output.file",
				)

				return errors.New("missing path for output.file")
			}

			if c.IsSet("scw.access_key") && c.IsSet("scw.secret_key") {
				credentials := config.Credential{
					Project:   "default",
					AccessKey: c.String("scw.access_key"),
					SecretKey: c.String("scw.secret_key"),
					Org:       c.String("scw.org"),
					Zone:      c.String("scw.zone"),
				}

				cfg.Target.Credentials = append(
					cfg.Target.Credentials,
					credentials,
				)

				if credentials.AccessKey == "" {
					level.Error(logger).Log(
						"msg", "Missing required scw.access_key",
					)

					return errors.New("missing required scw.access_key")
				}

				if credentials.SecretKey == "" {
					level.Error(logger).Log(
						"msg", "Missing required scw.secret_key",
					)

					return errors.New("missing required scw.secret_key")
				}
			}

			if len(cfg.Target.Credentials) == 0 {
				level.Error(logger).Log(
					"msg", "Missing any credentials",
				)

				return errors.New("missing any credentials")
			}

			return action.Server(cfg, logger)
		},
	}
}

// ServerFlags defines the available server flags.
func ServerFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
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
			Name:        "web.config",
			Value:       "",
			Usage:       "Path to web-config file",
			EnvVars:     []string{"PROMETHEUS_SCW_WEB_CONFIG"},
			Destination: &cfg.Server.Web,
		},
		&cli.StringFlag{
			Name:        "output.engine",
			Value:       "file",
			Usage:       "Enabled engine like file or http",
			EnvVars:     []string{"PROMETHEUS_SCW_OUTPUT_ENGINE"},
			Destination: &cfg.Target.Engine,
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
		&cli.BoolFlag{
			Name:        "scw.check_instance",
			Value:       true,
			Usage:       "Enable instance gathering",
			EnvVars:     []string{"PROMETHEUS_SCW_CHECK_INSTANCE"},
			Destination: &cfg.Target.CheckInstance,
		},
		&cli.BoolFlag{
			Name:        "scw.check_baremetal",
			Value:       true,
			Usage:       "Enable baremetal gathering",
			EnvVars:     []string{"PROMETHEUS_SCW_CHECK_BAREMETAL"},
			Destination: &cfg.Target.CheckBaremetal,
		},
		&cli.StringFlag{
			Name:    "scw.access_key",
			Value:   "",
			Usage:   "Access key for the Scaleway API",
			EnvVars: []string{"PROMETHEUS_SCW_ACCESS_KEY"},
		},
		&cli.StringFlag{
			Name:    "scw.secret_key",
			Value:   "",
			Usage:   "Secret key for the Scaleway API",
			EnvVars: []string{"PROMETHEUS_SCW_SECRET_KEY"},
		},
		&cli.StringFlag{
			Name:    "scw.org",
			Value:   "",
			Usage:   "Organization for the Scaleway API",
			EnvVars: []string{"PROMETHEUS_SCW_ORG"},
		},
		&cli.StringFlag{
			Name:    "scw.zone",
			Value:   "",
			Usage:   "Zone for the Scaleway API",
			EnvVars: []string{"PROMETHEUS_SCW_ZONE"},
		},
		&cli.StringFlag{
			Name:    "scw.config",
			Value:   "",
			Usage:   "Path to Scaleway configuration file",
			EnvVars: []string{"PROMETHEUS_SCW_CONFIG"},
		},
		&cli.StringSliceFlag{
			Name:    "scw.instance_zone",
			Value:   cli.NewStringSlice("fr-par-1", "nl-ams-1"),
			Usage:   "List of available zones for instance API",
			EnvVars: []string{"PROMETHEUS_SCW_INSTANCE_ZONES"},
			Hidden:  true,
		},
		&cli.StringSliceFlag{
			Name:    "scw.baremetal_zone",
			Value:   cli.NewStringSlice("fr-par-2"),
			Usage:   "List of available zones for baremetal API",
			EnvVars: []string{"PROMETHEUS_SCW_BAREMETAL_ZONES"},
			Hidden:  true,
		},
	}
}
