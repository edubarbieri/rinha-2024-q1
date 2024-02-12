package internal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/edubarbieri/rinha-2024-q1/internal/entity"
	"github.com/edubarbieri/rinha-2024-q1/internal/repository"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type WebApplication struct {
	repo repository.Repository
}

func NewWebApplication(repo repository.Repository) *WebApplication {
	return &WebApplication{
		repo: repo,
	}
}
func (a *WebApplication) Server(port int) error {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	app.Post("/clientes/:clientID/transacoes", a.handleCreateTransaction)
	app.Get("/clientes/:clientID/extrato", a.handleGetStatement)

	return app.Listen(fmt.Sprintf(":%d", port))
}

func (a *WebApplication) handleCreateTransaction(c *fiber.Ctx) error {
	clientID, err := c.ParamsInt("clientID")
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	var request entity.TransactionInput

	err = c.BodyParser(&request)
	if err != nil {
		log.Print(err)
		return c.SendStatus(http.StatusUnprocessableEntity)
	}

	if request.Validate() != nil {
		return c.SendStatus(http.StatusUnprocessableEntity)
	}

	resp, err := a.repo.SaveTransaction(c.Context(), clientID, &request)
	if err != nil {
		return c.SendStatus(a.errorStatusCode(err))
	}

	return c.JSON(resp)
}

func (a *WebApplication) handleGetStatement(c *fiber.Ctx) error {
	clientID, err := c.ParamsInt("clientID")
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	resp, err := a.repo.GetUserStatement(c.Context(), clientID)
	if err != nil {
		return c.SendStatus(a.errorStatusCode(err))
	}

	return c.JSON(resp)
}

func (a *WebApplication) errorStatusCode(err error) int {

	switch err {
	case repository.ErrClientNotExist:
		return http.StatusNotFound
	case repository.ErrLimitExceeded:
		return http.StatusUnprocessableEntity
	default:
		log.Println(err)
		return http.StatusInternalServerError
	}

}
