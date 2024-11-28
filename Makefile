.PHONY: run clean build

run:
	@mkdir -p data
	@if [ ! -f config/bots.yaml ]; then \
		echo "Error: config/bots.yaml file not found. Please create it first."; \
		exit 1; \
	fi
	go run main.go

build:
	@mkdir -p bin
	go build -o bin/nostr-bot main.go

clean:
	@echo "Cleaning database..."
	@rm -rf data/
	@rm -rf bin/
	@echo "Database cleaned successfully" 