package main

import (
	"flag"
	"log"
	"os"

	"path"

	"github.com/bobrik/collectd-docker/collector"
	"github.com/fsouza/go-dockerclient"
)

func main() {
	e := flag.String("endpoint", "unix:///var/run/docker.sock", "docker endpoint")
	c := flag.String("cert", "", "cert path for tls")
	h := flag.String("host", "", "host to report")
	i := flag.Int("interval", 1, "interval to report")
	flag.Parse()

	if *h == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var client *docker.Client
	var err error

	if *c != "" {
		client, err = docker.NewTLSClient(*e, path.Join(*c, "cert.pem"), path.Join(*c, "key.pem"), path.Join(*c, "ca.pem"))
	} else {
		client, err = docker.NewClient(*e)
	}

	if err != nil {
		log.Fatal(err)
	}

	writer := collector.NewCollectdWriter(*h, os.Stdout)

	collector := collector.NewCollector(client, writer, *i)

	err = collector.Run(5)
	if err != nil {
		log.Fatal(err)
	}
}
