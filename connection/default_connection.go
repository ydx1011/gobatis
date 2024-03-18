package connection

import (
	"context"
	"database/sql"
	"github.com/ydx1011/gobatis/common"
	"github.com/ydx1011/gobatis/errors"
	"github.com/ydx1011/gobatis/reflection"
	"github.com/ydx1011/gobatis/statement"
	"github.com/ydx1011/gobatis/util"
)

type DefaultConnection sql.DB
type DefaultStatement sql.Stmt

func (c *DefaultConnection) Prepare(sqlStr string) (statement.Statement, error) {
	db := (*sql.DB)(c)
	s, err := db.Prepare(sqlStr)
	if err != nil {
		return nil, errors.CONNECTION_PREPARE_ERROR
	}
	return (*DefaultStatement)(s), nil
}

func (c *DefaultConnection) Query(ctx context.Context, result reflection.Object, sqlStr string, params ...interface{}) error {
	db := (*sql.DB)(c)
	rows, err := db.QueryContext(ctx, sqlStr, params...)
	if err != nil {
		return errors.STATEMENT_QUERY_ERROR
	}
	defer rows.Close()

	util.ScanRows(rows, result)
	return nil
}

func (c *DefaultConnection) Exec(ctx context.Context, sqlStr string, params ...interface{}) (common.Result, error) {
	db := (*sql.DB)(c)
	return db.ExecContext(ctx, sqlStr, params...)
}

func (s *DefaultStatement) Query(ctx context.Context, result reflection.Object, params ...interface{}) error {
	stmt := (*sql.Stmt)(s)
	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return errors.STATEMENT_QUERY_ERROR
	}
	defer rows.Close()

	util.ScanRows(rows, result)
	return nil
}

func (s *DefaultStatement) Exec(ctx context.Context, params ...interface{}) (common.Result, error) {
	stmt := (*sql.Stmt)(s)
	return stmt.ExecContext(ctx, params...)
}

func (s *DefaultStatement) Close() {
	stmt := (*sql.Stmt)(s)
	stmt.Close()
}
