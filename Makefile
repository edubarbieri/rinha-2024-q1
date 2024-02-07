build:
	go build -o bin/main main.go

run:
	PORT=3000 go run cmd/main.go