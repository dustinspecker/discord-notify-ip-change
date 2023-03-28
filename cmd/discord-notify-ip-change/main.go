package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/dustinspecker/discord-notify-ip-change/internal/ip"
)

func main() {
	var ipURL string
	flag.StringVar(&ipURL, "ip-url", "", `URL to retrieve public IP in format of {"ip": "0.0.0.0"}`)

	var timeout string
	flag.StringVar(&timeout, "timeout", "60s", "amount of time to wait for response from -ip-url")

	flag.Parse()

	parsedTimeout, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatalf("unable to parse timeout: %v", err)
	}

	publicIp, err := ip.Get(ipURL, parsedTimeout)
	if err != nil {
		log.Fatalf("error getting public IP: %v", err)
	}

	fmt.Println(publicIp)
}
