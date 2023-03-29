# discord-notify-ip-change

> Send a message to a Discord channel notifying Public IP Address changed

## Build

1. `make build`

## Usage

```bash
./.bin/discord-notify-ip-change -ip-url "https://api.ipify.org/?format=json" -discord-webhook-url "https://discord.com/webhooks/webhooktoken"
```

## Test

```bash
make test
make int-test
```

## License

MIT
