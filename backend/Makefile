.PHONY: clean critic security lint test build run

APP_NAME = apiserver
BUILD_DIR = I:/whatsappscrapping-go/test-ecommerce/backend/build
MIGRATIONS_FOLDER =  I:/whatsappscrapping-go/test-ecommerce/backend/platform/migrations
DATABASE_URL = mysql://root:mysqlpw@tcp(192.168.90.10:3306)/testing

clean:
	rm -rf ./build

critic:
	gocritic check -enableAll ./...

security:
	gosec ./...

lint:
	golangci-lint run ./...

test: clean critic security lint
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

build: test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

run: swag build
	$(BUILD_DIR)/$(APP_NAME)

migrate.up:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" up

migrate.down:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" down

migrate.force:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" force $(version)

docker.run: docker.network docker.mysql swag docker.fiber docker.redis migrate.up

docker.network:
	docker network create -d bridge dev-network

docker.fiber.build:
	docker build -t fiber .

docker.fiber: docker.fiber.build
	docker run --rm -d \
		--name cgapp-fiber \
		--network dev-network \
		-p 5000:5000 \
		fiber

docker.mysql:
	docker run --rm -d \
		--name cgapp-mysql \
		--network dev-network \
		-e MYSQL_ROOT_PASSWORD=mysqlpw \
		-e MYSQL_DATABASE=testing \
		-p 3306:3306 \
		mysql

docker.redis:
	docker run --rm -d \
		--name cgapp-redis \
		--network dev-network \
		-p 6379:6379 \
		redis

docker.stop: docker.stop.fiber docker.stop.mysql docker.stop.redis

docker.stop.fiber:
	docker stop cgapp-fiber

docker.stop.mysql:
	docker stop cgapp-mysql

docker.stop.redis:
	docker stop cgapp-redis

swag:
	swag init