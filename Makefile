build:
	go build -o ./.bin/ ./cmd/...

test:
	go test ./internal/...

int-test:
	ginkgo run ./integration/
