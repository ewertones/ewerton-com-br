.PHONY: dev build docker clean

# Run with live output; kill with Ctrl+C
dev:
	go run .

# Produce a local binary
build:
	go build -ldflags="-s -w" -o ewerton-go .

# Build and run Docker image
docker:
	docker build -t ewerton-go .
	docker run --rm -p 3000:3000 ewerton-go

clean:
	rm -f ewerton-go
