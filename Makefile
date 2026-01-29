.PHONY: build clean run

BINARY=uptime_monitor

build:
	go build -o $(BINARY) .

clean:
	rm -f $(BINARY)

run: build
	./$(BINARY)
