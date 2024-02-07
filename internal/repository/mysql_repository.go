package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/edubarbieri/rinha-2024-q1/internal/application/model"
	"github.com/go-sql-driver/mysql"
)

type mysqlRepository struct {
	db *sql.DB
}

func NewMysqlRepository(dbUser, dbPass, dbName, dbAddress string) (Repository, error) {
	// Capture connection properties.
	cfg := mysql.Config{
		User:      dbUser,
		Passwd:    dbPass,
		Net:       "tcp",
		Addr:      dbAddress,
		DBName:    dbName,
		ParseTime: true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	fmt.Println("Mysql connected success!")
	return &mysqlRepository{
		db: db,
	}, nil
}

func (m *mysqlRepository) UserExists(ctx context.Context, clientID int) (bool, error) {
	var clientDB int
	err := m.db.QueryRowContext(ctx, "select id from clients where id = ?", clientID).Scan(&clientDB)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (m *mysqlRepository) SaveTransaction(ctx context.Context, clientID int, data *model.TransactionInput) error {
	return WithTransaction(ctx, m.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE clients set balance = balance - ? where id =?",
			data.Value, clientID)

		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx,
			"INSERT INTO transactions (client_id, value, type, description) VALUES (?, ?, ?, ?)",
			clientID, data.Value, data.Type, data.Description)

		if err != nil {
			return err
		}
		return nil
	})
}

func (m *mysqlRepository) GetUserStatement(ctx context.Context, clientID int) (*model.StatementOutput, error) {
	balance := model.Balance{
		Date: time.Now(),
	}
	err := m.db.QueryRowContext(ctx, "SELECT balance, c_limit from clients where id = ?", clientID).
		Scan(&balance.Total, &balance.Limit)

	if err != nil {
		return nil, err
	}

	rows, err := m.db.QueryContext(ctx, "select value, type, description, create_at from transactions where client_id = ?", clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := make([]*model.Transaction, 0)

	for rows.Next() {
		var tx model.Transaction
		err = rows.Scan(&tx.Value, &tx.Type, &tx.Description, &tx.Date)
		if err != nil {
			return nil, err
		}
		txs = append(txs, &tx)
	}

	return &model.StatementOutput{
		Balance:      &balance,
		Transactions: txs,
	}, nil
}
