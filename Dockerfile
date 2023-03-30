# builder
FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum .

RUN go mod download

COPY . .

ENV CGO_ENABLED=0

RUN go build -o discord-notify-ip-change ./cmd/discord-notify-ip-change/main.go

# final
FROM scratch

COPY --from=builder /app/discord-notify-ip-change .

ENTRYPOINT ["/discord-notify-ip-change"]
