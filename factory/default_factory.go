package factory

import (
	"database/sql"
	"github.com/ydx1011/gobatis/datasource"
	"github.com/ydx1011/gobatis/errors"
	"github.com/ydx1011/gobatis/executor"
	"github.com/ydx1011/gobatis/logging"
	"github.com/ydx1011/gobatis/session"
	"github.com/ydx1011/gobatis/transaction"
	"sync"
	"time"
)

type DefaultFactory struct {
	MaxConn         int
	MaxIdleConn     int
	ConnMaxLifetime time.Duration
	Log             logging.LogFunc

	DataSource datasource.DataSource

	db    *sql.DB
	mutex sync.Mutex
}

func (f *DefaultFactory) Open(ds datasource.DataSource) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.db != nil {
		return errors.FACTORY_INITED
	}

	if ds != nil {
		f.DataSource = ds
	}

	db, err := sql.Open(f.DataSource.DriverName(), f.DataSource.DriverInfo())
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(f.MaxConn)
	db.SetMaxIdleConns(f.MaxIdleConn)
	db.SetConnMaxLifetime(f.ConnMaxLifetime)

	f.db = db
	return nil
}

func (f *DefaultFactory) Close() error {
	if f.db != nil {
		return f.db.Close()
	}
	return nil
}

func (f *DefaultFactory) GetDataSource() datasource.DataSource {
	return f.DataSource
}

func (f *DefaultFactory) CreateTransaction() transaction.Transaction {
	return transaction.NewDefaultTransaction(f.DataSource, f.db)
}

func (f *DefaultFactory) CreateExecutor(transaction transaction.Transaction) executor.Executor {
	return executor.NewSimpleExecutor(transaction)
}

func (f *DefaultFactory) CreateSession() session.SqlSession {
	tx := f.CreateTransaction()
	return session.NewDefaultSqlSession(f.Log, tx, f.CreateExecutor(tx), false)
}

func (f *DefaultFactory) LogFunc() logging.LogFunc {
	return f.Log
}

func (f *DefaultFactory) WithLock(lockFunc func(fac *DefaultFactory)) {
	f.mutex.Lock()
	lockFunc(f)
	f.mutex.Unlock()
}
