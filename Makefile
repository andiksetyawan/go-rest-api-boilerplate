APP_NAME=go-rest-api-boilerplate
#IMAGE_REGISTRY=docker.io/andiksetyawan
#IMAGE_NAME=$(IMAGE_REGISTRY)/$(APP_NAME)
IMAGE_NAME=$(APP_NAME)
IMAGE_TAG=$(shell git rev-parse --short HEAD)

.PHONY:

build:
	go build -o $(APP_NAME) cmd/main.go

test:
	go test -v -cover -covermode=atomic ./...

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_NAME):latest

docker-push:
	docker push $(IMAGE_NAME):latest

migrate-up:
	[ -f ./$(APP_NAME) ] && ./$(APP_NAME) migrate up || go build -o $(APP_NAME) cmd/main.go;\
	./$(APP_NAME) migrate up

migrate-down:
	[ -f ./$(APP_NAME) ] && ./$(APP_NAME) migrate down || go build -o $(APP_NAME) cmd/main.go;\
	./$(APP_NAME) migrate down

docker-compose-run:
	docker-compose up --build -d

docker-compose-stop:
	docker-compose down