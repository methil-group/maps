.PHONY: maps run build help

# Variables
DATA_DIR=data
MAPS_FILE=MAPS.txt
BINARY_NAME=server_bin

help:
	@echo "Usage:"
	@echo "  make maps   - Download .osm.pbf maps listed in MAPS.txt"
	@echo "  make run    - Run the Go server"
	@echo "  make build  - Build the Go server binary"

maps:
	@echo "⏳ Checking for maps to download..."
	@mkdir -p $(DATA_DIR)
	@while read -r url || [ -n "$$url" ]; do \
		case "$$url" in \
			\#*) continue ;; \
			"") continue ;; \
		esac; \
		filename=$$(basename "$$url"); \
		if [ ! -f "$(DATA_DIR)/$$filename" ]; then \
			echo "🌐 Downloading $$filename..."; \
			wget -q --show-progress -O "$(DATA_DIR)/$$filename" "$$url"; \
		else \
			echo "✅ $$filename already exists."; \
		fi; \
	done < $(MAPS_FILE)
	@echo "🎯 All maps are ready."

run:
	go run .

build:
	go build -o $(BINARY_NAME) .
