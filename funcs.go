package gojs

import "github.com/dop251/goja"

// IsValid 判断是否为一个有效值
// 如果 value 为 nil, 或者为 js 中的 undefined, null 值返回 false
func IsValid(value goja.Value) bool {
	return !(value == nil || goja.IsUndefined(value) || goja.IsNaN(value) || goja.IsNull(value))
}

// IsBuffer 判断是否为 Buffer 对象
func IsBuffer(value goja.Value) bool {
	if !IsValid(value) {
		return false
	}
	obj, ok := value.(*goja.Object)
	if !ok {
		return false
	}
	return obj.Get("buffer") != nil
}
