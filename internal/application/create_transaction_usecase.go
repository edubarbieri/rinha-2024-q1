package application

import (
	"context"

	"github.com/edubarbieri/rinha-2024-q1/internal/application/model"
)

type CreateTransactionUseCase struct {
}

func NewCreateTransactionUseCase() *CreateTransactionUseCase {
	return &CreateTransactionUseCase{}
}

func (c *CreateTransactionUseCase) Execute(
	ctx context.Context, userID int, transaction *model.TransactionInput) (*model.TransactionOutput, error) {

	return &model.TransactionOutput{
		Limit:   11,
		Balance: 100,
	}, nil
}
