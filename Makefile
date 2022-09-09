APP_NAME=go-rest-api-boilerplate
IMAGE_TAG=$(shell git rev-parse --short HEAD)

.PHONY:

build:
	go build -o $(APP_NAME) cmd/main.go

test:
	go test -v -cover -covermode=atomic ./...

migrate_up:
	[ -f ./$(APP_NAME) ] && ./$(APP_NAME) migrate up || go build -o $(APP_NAME) cmd/main.go;\
	DB_PASSWORD=root_password ./$(APP_NAME) migrate up

migrate_down:
	[ -f ./$(APP_NAME) ] && ./$(APP_NAME) migrate down || go build -o $(APP_NAME) cmd/main.go;\
	DB_PASSWORD=root_password ./$(APP_NAME) migrate down