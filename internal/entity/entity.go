package entity

import (
	"fmt"
	"time"
)

type TransactionInput struct {
	Value       int    `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

func (t *TransactionInput) Validate() error {
	if t.Value <= 0 {
		return fmt.Errorf("\"valor\" deve ser maior que 0")
	}

	if t.Type != "c" && t.Type != "d" {
		return fmt.Errorf("\"tipo\" deve ser \"c\" ou \"d\"")
	}

	descLen := len(t.Description)
	if descLen < 1 || descLen > 10 {
		return fmt.Errorf("\"descricao\" deve possuir de 1 a 10 caractéres")
	}

	return nil
}

type TransactionOutput struct {
	Limit   int `json:"limite"`
	Balance int `json:"saldo"`
}

type StatementOutput struct {
	Balance      *Balance       `json:"saldo"`
	Transactions []*Transaction `json:"ultimas_transacoes"`
}

type Balance struct {
	Total int       `json:"total"`
	Date  time.Time `json:"data_extrato"`
	Limit int       `json:"limite"`
}

type Transaction struct {
	Value       int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	Date        time.Time `json:"realizada_em"`
}
