package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionInput_Validate(t *testing.T) {
	t.Run("should return nil if have no validation erros", func(t *testing.T) {
		// Arrange
		underTest := TransactionInput{
			Value:       200,
			Type:        "c",
			Description: "test",
		}
		// Act
		result := underTest.Validate()
		// Assert
		assert.NoError(t, result)
	})

	t.Run("should validate if value is greater than zero", func(t *testing.T) {
		// Arrange
		underTest := TransactionInput{
			Value:       0,
			Type:        "c",
			Description: "test",
		}
		// Act
		result := underTest.Validate()
		// Assert
		assert.Error(t, result)
	})

	t.Run("should validate if type ", func(t *testing.T) {
		// Arrange
		underTest := TransactionInput{
			Value:       100,
			Type:        "r",
			Description: "test",
		}
		// Act
		result := underTest.Validate()
		// Assert
		assert.Error(t, result)
	})

	t.Run("should validate if description", func(t *testing.T) {
		// Arrange
		underTest := TransactionInput{
			Value:       100,
			Type:        "r",
			Description: "test dd dd d ds",
		}
		// Act
		result := underTest.Validate()
		// Assert
		assert.Error(t, result)
	})
}
