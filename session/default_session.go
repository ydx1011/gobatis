package session

import (
	"context"
	"fmt"
	"github.com/ydx1011/gobatis/common"
	"github.com/ydx1011/gobatis/executor"
	"github.com/ydx1011/gobatis/logging"
	"github.com/ydx1011/gobatis/reflection"
	"github.com/ydx1011/gobatis/transaction"
)

type DefaultSqlSession struct {
	Log        logging.LogFunc
	tx         transaction.Transaction
	executor   executor.Executor
	autoCommit bool
}

func NewDefaultSqlSession(log logging.LogFunc, tx transaction.Transaction, e executor.Executor, autoCommit bool) *DefaultSqlSession {
	ret := &DefaultSqlSession{
		Log:        log,
		tx:         tx,
		executor:   e,
		autoCommit: autoCommit,
	}
	return ret
}

func (sess *DefaultSqlSession) Close(rollback bool) {
	sess.executor.Close(rollback)
}

func (sess *DefaultSqlSession) Query(ctx context.Context, result reflection.Object, sql string, params ...interface{}) error {
	sess.logLastSql(sql, params...)
	return sess.executor.Query(ctx, result, sql, params...)
}

func (sess *DefaultSqlSession) Insert(ctx context.Context, sql string, params ...interface{}) (int64, int64, error) {
	sess.logLastSql(sql, params...)
	ret, err := sess.exec(ctx, sql, params...)
	if err != nil {
		return 0, -1, err
	}

	count, err := ret.RowsAffected()
	if err != nil {
		return 0, -1, err
	}

	id, err := ret.LastInsertId()
	if err != nil {
		return count, id, err
	}
	return count, id, nil
}

func (sess *DefaultSqlSession) Update(ctx context.Context, sql string, params ...interface{}) (int64, error) {
	sess.logLastSql(sql, params...)
	ret, err := sess.exec(ctx, sql, params...)
	if err != nil {
		return 0, err
	}

	count, err := ret.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (sess *DefaultSqlSession) Delete(ctx context.Context, sql string, params ...interface{}) (int64, error) {
	sess.logLastSql(sql, params...)
	ret, err := sess.exec(ctx, sql, params...)
	if err != nil {
		return 0, err
	}

	count, err := ret.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (sess *DefaultSqlSession) Begin() error {
	sess.logLastSql("Begin", "")
	return sess.tx.Begin()
}

func (sess *DefaultSqlSession) Commit() error {
	sess.logLastSql("Commit", "")
	return sess.tx.Commit()
}

func (sess *DefaultSqlSession) Rollback() error {
	sess.logLastSql("Rollback", "")
	return sess.tx.Rollback()
}

func (sess *DefaultSqlSession) logLastSql(sql string, params ...interface{}) {
	sess.Log(logging.INFO, "sql: [%s], param: %s\n", sql, fmt.Sprint(params))
}

func (sess *DefaultSqlSession) exec(ctx context.Context, sql string, params ...interface{}) (common.Result, error) {
	return sess.executor.Exec(ctx, sql, params...)
}
