package main

import (
	"bytes"
	"flag"
	"log"
	"time"

	"github.com/dustinspecker/discord-notify-ip-change/internal/discord"
	"github.com/dustinspecker/discord-notify-ip-change/internal/ip"
)

func main() {
	var discordWebhookURL string
	flag.StringVar(&discordWebhookURL, "discord-webhook-url", "", "Discord Webhook URL to send message to")

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

	if err := discord.SendMessage(discordWebhookURL, bytes.NewReader([]byte(publicIp))); err != nil {
		log.Fatalf("error sending message to discord: %v", err)
	}
}
