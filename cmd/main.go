package main

import (
	"os"
	"strconv"

	"github.com/edubarbieri/rinha-2024-q1/internal"
	"github.com/edubarbieri/rinha-2024-q1/internal/repository"
)

func main() {
	repo, err := repository.NewMysqlRepository(
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_DATABASE"),
		os.Getenv("MYSQL_ADDRESS"))
	if err != nil {
		panic(err)
	}

	app := internal.NewWebApplication(repo)

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
	err = app.Server(port)
	if err != nil {
		panic(err)
	}
}
