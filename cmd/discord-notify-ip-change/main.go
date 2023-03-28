package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dustinspecker/discord-notify-ip-change/internal/ip"
)

func main() {
	var ipURL string
	flag.StringVar(&ipURL, "ip-url", "", `URL to retrieve public IP in format of {"ip": "0.0.0.0"}`)
	flag.Parse()

	publicIp, err := ip.Get(ipURL)
	if err != nil {
		log.Fatalf("error getting public IP: %v", err)
	}

	fmt.Println(publicIp)
}
