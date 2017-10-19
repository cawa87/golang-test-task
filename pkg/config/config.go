package config

import "github.com/l-vitaly/goenv"

// env name constants
const (
	MaxWorkersEnvName = "CRAWL_MAX_WORKERS"
	BindAddrEnvName   = "CRAWL_BIND_ADDR"
)

// Config service configuration
type Config struct {
	MaxWorkers int
	BindAddr   string
}

// Parse parse env config vars
func Parse() *Config {
	cfg := &Config{}
	goenv.IntVar(&cfg.MaxWorkers, MaxWorkersEnvName, 100)
	goenv.StringVar(&cfg.BindAddr, BindAddrEnvName, ":9000")
	goenv.Parse()
	return cfg
}

// Validate validate config
func Validate(cfg *Config) error {
	return nil
}
