build:
	go build -o bin/main main.go

run:
	MYSQL_ADDRESS=localhost:3306 \
  MYSQL_DATABASE=rinha \
  MYSQL_USER=rinha \
  MYSQL_PASSWORD=rinha \
	PORT=3000 go run cmd/main.go

docker-build:
	docker build . -t duduardo23/rinha:latest