build:
	CGO_ENABLED=0 GOOS=linux go build -o bin ./...

run:
	go run cmd/meisterwerk/main.go

dock:
	CGO_ENABLED=0 GOOS=linux go build -o bin ./...
	docker build . -t meisterwerk-test
	docker run -p 8000:8000 meisterwerk-test:latest

compose:
	CGO_ENABLED=0 GOOS=linux go build -o bin ./...
	docker-compose build --no-cache
	docker-compose up --remove-orphans

mock:
	GO111MODULE=on go get github.com/golang/mock/mockgen@v1.4.4
	mockgen -source=internal/gh/gh.go -destination=internal/mock/mock_gh.go -package=mock
	mockgen -source=internal/store/store.go -destination=internal/mock/mock_store.go -package=mock

test:
	go test  -cover ./...