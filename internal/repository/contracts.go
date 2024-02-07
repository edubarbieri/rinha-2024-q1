package repository

import (
	"context"

	"github.com/edubarbieri/rinha-2024-q1/internal/application/model"
)

type Repository interface {
	SaveTransaction(ctx context.Context, clientID int, data *model.TransactionInput) error
	GetUserStatement(ctx context.Context, clientID int) (*model.StatementOutput, error)
	UserExists(ctx context.Context, clientID int) (bool, error)
}
