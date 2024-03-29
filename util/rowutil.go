package util

import (
	"database/sql"
	"github.com/ydx1011/gobatis/reflection"
	"reflect"
)

func ScanRows(rows *sql.Rows, result reflection.Object) int64 {
	columns, _ := rows.Columns()

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for j := range values {
		scanArgs[j] = &values[j]
	}

	var index int64 = 0
	for rows.Next() {
		if err := rows.Scan(scanArgs...); err == nil {
			//for _, col := range values {
			//    logging.Debug("%v", col)
			//}
			if !deserialize(result, columns, values) {
				break
			}
			index++
		}
	}
	return index
}

func deserialize(result reflection.Object, columns []string, values []interface{}) bool {
	obj := result
	if result.CanAddValue() {
		obj = result.NewElem()
	}
	for i := range columns {
		if obj.CanSetField() {
			obj.SetField(columns[i], values[i])
		} else {
			obj.SetValue(reflect.ValueOf(values[0]))
			break
		}
	}
	if result.CanAddValue() {
		result.AddValue(obj.GetValue())
		return true
	}
	return false
}
