package router

import (
	"log"
	"reflect"
	"strconv"
	"strings"
)

// 参数校验器
var validMap = make(map[string]func(f reflect.StructField, v reflect.Value, other string))

func init() {
	validMap["not nil"] = notNilHandler
	validMap["not black"] = notBlackHandler
	validMap["len"] = lenHandler
	validMap["min"] = minHandler
}

// 不能为空字符串
func notBlackHandler(f reflect.StructField, v reflect.Value, length string){
	t := f.Type
	if t.Kind() == reflect.Pointer {
		if v.IsNil() {
			return
		}
		t = t.Elem()
	}

	if t.Kind() != reflect.String {
		return
	}

	if v.String() == "" {
		jsonTag := f.Tag.Get("json")
		split := strings.Split(jsonTag, ",")
		fieldName := f.Name
		if len(split) > 0 {
			fieldName = split[0]
		}
		log.Panicf("字段 %s 值不能是空字符串", fieldName)
	}

}

// 最小值
func minHandler(f reflect.StructField, v reflect.Value, length string) {
	if length == "" {
		return
	}
	t := f.Type
	if f.Type.Kind() == reflect.Pointer {
		if v.IsNil() {
			return
		}
		t = t.Elem()
	}

	if t.Kind() != reflect.Int {
		return
	}

	n, err := strconv.Atoi(length)
	if err != nil {
		log.Panic(err)
	}
	if v.Int() < int64(n) {
		jsonTag := f.Tag.Get("json")
		split := strings.Split(jsonTag, ",")
		fieldName := f.Name
		if len(split) > 0 {
			fieldName = split[0]
		}
		log.Panicf("字段 %s 值不能小于 %d", fieldName, n)
	}
}

// 指针不能为 nil
func notNilHandler(f reflect.StructField, v reflect.Value, other string) {
	if f.Type.Kind() == reflect.Pointer && v.IsNil() {
		jsonTag := f.Tag.Get("json")
		split := strings.Split(jsonTag, ",")
		fieldName := f.Name
		if len(split) > 0 {
			fieldName = split[0]
		}
		log.Panicf("字段 %s 不能为空", fieldName)
	}
}

// 字符串长度限制
func lenHandler(f reflect.StructField, v reflect.Value, length string) {
	if length == "" {
		return
	}

	t := f.Type
	if f.Type.Kind() == reflect.Pointer {
		if v.IsNil() {
			return
		}
		t = f.Type.Elem()
	}

	if t.Kind() != reflect.String {
		return
	}

	// 长度转成 int
	n, err := strconv.Atoi(length)
	if err != nil {
		log.Panic(err)
	}

	s := v.String()
	if len(s) > n {
		jsonTag := f.Tag.Get("json")
		split := strings.Split(jsonTag, ",")
		fieldName := f.Name
		if len(split) > 0 {
			fieldName = split[0]
		}
		log.Panicf("字段 %s 长度不能大于 %d", fieldName, n)
	}

}
