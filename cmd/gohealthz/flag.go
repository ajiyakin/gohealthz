package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type config struct {
	updaterInterval   time.Duration
	httpClientTimeout time.Duration
}

func (c config) String() string {
	return fmt.Sprintf("update_interval=%s http_client_timeout=%s", c.updaterInterval.String(), c.httpClientTimeout.String())
}

func parseFlag() (*config, error) {
	updaterIntervalFlag := flag.String("interval", "5m", "Updater interval")
	httpClientTimeoutFlag := flag.String("timeout", "800ms", "Updater interval")
	helpFlag := flag.Bool("help", false, "print this message")
	flag.Parse()
	if *helpFlag || len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

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
