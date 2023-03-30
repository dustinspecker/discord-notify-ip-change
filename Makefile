build:
	go build -o ./.bin/ ./cmd/...

build-docker:
	docker build . --tag discord-notify-ip-change

test:
	go test ./internal/...

int-test:
	ginkgo run -p ./integration/
