package transaction

import "github.com/ydx1011/gobatis/connection"

type Transaction interface {
	Close()

	GetConnection() connection.Connection

	Begin() error

	Commit() error

	Rollback() error
}
