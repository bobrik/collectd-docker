package main

import (
	"flag"
	"log"
	"os"

	"github.com/bobrik/collectd-docker/collector"
	"github.com/fsouza/go-dockerclient"
)

func main() {
	endpoint := flag.String("endpoint", "unix:///var/run/docker.sock", "docker endpoint")
	host := flag.String("host", "", "host to report")
	interval := flag.Int("interval", 1, "interval to report")
	flag.Parse()

	if *host == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	client, err := docker.NewClient(*endpoint)
	if err != nil {
		log.Fatal(err)
	}

	writer := collector.NewCollectdWriter(*host, os.Stdout)

	collector := collector.NewCollector(client, writer, *interval)

	err = collector.Run(5)
	if err != nil {
		log.Fatal(err)
	}
}
