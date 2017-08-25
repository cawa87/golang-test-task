package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	flagCfg := flag.String("config", "config.yaml", "Path to YAML configuration file")
	flag.Parse()

	logger := logrus.StandardLogger()

	cfg, err := readConfig(*flagCfg)
	if err != nil {
		logger.WithError(err).Fatal("Failed to parse configuration file")
	}

	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logger.WithError(err).Fatalf("Unknown log level %s", cfg.LogLevel)
	}
	logger.Level = logLevel

	fetcher, err := NewDocumentFetcher(cfg.Fetcher, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to init document fetcher")
	}

	parser, err := NewDocumentParser(cfg.Parser, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to init document parser")
	}

	server, err := NewServer(cfg.Server, fetcher, parser, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to init HTTP server")
	}

	go func() {
		if err = server.Start(); err != nil {
			logger.WithError(err).Fatal("Failed to start HTTP server")
		}
	}()

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, os.Interrupt, os.Kill)

	<-stopCh

	if err = server.Shutdown(); err != nil {
		logger.WithError(err).Error("Failed to shutdwon HTTP server")
	}
}

func readConfig(path string) (Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
