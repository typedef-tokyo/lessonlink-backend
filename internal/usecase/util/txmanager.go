package util

import (
	"context"
	"database/sql"
)

type TxManager interface {
	Do(ctx context.Context, fn func(tx *sql.Tx) error) error
}
