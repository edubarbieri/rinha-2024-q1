package model

import "time"

type StatementOutput struct {
	Balance      Balance       `json:"saldo"`
	Transactions []Transaction `json:"ultimas_transacoes"`
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
