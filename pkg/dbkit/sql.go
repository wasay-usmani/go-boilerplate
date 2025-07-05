package dbkit

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/wasay-usmani/go-boilerplate/pkg/logkit"
)

type SqlConn struct {
	*sql.DB
	logger logkit.Logger
}

func (c *SqlConn) Ping() error {
	return c.DB.Ping()
}

func (c *SqlConn) Close() error {
	return c.DB.Close()
}

func (c *SqlConn) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	txn, err := c.DB.BeginTx(ctx, opts)
	if err != nil {
		c.logger.Error("begin transaction error", err)
		return nil, fmt.Errorf("begin transaction failure %w", err)
	}

	return txn, nil
}

func (c *SqlConn) Atomic(ctx context.Context, opts *sql.TxOptions, fn func(txn *sql.Tx) error) (err error) {
	txn, err := c.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		if e := recover(); e != nil {
			err = c.resolveTxn(txn, fmt.Errorf("panic occurred, cause: %v", e))
			return
		}

		err = c.resolveTxn(txn, err)
	}()

	if err := fn(txn); err != nil {
		return err
	}

	return nil
}

func (c *SqlConn) resolveTxn(txn *sql.Tx, err error) error {
	if err != nil {
		if rollbackErr := txn.Rollback(); rollbackErr != nil {
			c.logger.Error("transaction rollback failed", err)
			return rollbackErr
		}

		if !errors.Is(err, sql.ErrNoRows) {
			c.logger.Error("transaction operation failed", err)
		}

		return err
	}

	err = txn.Commit()
	if err != nil {
		c.logger.Error("transaction commit failed", err)
		return err
	}

	return nil
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type AnyString struct{}

func (a AnyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}

func NewMockDB(l logkit.Logger) (*SqlConn, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		l.Error("unexpected error while opening a mock db connection", err)
	}

	return &SqlConn{db, l}, mock
}
