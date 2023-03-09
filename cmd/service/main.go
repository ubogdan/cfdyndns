package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/ubogdan/cfdyndns"
	"sigs.k8s.io/yaml"
)

type Config struct {
	PublicIPResolver string `yaml:"PublicIPResolver"`
	Cloudflare       struct {
		Token  string `yaml:"token"`
		Zone   string `yaml:"zone"`
		Record string `yaml:"record"`
	} `yaml:"Cloudflare"`
}

var config string

func main() {
	flag.StringVar(&config, "c", "config.yml", "Path to config file")
	flag.Parse()

	f, err := os.Open(config)
	if err != nil {
		log.Fatalf("Failed to open config file: %s", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Failed to unmarshal config file: %s", err)
	}

	if len(cfg.Cloudflare.Token) == 0 {
		log.Fatalf("Cloudflare token is missing")
	}

	if len(cfg.Cloudflare.Zone) == 0 {
		log.Fatalf("Cloudflare zone is missing")

	}
	if len(cfg.Cloudflare.Record) == 0 {
		log.Fatalf("Cloudflare record is missing")
	}

	if len(cfg.PublicIPResolver) == 0 {
		cfg.PublicIPResolver = "https://ifconfig.io/ip"
	}

	cli := cfdyndns.New(cfg.Cloudflare.Token)

retryGetZone:
	zone, err := cli.GetZone(cfg.Cloudflare.Zone)
	if err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			var netOpErr *net.OpError
			if errors.As(urlErr.Err, &netOpErr) {
				log.Printf("Internet connection problem %s", netOpErr)
				time.Sleep(3 * time.Second)
				goto retryGetZone
			}
		}

		log.Printf("Failed to retrieve zone %s", err)
		return
	}

retryGetRecord:
	record, err := cli.GetRecord(zone.Id, cfg.Cloudflare.Record)
	if err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			var netOpErr *net.OpError
			if errors.As(urlErr.Err, &netOpErr) {
				log.Printf("Internet connection problem %s", netOpErr)
				time.Sleep(3 * time.Second)
				goto retryGetRecord
			}
		}

		log.Printf("Failed to retrieve record %s", err)
		return
	}

	publicIP := net.ParseIP(record.Content)

	checkFunc := func() {
		remoteIP, err := cfdyndns.GetOutboundIP(cfg.PublicIPResolver, 1*time.Second)
		if err != nil {
			log.Printf("Error %s", err)
			return
		}
		if remoteIP.Equal(publicIP) {
			log.Printf("Record %s is pointing to %s (up to date)", record.Name, remoteIP.String())
			return
		}

		update := *record
		update.Content = remoteIP.String()

		log.Printf("Updating record %s to %s", cfg.Cloudflare.Record, remoteIP.String())

		err = cli.UpdateRecord(update)
		if err != nil {
			log.Printf("failed to update record %s", err)

			return
		}

		publicIP = remoteIP
	}

	checkFunc()

	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			checkFunc()
		}
	}
}
