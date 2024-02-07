package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/edubarbieri/rinha-2024-q1/internal/application"
	"github.com/edubarbieri/rinha-2024-q1/internal/application/model"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type WebApplication struct {
	createTxUseCase     *application.CreateTransactionUseCase
	getStatementUseCase *application.GetStatementUseCase
}

func NewWebApplication(
	createTxUseCase *application.CreateTransactionUseCase,
	getStatementUseCase *application.GetStatementUseCase,
) *WebApplication {
	return &WebApplication{
		createTxUseCase:     createTxUseCase,
		getStatementUseCase: getStatementUseCase,
	}
}
func (a *WebApplication) Server(port int) error {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Post("/clientes/:clientID/transacoes", a.handleCreateTransaction)
	app.Get("/clientes/:clientID/extrato", a.handleGetStatement)

	return app.Listen(fmt.Sprintf(":%d", port))
}

func (a *WebApplication) handleCreateTransaction(c *fiber.Ctx) error {
	clientID, err := c.ParamsInt("clientID")
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	var request model.TransactionInput

	err = c.BodyParser(&request)
	if err != nil {
		log.Print(err)
		return c.SendStatus(http.StatusBadRequest)
	}
	resp, err := a.createTxUseCase.Execute(c.Context(), clientID, &request)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.JSON(resp)
}

func (a *WebApplication) handleGetStatement(c *fiber.Ctx) error {
	clientID, err := c.ParamsInt("clientID")
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	resp, err := a.getStatementUseCase.Execute(c.Context(), clientID)
	if err != nil {
		log.Print(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.JSON(resp)
}
