.PHONY: build
build:
	@echo "build into ./dist/ktpready"
	@go build -o dist/ktpready cmd/client/*
