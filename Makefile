-include .env
export

APP_NAME = aeza-promo-instances-watchdog

.PHONY: bin
build:
	go build \
	-o ./bin/$(APP_NAME) ./cmd/$(APP_NAME)


.PHONY: app
app: build
	./bin/$(APP_NAME)
