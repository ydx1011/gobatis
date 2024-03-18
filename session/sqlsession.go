package session

import (
	"context"
	"github.com/ydx1011/gobatis/reflection"
)

type SqlSession interface {
	Close(rollback bool)

	Query(ctx context.Context, result reflection.Object, sql string, params ...interface{}) error

	Insert(ctx context.Context, sql string, params ...interface{}) (int64, int64, error)

	Update(ctx context.Context, sql string, params ...interface{}) (int64, error)

	Delete(ctx context.Context, sql string, params ...interface{}) (int64, error)

	Begin() error

	Commit() error

	Rollback() error
}
