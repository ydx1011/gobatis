package common

//参数
//  idx:迭代数
//  bean:序列化后的值
//返回值:
//  打断迭代返回true
type IterFunc func(idx int64, bean interface{}) bool

const (
	FIELD_NAME = "yfield"
)
