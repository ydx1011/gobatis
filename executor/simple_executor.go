package executor

import (
	"context"
	"github.com/ydx1011/gobatis/common"
	"github.com/ydx1011/gobatis/errors"
	"github.com/ydx1011/gobatis/reflection"
	"github.com/ydx1011/gobatis/transaction"
)

type SimpleExecutor struct {
	transaction transaction.Transaction
	closed      bool
}

func NewSimpleExecutor(transaction transaction.Transaction) *SimpleExecutor {
	return &SimpleExecutor{transaction: transaction}
}

func (exec *SimpleExecutor) Close(rollback bool) {
	defer func() {
		if exec.transaction != nil {
			exec.transaction.Close()
		}
		exec.transaction = nil
		exec.closed = true
	}()

	if rollback {
		exec.Rollback(true)
	}
}

func (exec *SimpleExecutor) Query(ctx context.Context, result reflection.Object, sql string, params ...interface{}) error {
	if exec.closed {
		return errors.EXECUTOR_QUERY_ERROR
	}

	conn := exec.transaction.GetConnection()
	if conn == nil {
		return errors.EXECUTOR_GET_CONNECTION_ERROR
	}

	return conn.Query(ctx, result, sql, params...)
}

func (exec *SimpleExecutor) Exec(ctx context.Context, sql string, params ...interface{}) (common.Result, error) {
	if exec.closed {
		return nil, errors.EXECUTOR_QUERY_ERROR
	}

	conn := exec.transaction.GetConnection()
	if conn == nil {
		return nil, errors.EXECUTOR_GET_CONNECTION_ERROR
	}

	return conn.Exec(ctx, sql, params...)
}

func (exec *SimpleExecutor) Begin() error {
	if exec.closed {
		return errors.EXECUTOR_BEGIN_ERROR
	}

	return exec.transaction.Begin()
}

func (exec *SimpleExecutor) Commit(require bool) error {
	if exec.closed {
		return errors.EXECUTOR_COMMIT_ERROR
	}

	if require {
		return exec.transaction.Commit()
	}

	return nil
}

func (exec *SimpleExecutor) Rollback(require bool) error {
	if !exec.closed {
		if require {
			return exec.transaction.Rollback()
		}
	}
	return nil
}
