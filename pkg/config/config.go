package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

// Credential defines a single project credential.
type Credential struct {
	Project   string `json:"project" yaml:"project"`
	AccessKey string `json:"access_key" yaml:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	Org       string `json:"org" yaml:"org"`
	Zone      string `json:"zone" yaml:"zone"`
}

// Server defines the general server configuration.
type Server struct {
	Addr string `json:"addr" yaml:"addr"`
	Path string `json:"path" yaml:"path"`
	Web  string `json:"web_config" yaml:"web_config"`
}

// Logs defines the level and color for log configuration.
type Logs struct {
	Level  string `json:"level" yaml:"level"`
	Pretty bool   `json:"pretty" yaml:"pretty"`
}

// Zones defines the available zones per api endpoint.
type Zones struct {
	Instance  []string `json:"instance" yaml:"instance"`
	Baremetal []string `json:"baremetal" yaml:"baremetal"`
}

// Target defines the target specific configuration.
type Target struct {
	Engine         string       `json:"engine" yaml:"engine"`
	File           string       `json:"file" yaml:"file"`
	Refresh        int64        `json:"refresh" yaml:"refresh"`
	CheckInstance  bool         `json:"check_instance" yaml:"check_instance"`
	CheckBaremetal bool         `json:"check_baremetal" yaml:"check_baremetal"`
	Credentials    []Credential `json:"credentials" yaml:"credentials"`
}

// Config is a combination of all available configurations.
type Config struct {
	Server Server `json:"server" yaml:"server"`
	Logs   Logs   `json:"logs" yaml:"logs"`
	Zones  Zones  `json:"zones" yaml:"zones"`
	Target Target `json:"target" yaml:"target"`
}

// Load initializes a default configuration struct.
func Load() *Config {
	return &Config{
		Target: Target{
			Credentials: make([]Credential, 0),
		},
	}
}

// Value returns the config value based on a DSN.
func Value(val string) (string, error) {
	if strings.HasPrefix(val, "file://") {
		content, err := os.ReadFile(
			strings.TrimPrefix(val, "file://"),
		)

		if err != nil {
			return "", fmt.Errorf("failed to parse secret file: %w", err)
		}

		return string(content), nil
	}

	if strings.HasPrefix(val, "base64://") {
		content, err := base64.StdEncoding.DecodeString(
			strings.TrimPrefix(val, "base64://"),
		)

		if err != nil {
			return "", fmt.Errorf("failed to parse base64 value: %w", err)
		}

		return string(content), nil
	}

	return val, nil
}
