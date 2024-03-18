package gobatis

import (
	"github.com/ydx1011/gobatis/factory"
	"github.com/ydx1011/gobatis/logging"
)

type FacOpt func(f *factory.DefaultFactory)

func NewFactory(opts ...FacOpt) factory.Factory {
	f, _ := CreateFactory(opts...)
	return f
}

func CreateFactory(opts ...FacOpt) (factory.Factory, error) {
	f := &factory.DefaultFactory{
		Log: logging.DefaultLogf,
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(f)
		}
	}

	err := f.Open(f.DataSource)
	if err != nil {
		return nil, err
	}

	return f, nil
}
