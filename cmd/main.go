package main

import (
	"os"
	"strconv"

	"github.com/edubarbieri/rinha-2024-q1/internal/application"
	"github.com/edubarbieri/rinha-2024-q1/internal/application/web"
)

func main() {
	createTxUseCase := application.NewCreateTransactionUseCase()
	getStatementUseCase := application.NewGetStatementUseCase()
	app := web.NewWebApplication(createTxUseCase, getStatementUseCase)

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
	err = app.Server(port)
	if err != nil {
		panic(err)
	}
}
