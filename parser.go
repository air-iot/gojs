package gojs

import (
	"fmt"
	"github.com/dop251/goja"
)

// ParseResult  数据处理函数返回结果数据格式
type ParseResult struct {
	// 设备ID
	ID string `json:"id"`
	// 设备所属工作表标识
	Table string `json:"table"`
	// 子设备ID
	CID string `json:"cid"`
	// 上数时间
	Time int64 `json:"time"`
	// 数据点信息
	// key: 数据点标识
	// value: 数据点的值
	Values map[string]interface{} `json:"values"`
}

type Parser struct {
	IDKey     string `json:"id"`
	TableKey  string `json:"table"`
	CIDKey    string `json:"cid"`
	TimeKey   string `json:"time"`
	ValuesKey string `json:"values"`
}

func NewParser() *Parser {
	return &Parser{
		IDKey:     "id",
		TableKey:  "table",
		CIDKey:    "cid",
		TimeKey:   "time",
		ValuesKey: "values",
	}
}

// Parse 解析脚本返回值.
// 从返回值中解析出设备编号, 表标识, 时间和数据等信息
func (p Parser) Parse(result goja.Value) ([]ParseResult, error) {
	if !IsValid(result) {
		return nil, nil
	}

	obj, ok := result.(*goja.Object)
	if !ok {
		return nil, fmt.Errorf("脚本返回值不是有效的对象")
	}

	if obj.ClassName() != "Array" {
		return nil, fmt.Errorf("不是有效的数组")
	}

	keys := obj.Keys()
	results := make([]ParseResult, 0, len(keys))
	for _, key := range keys {
		value, ok := obj.Get(key).(*goja.Object)
		if !ok {
			return nil, fmt.Errorf("数组内的元素不是有效的对象")
		}

		id := value.Get(p.IDKey)
		table := value.Get(p.TableKey)
		time := value.Get(p.TimeKey)
		cid := value.Get(p.CIDKey)
		values := value.Get(p.ValuesKey)

		if !IsValid(id) {
			return nil, fmt.Errorf("数组内的元素 [%s] 缺少 id 字段", key)
		}
		if !IsValid(values) {
			return nil, fmt.Errorf("数组内的元素 [%s] 缺少 values 字段", key)
		}

		idStr, ok := id.Export().(string)
		if !ok {
			return nil, fmt.Errorf("数组内的元素 [%s] 的 id 字段 %v 不是有效的字符串", key, id.Export())
		}

		var tableId string
		if IsValid(table) {
			v, ok := table.Export().(string)
			if !ok {
				return nil, fmt.Errorf("数组内的元素 [%s] 的 table 字段 %v 不是有效的字符串", key, table.Export())
			}
			tableId = v
		}

		valuesObj, ok := values.(*goja.Object)
		if !ok {
			return nil, fmt.Errorf("数组内的元素 [%s] 的 values 字段不是有效的对象", key)
		}

		if !IsValid(valuesObj) {
			return nil, fmt.Errorf("数组内的元素 [%s] 的 values 字段不是有效的对象", key)
		}

		cidValue := ""
		if IsValid(cid) {
			v, ok := cid.Export().(string)
			if !ok {
				return nil, fmt.Errorf("数组内的元素 [%s] 的 cid 字段 %v 不是有效的字符串", key, cid.Export())
			}
			cidValue = v
		}

		var timeValue int64
		if IsValid(time) {
			if v, ok := time.Export().(int64); ok {
				timeValue = v
			} else if v, ok := time.Export().(float64); ok {
				timeValue = int64(v)
			} else {
				return nil, fmt.Errorf("数组内的元素 [%s] 的 time 字段 %v 不是有效的时间戳(ms)", key, time.Export())
			}
		}

		valueKeys := valuesObj.Keys()
		fieldValues := make(map[string]interface{}, len(valueKeys))
		for _, fieldKey := range valueKeys {
			fieldValue := valuesObj.Get(fieldKey)
			if !IsValid(fieldValue) {
				continue
			}

			if IsBuffer(fieldValue) {
				bytes, err := BufferToBytes(fieldValue)
				if err != nil {
					return nil, fmt.Errorf("数组内的元素 [%s] 的 values 字段的 [%s] 字段不是有效的 Buffer", key, fieldKey)
				}
				fieldValues[fieldKey] = bytes
			} else {
				fieldValues[fieldKey] = fieldValue.Export()
			}
		}

		results = append(results, ParseResult{
			ID:     idStr,
			Table:  tableId,
			CID:    cidValue,
			Time:   timeValue,
			Values: fieldValues,
		})
	}

	return results, nil
}
