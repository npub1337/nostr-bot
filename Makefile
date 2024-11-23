.PHONY: run clean

run:
	@mkdir -p data
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Please create it first."; \
		exit 1; \
	fi
	go run main.go

clean:
	@echo "Cleaning database..."
	@rm -rf data/
	@echo "Database cleaned successfully" 