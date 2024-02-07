package application

import (
	"context"

	"github.com/edubarbieri/rinha-2024-q1/internal/application/model"
)

type GetStatementUseCase struct {
}

func NewGetStatementUseCase() *GetStatementUseCase {
	return &GetStatementUseCase{}
}

func (c *GetStatementUseCase) Execute(
	ctx context.Context, userID int) (*model.StatementOutput, error) {

	return &model.StatementOutput{}, nil
}
