.PHONY: build
build:
	@echo "build client ./dist/ktpc"
	@go build -o dist/ktpc cmd/ktpc/*
	@
	@echo "build server ./dist/ktpd"
	@go build -o dist/ktpd cmd/ktpd/*

.PHONY: deploy-fly
deploy-fly:
	fly deploy

.PHONY: deploy
deploy: build deploy-fly