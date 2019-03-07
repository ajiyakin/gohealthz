package main

import (
	"flag"
	"time"
)

type config struct {
	updaterInterval   time.Duration
	httpClientTimeout time.Duration
}

func parseFlag() (*config, error) {
	updaterIntervalFlag := flag.String("interval", "5m", "Updater interval (default: 5m)")
	httpClientTimeoutFlag := flag.String("timeout", "800ms", "Updater interval (default: 800ms)")
	flag.Parse()

	updaterInterval, err := time.ParseDuration(*updaterIntervalFlag)
	if err != nil {
		return nil, err
	}

	httpClientTimeout, err := time.ParseDuration(*httpClientTimeoutFlag)
	if err != nil {
		return nil, err
	}

	c := config{
		updaterInterval:   updaterInterval,
		httpClientTimeout: httpClientTimeout,
	}
	return &c, nil
}
