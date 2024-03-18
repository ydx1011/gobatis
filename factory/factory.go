package factory

import (
	"github.com/ydx1011/gobatis/datasource"
	"github.com/ydx1011/gobatis/executor"
	"github.com/ydx1011/gobatis/logging"
	"github.com/ydx1011/gobatis/session"
	"github.com/ydx1011/gobatis/transaction"
)

type Factory interface {
	Open(datasource.DataSource) error
	Close() error

	GetDataSource() datasource.DataSource

	CreateTransaction() transaction.Transaction
	CreateExecutor(transaction.Transaction) executor.Executor
	CreateSession() session.SqlSession
	LogFunc() logging.LogFunc
}
