package datasource

type CommonDataSource struct {
	Name string
	Info string
}

func (ds *CommonDataSource) DriverName() string {
	return ds.Name
}

func (ds *CommonDataSource) DriverInfo() string {
	return ds.Info
}
