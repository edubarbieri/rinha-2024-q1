package repository

import (
	"context"
	"testing"

	"github.com/edubarbieri/rinha-2024-q1/internal/application/model"
	"github.com/stretchr/testify/assert"
)

func TestMysqlRepository(t *testing.T) {
	underTest, err := NewMysqlRepository("root", "", "rinha", "localhost:3306")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should return true if client existe", func(t *testing.T) {
		result, err := underTest.UserExists(context.Background(), 1)
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return false if client existe", func(t *testing.T) {
		result, err := underTest.UserExists(context.Background(), 10)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should create a transction", func(t *testing.T) {
		// Arrange
		clientID := 1
		tx := &model.TransactionInput{
			Value:       100,
			Type:        "c",
			Description: "test",
		}
		err := underTest.SaveTransaction(context.Background(), clientID, tx)
		assert.NoError(t, err)
	})

	t.Run("should return statement", func(t *testing.T) {
		clientID := 1
		result, err := underTest.GetUserStatement(context.Background(), clientID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}
