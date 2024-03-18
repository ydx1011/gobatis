package reflection

import (
	"fmt"
	"reflect"
	"strconv"
)

type paramParser struct {
	ret   map[string]interface{}
	index int
}

func ParseParams(params ...interface{}) map[string]interface{} {
	parser := paramParser{
		ret:   map[string]interface{}{},
		index: 0,
	}
	parser.innerParse(params...)
	return parser.ret
}

func (parser *paramParser) innerParse(params ...interface{}) {
	for i := range params {
		parser.parseOne("", params[i])
	}
}

func (parser *paramParser) parseOne(parentKey string, v interface{}) {
	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	if IsSimpleType(rt) {
		if parentKey == "" {
			parser.ret[parentKey+strconv.Itoa(parser.index)] = v
			parser.index++
		} else {
			parser.ret[parentKey[:len(parentKey)-1]] = v
		}
	} else if rt.Kind() == reflect.Struct {
		oi, _ := GetStructInfo(v)
		structMap := oi.MapValue()
		for key, value := range structMap {
			parser.ret[parentKey+structKey(oi, key)] = value
		}
	} else if rt.Kind() == reflect.Slice {
		l := rv.Len()
		for i := 0; i < l; i++ {
			elemV := rv.Index(i)
			if !elemV.CanInterface() {
				elemV = reflect.Indirect(elemV)
			}
			parser.parseOne(fmt.Sprintf("%s%d[%d].", parentKey, parser.index, i), elemV.Interface())
		}
		parser.ret[strconv.Itoa(parser.index)] = l
		parser.index++
		//l := rv.Len()
		//builder := strings.Builder{}
		//for i := 0; i < l; i++ {
		//	elemV := rv.Index(i)
		//	if !elemV.CanInterface() {
		//		elemV = reflect.Indirect(elemV)
		//	}
		//	if elemV.Kind() == reflect.String {
		//		builder.WriteString(elemV.String())
		//	} else {
		//		var str string
		//		if SafeSetValue(reflect.ValueOf(&str), elemV.Interface()) {
		//			builder.WriteString(str)
		//		} else {
		//			//log
		//		}
		//	}
		//
		//	if i < l-1 {
		//		builder.WriteString(slice_param_separator)
		//	}
		//}
		//parser.ret[strconv.Itoa(parser.index)] = builder.String()
		//parser.index++
	} else if rt.Kind() == reflect.Map {
		keys := rv.MapKeys()
		for _, key := range keys {
			if key.Kind() == reflect.String {
				value := rv.MapIndex(key)
				value = value.Elem()
				if IsSimpleType(value.Type()) {
					if !value.CanInterface() {
						value = reflect.Indirect(value)
					}
					parser.ret[parentKey+key.String()] = value
				}
			}
		}
	}
}

//是否是数据库使用的简单类型，注意不能是PTR
func IsSimpleType(t reflect.Type) bool {
	switch t.Kind() {
	case IntKind, Int8Kind, Int16Kind, Int32Kind, Int64Kind, UintKind, Uint8Kind, Uint16Kind, Uint32Kind, Uint64Kind,
		Float32Kind, Float64Kind, Complex64Kind, Complex128Kind, StringKind, BoolKind, ByteKind /*, BytesKind, TimeKind*/ :
		return true
	}

	if t.ConvertibleTo(BytesType) || t.ConvertibleTo(TimeType) {
		return true
	}
	return false
}

func structKey(oi *StructInfo, field string) string {
	return oi.Name + "." + field
}
