.PHONY: all build

tests:
	go test ./... -count 1 -v

all: tests