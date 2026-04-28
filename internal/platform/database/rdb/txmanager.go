package rdb

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database"
)

type TxManager struct {
	db *sql.DB
}

func NewTxManager(m IMySQL) database.TxManager {
	return &TxManager{db: m.GetConn()}
}

func (m *TxManager) Do(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {

		return err
	}

	return tx.Commit()
}
