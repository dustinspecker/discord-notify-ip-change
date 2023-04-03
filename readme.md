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

## Test

```bash
make test
go install github.com/onsi/ginkgo/v2/ginkgo
make int-test
```

## License

MIT
