.PHONY: run build test tidy clean docker-up docker-down docker-logs

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

test:
	go test -v ./...

tidy:
	go mod tidy

clean:
	rm -rf bin/

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app
