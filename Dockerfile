# intended to be built via goreleaser
FROM scratch
COPY --from=alpine:3.17.3 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY discord-notify-ip-change .
ENTRYPOINT ["/discord-notify-ip-change"]
