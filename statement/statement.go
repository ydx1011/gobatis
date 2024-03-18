package statement

import (
	"context"
	"github.com/ydx1011/gobatis/common"
	"github.com/ydx1011/gobatis/reflection"
)

type Statement interface {
	Query(ctx context.Context, result reflection.Object, params ...interface{}) error
	Exec(ctx context.Context, params ...interface{}) (common.Result, error)
	Close()
}
