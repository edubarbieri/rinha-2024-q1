package repository

import (
	"context"

	"github.com/edubarbieri/rinha-2024-q1/internal/entity"
)

type Repository interface {
	SaveTransaction(ctx context.Context, clientID int, data *entity.TransactionInput) (*entity.TransactionOutput, error)
	GetUserStatement(ctx context.Context, clientID int) (*entity.StatementOutput, error)
}
