package config

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
	File           string       `json:"file" yaml:"file"`
	Refresh        int          `json:"refresh" yaml:"refresh"`
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
