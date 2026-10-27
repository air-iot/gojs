package gojs

import (
	"reflect"

	"github.com/dop251/goja"
)

// IsBool 判断 goja.Value 是否为 bool 类型
func IsBool(value goja.Value) bool {
	if !IsValid(value) {
		return false
	}
	switch value.ExportType().Kind() {
	case reflect.Bool:
		return true
	default:
		return false
	}
}

// IsString 判断 goja.Value 是否为 string 类型
func IsString(value goja.Value) bool {
	if !IsValid(value) {
		return false
	}
	switch value.ExportType().Kind() {
	case reflect.String:
		return true
	default:
		return false
	}
}

// IsObject 判断 goja.Value 是否为 Object 类型
func IsObject(value goja.Value) bool {
	if !IsValid(value) {
		return false
	}
	switch value.ExportType().Kind() {
	case reflect.Map:
		return true
	default:
		return false
	}
}

// IsArray 判断 goja.Value 是否为 Array 类型
func IsArray(value goja.Value) bool {
	if !IsValid(value) {
		return false
	}
	switch value.ExportType().Kind() {
	case reflect.Array, reflect.Slice:
		return true
	default:
		return false
	}
}

// IsObjectArray 判断 goja.Value 是否为 Array<Object> 类型
func IsObjectArray(value goja.Value) bool {
	if !IsValid(value) {
		return false
	}
	vt := value.ExportType()
	switch vt.Kind() {
	case reflect.Array, reflect.Slice:
		switch vt.Elem().Kind() {
		case reflect.Map:
			return true
		case reflect.Interface:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

// IsStringArray 判断 goja.Value 是否为 Array<string> 类型
func IsStringArray(value goja.Value) bool {
	if !IsValid(value) {
		return false
	}
	vt := value.ExportType()
	switch vt.Kind() {
	case reflect.Array, reflect.Slice:
		switch vt.Elem().Kind() {
		case reflect.Interface:
			return true
		case reflect.String:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

// GetAsArrayObject 将 value 转换为 []map[string]any, 如果类型不匹配返回 nil, false
func GetAsArrayObject(value goja.Value) ([]map[string]any, bool) {
	if !IsObjectArray(value) {
		return nil, false
	}

	values := value.Export()
	switch v := values.(type) {
	case []map[string]any:
		return v, true
	case []any:
		objs := make([]map[string]any, len(v))
		for i := range v {
			if vv, ok := v[i].(map[string]any); ok {
				objs[i] = vv
			} else {
				return nil, false
			}
		}
		return objs, true
	}
	return nil, false
}

// GetAsArrayString 将 value 转换为 []map[string]any, 如果类型不匹配返回 nil, false
func GetAsArrayString(value goja.Value) ([]string, bool) {
	if !IsStringArray(value) {
		return nil, false
	}

	values := value.Export()
	switch v := values.(type) {
	case []string:
		return v, true
	case []any:
		objs := make([]string, len(v))
		for i := range v {
			if vv, ok := v[i].(string); ok {
				objs[i] = vv
			} else {
				return nil, false
			}
		}
		return objs, true
	}
	return nil, false
}
