package main

import (
	"flag"
	"log"
	"time"

	"github.com/dustinspecker/discord-notify-ip-change/internal"
)

type messageData struct {
	PublicIP string
}

func main() {
	var discordWebhookURL string
	flag.StringVar(&discordWebhookURL, "discord-webhook-url", "", "Discord Webhook URL to send message to")

	var format string
	flag.StringVar(&format, "format", `{"content": "{{ .PublicIP }}"}`, "template for rendering message to send to Discord")

	var interval string
	flag.StringVar(&interval, "interval", "4h", "time to wait between checking if IP has changed")

	var ipURL string
	flag.StringVar(&ipURL, "ip-url", "", `URL to retrieve public IP in format of {"ip": "0.0.0.0"}`)

	var timeout string
	flag.StringVar(&timeout, "timeout", "60s", "amount of time to wait for response from -ip-url")

	flag.Parse()

	parsedInterval, err := time.ParseDuration(interval)
	if err != nil {
		log.Fatalf("unable to parse interval: %v", err)
	}

	parsedTimeout, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatalf("unable to parse timeout: %v", err)
	}

	var lastPublicIP string

	for {
		publicIp, err := internal.GetIP(ipURL, parsedTimeout)
		if err != nil {
			log.Printf("error getting public IP: %v", err)
		}

		if publicIp != lastPublicIP {
			lastPublicIP = publicIp

			renderedMessageStr, err := internal.RenderMessage(format, messageData{PublicIP: publicIp})
			if err != nil {
				log.Printf("error rendering message: %v", err)
			}

			if err := internal.SendMessage(discordWebhookURL, renderedMessageStr); err != nil {
				log.Printf("error sending message to discord: %v", err)
			}
		}

		time.Sleep(parsedInterval)
	}
}
