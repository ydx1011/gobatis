package transaction

import (
	"context"
	"database/sql"
	"github.com/ydx1011/gobatis/common"
	"github.com/ydx1011/gobatis/connection"
	"github.com/ydx1011/gobatis/datasource"
	"github.com/ydx1011/gobatis/errors"
	"github.com/ydx1011/gobatis/reflection"
	"github.com/ydx1011/gobatis/statement"
	"github.com/ydx1011/gobatis/util"
)

type DefaultTransaction struct {
	ds datasource.DataSource
	db *sql.DB
	tx *sql.Tx
}

func NewDefaultTransaction(ds datasource.DataSource, db *sql.DB) *DefaultTransaction {
	ret := &DefaultTransaction{ds: ds, db: db}
	return ret
}

func (trans *DefaultTransaction) Close() {

}

func (trans *DefaultTransaction) GetConnection() connection.Connection {
	if trans.tx == nil {
		return (*connection.DefaultConnection)(trans.db)
	} else {
		return &TransactionConnection{tx: trans.tx}
	}
}

func (trans *DefaultTransaction) Begin() error {
	tx, err := trans.db.Begin()
	if err != nil {
		return err
	}
	trans.tx = tx
	return nil
}

func (trans *DefaultTransaction) Commit() error {
	if trans.tx == nil {
		return errors.TRANSACTION_WITHOUT_BEGIN
	}

	err := trans.tx.Commit()
	if err != nil {
		return errors.TRANSACTION_COMMIT_ERROR
	}
	return nil
}

func (trans *DefaultTransaction) Rollback() error {
	if trans.tx == nil {
		return errors.TRANSACTION_WITHOUT_BEGIN
	}

	err := trans.tx.Rollback()
	if err != nil {
		return errors.TRANSACTION_COMMIT_ERROR
	}
	return nil
}

type TransactionConnection struct {
	tx *sql.Tx
}

type TransactionStatement struct {
	tx  *sql.Tx
	sql string
}

func (c *TransactionConnection) Prepare(sqlStr string) (statement.Statement, error) {
	ret := &TransactionStatement{
		tx:  c.tx,
		sql: sqlStr,
	}
	return ret, nil
}

func (c *TransactionConnection) Query(ctx context.Context, result reflection.Object, sqlStr string, params ...interface{}) error {
	db := c.tx
	rows, err := db.QueryContext(ctx, sqlStr, params...)
	if err != nil {
		return errors.STATEMENT_QUERY_ERROR
	}
	defer rows.Close()

	util.ScanRows(rows, result)
	return nil
}

func (c *TransactionConnection) Exec(ctx context.Context, sqlStr string, params ...interface{}) (common.Result, error) {
	db := c.tx
	return db.ExecContext(ctx, sqlStr, params...)
}

func (s *TransactionStatement) Query(ctx context.Context, result reflection.Object, params ...interface{}) error {
	rows, err := s.tx.QueryContext(ctx, s.sql, params...)
	if err != nil {
		return errors.STATEMENT_QUERY_ERROR
	}
	defer rows.Close()

	util.ScanRows(rows, result)
	return nil
}

func (s *TransactionStatement) Exec(ctx context.Context, params ...interface{}) (common.Result, error) {
	return s.tx.ExecContext(ctx, s.sql, params...)
}

func (s *TransactionStatement) Close() {
	//Will be closed when commit or rollback
}
