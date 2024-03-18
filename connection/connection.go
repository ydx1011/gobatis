package connection

import (
	"context"
	"github.com/ydx1011/gobatis/common"
	"github.com/ydx1011/gobatis/reflection"
	"github.com/ydx1011/gobatis/statement"
)

type Connection interface {
	Prepare(sql string) (statement.Statement, error)
	Query(ctx context.Context, result reflection.Object, sql string, params ...interface{}) error
	Exec(ctx context.Context, sql string, params ...interface{}) (common.Result, error)
}
