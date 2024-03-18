package executor

import (
	"context"
	"github.com/ydx1011/gobatis/common"
	"github.com/ydx1011/gobatis/reflection"
)

type Executor interface {
	Close(rollback bool)

	Query(ctx context.Context, result reflection.Object, sql string, params ...interface{}) error

	Exec(ctx context.Context, sql string, params ...interface{}) (common.Result, error)

	Begin() error

	Commit(require bool) error

	Rollback(require bool) error
}
