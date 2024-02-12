package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/edubarbieri/rinha-2024-q1/internal/entity"
	"github.com/go-sql-driver/mysql"
)

type mysqlRepository struct {
	db          *sql.DB
	mu          sync.Mutex
	clientCache map[int]bool
}

var ErrClientNotExist = fmt.Errorf("client not exists")
var ErrLimitExceeded = fmt.Errorf("client does not have sufficient limit")

func NewMysqlRepository(dbUser, dbPass, dbName, dbAddress string) (Repository, error) {
	// Capture connection properties.
	cfg := mysql.Config{
		User:                 dbUser,
		Passwd:               dbPass,
		Net:                  "tcp",
		Addr:                 dbAddress,
		DBName:               dbName,
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	var pingErr error
	for range 10 {
		pingErr = db.Ping()
		if pingErr != nil {
			log.Println("waiting db")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	if pingErr != nil {
		return nil, pingErr
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)

	return &mysqlRepository{
		db:          db,
		clientCache: make(map[int]bool, 0),
	}, nil
}

func (m *mysqlRepository) clientExists(ctx context.Context, clientID int) (bool, error) {
	if exist, ok := m.clientCache[clientID]; ok {
		return exist, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if exist, ok := m.clientCache[clientID]; ok {
		return exist, nil
	}

	var clientDB int
	err := m.db.QueryRowContext(ctx, "select id from clients where id = ?", clientID).Scan(&clientDB)
	switch {
	case err == sql.ErrNoRows:
		m.clientCache[clientID] = false
		return false, nil
	case err != nil:
		return false, err
	default:
		m.clientCache[clientID] = true
		return true, nil
	}
}

func (m *mysqlRepository) SaveTransaction(ctx context.Context, clientID int, data *entity.TransactionInput) (*entity.TransactionOutput, error) {
	exists, err := m.clientExists(ctx, clientID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrClientNotExist
	}

	err = WithTransaction(ctx, m.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, "CALL create_transaction(?, ?, ?, ?)",
			clientID, data.Value, data.Type, data.Description)

		if err != nil {
			if strings.Contains(err.Error(), "check_balance_positive") {
				return ErrLimitExceeded
			}
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	var output entity.TransactionOutput
	err = m.db.QueryRowContext(ctx, "SELECT balance, c_limit from clients where id = ?", clientID).
		Scan(&output.Balance, &output.Limit)
	if err != nil {
		return nil, err
	}

	return &output, nil

}

func (m *mysqlRepository) GetUserStatement(ctx context.Context, clientID int) (*entity.StatementOutput, error) {
	exists, err := m.clientExists(ctx, clientID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrClientNotExist
	}

	balance := entity.Balance{
		Date: time.Now(),
	}
	err = m.db.QueryRowContext(ctx, "SELECT balance, c_limit from clients where id = ?", clientID).
		Scan(&balance.Total, &balance.Limit)

	if err != nil {
		return nil, err
	}

	rows, err := m.db.QueryContext(ctx, "select value, type, description, create_at from transactions where client_id = ? order by create_at DESC LIMIT 10", clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := make([]*entity.Transaction, 0)

	for rows.Next() {
		var tx entity.Transaction
		err = rows.Scan(&tx.Value, &tx.Type, &tx.Description, &tx.Date)
		if err != nil {
			return nil, err
		}
		txs = append(txs, &tx)
	}

	return &entity.StatementOutput{
		Balance:      &balance,
		Transactions: txs,
	}, nil
}
