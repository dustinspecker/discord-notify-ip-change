# discord-notify-ip-change

> Send a message to a Discord channel notifying Public IP Address changed

## Build

1. `make build`

## Usage

```bash
./.bin/discord-notify-ip-change \
  -discord-webhook-url "https://discord.com/webhooks/webhooktoken" \
  -ip-url "https://api.ipify.org/?format=json"
```

### Arguments

#### Required

- `-discord-webhook-url` Discord Webhook URL to send message to
- `-ip-url` URL to retrieve

#### Optional

- `-format` template for rendering message to send to Discord (default: `{"content": "{{ .PublicIP }}"}`)
- `-interval` time to wait between checking if IP has changed (default: `4h`)
- `-timeout` amount of time to wait for response from -ip-url (default: `60s`)

## Test

```bash
make test
go install github.com/onsi/ginkgo/v2/ginkgo
make int-test
```

## License

MIT
