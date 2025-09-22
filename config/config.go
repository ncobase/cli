package config

// Config holds application configuration
type Config struct {
	Database  DatabaseConfig  `yaml:"database"`
	Server    ServerConfig    `yaml:"server"`
	Extension ExtensionConfig `yaml:"extension"`
}

type DatabaseConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ExtensionConfig struct {
	Path string `yaml:"path"`
}

// LoadConfig loads configuration from file
func LoadConfig(path string) (*Config, error) {
	return &Config{
		Database: DatabaseConfig{
			Driver: "postgres",
			DSN:    "postgres://localhost/db",
		},
		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Extension: ExtensionConfig{
			Path: "./plugins",
		},
	}, nil
}
