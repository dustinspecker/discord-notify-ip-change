# intended to be built via goreleaser
FROM scratch
COPY discord-notify-ip-change .
ENTRYPOINT ["/discord-notify-ip-change"]
