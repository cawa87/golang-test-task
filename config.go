package main

import "time"

const (
	logPlace = "place"
)

type Config struct {
	LogLevel string        `yaml:"logLevel"`
	Server   ServerConfig  `yaml:"server"`
	Fetcher  FetcherConfig `yaml:"fetcher"`
	Parser   ParserConfig  `yaml:"parser"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type FetcherConfig struct {
	Timeout     time.Duration `yaml:"timeout"`
	WorkerCount int           `yaml:"workerCount"`
}

type ParserConfig struct {
	WorkerCount int `yaml:"workerCount"`
}
