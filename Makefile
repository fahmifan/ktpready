.PHONY: build
build:
	@echo "build client ./dist/ktpc"
	@go build -o dist/client cmd/ktpc/*
	@
	@echo "build server ./dist/ktpd"
	@go build -o dist/server cmd/ktpd/*
