package parsing

import (
	"github.com/ydx1011/gobatis/logging"
	"github.com/ydx1011/gobatis/parsing/sqlparser"
	"github.com/ydx1011/gobatis/reflection"
	"reflect"
	"strings"
	"time"
)

type GetFunc func(key string) string

type DynamicElement interface {
	Format(func(key string) string) string
}

type DynamicData struct {
	OriginData     string
	DynamicElemMap map[string]DynamicElement
}

//需要外部确保param是一个struct
func (m *DynamicData) ReplaceWithMap(objParams map[string]interface{}) string {
	if len(m.DynamicElemMap) == 0 || len(objParams) == 0 {
		logging.Info("map is empty")
		//return m.OriginData
	}

	getFunc := func(s string) string {
		if o, ok := objParams[s]; ok {
			if str, ok := o.(string); ok {
				return str
			}

			//zero time convert to empty string (for <if> </if> element)
			if ti, ok := o.(time.Time); ok {
				if ti.IsZero() {
					return ""
				} else {
					return ti.String()
				}
			}

			var str string
			reflection.SafeSetValue(reflect.ValueOf(&str), o)
			return str
		}
		return ""
	}

	ret := m.OriginData
	for k, v := range m.DynamicElemMap {
		ret = strings.Replace(ret, k, v.Format(getFunc), -1)
	}
	return ret
}

func (m *DynamicData) ParseMetadata(driverName string, params ...interface{}) (*sqlparser.Metadata, error) {
	paramMap := reflection.ParseParams(params...)
	sqlStr := m.ReplaceWithMap(paramMap)
	return sqlparser.ParseWithParamMap(driverName, sqlStr, paramMap)
}
