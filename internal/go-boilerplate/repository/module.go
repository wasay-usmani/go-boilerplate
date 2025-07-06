package repository

import (
	"context"
	"database/sql"
)

type R interface {
	WrapAtomic(ctx context.Context, fn func(txn *sql.Tx) error) error
}

type Module struct {
}

func NewModule() *Module {
	return &Module{}
}
