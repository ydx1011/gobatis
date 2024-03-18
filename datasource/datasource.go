package datasource

type DataSource interface {
	DriverName() string
	DriverInfo() string
}
