package config

// Server defines the general server configuration.
type Server struct {
	Addr string
	Path string
}

// Logs defines the level and color for log configuration.
type Logs struct {
	Level  string
	Pretty bool
}

// Target defines the target specific configuration.
type Target struct {
	File    string
	Refresh int
	Token   string
	Org     string
	Region  string
}

// Config is a combination of all available configurations.
type Config struct {
	Server Server
	Logs   Logs
	Target Target
}

// Load initializes a default configuration struct.
func Load() *Config {
	return &Config{}
}
