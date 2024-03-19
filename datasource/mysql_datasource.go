package datasource

import "fmt"

type MysqlDataSource struct {
	Host     string
	Port     int
	DBName   string
	Username string
	Password string
	Charset  string
}

func (ds *MysqlDataSource) DriverName() string {
	return "mysql"
}

func (ds *MysqlDataSource) DriverInfo() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", ds.Username, ds.Password, ds.Host, ds.Port, ds.DBName, ds.Charset)
}
