package gobatis

import (
	"github.com/ydx1011/gobatis/reflection"
	"reflect"
)

func init() {
	var typeModelName ModelName
	reflection.SetModelNameType(reflect.TypeOf(typeModelName))
}
