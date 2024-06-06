package gojs

import (
	"encoding/binary"
	"fmt"
	"github.com/dop251/goja"
)

func checkBufferOffsetAndLength(buf []byte, offset, length int) error {
	if offset < 0 {
		return fmt.Errorf("the offset cannot be negative")
	} else if length <= 0 {
		return fmt.Errorf("the length cannot be negative")
	}

	if len(buf) < offset+length {
		return fmt.Errorf("the offset + length = %d is out of range buffer length %d", offset+length, len(buf))
	}

	return nil
}

func convertOffset(offsetValue goja.Value) (int, error) {
	var offset int
	switch v := offsetValue.Export().(type) {
	case int:
		offset = v
	case int32:
		offset = int(v)
	case int64:
		offset = int(v)
	default:
		return -1, fmt.Errorf("invalid offset '%v'", v)
	}

	if offset < 0 {
		return -1, fmt.Errorf("invalid offset %d, cannot be negative", offset)
	}

	return offset, nil
}

func convertToUint64(value goja.Value) (uint64, error) {
	var val uint64

	switch v := value.Export().(type) {
	case int:
		if v < 0 {
			return 0, fmt.Errorf("invalid value %d, must be negative", v)
		}
		val = uint64(v)
	case int32:
		if v < 0 {
			return 0, fmt.Errorf("invalid value %d, must be negative", v)
		}
		val = uint64(v)
	case int64:
		if v < 0 {
			return 0, fmt.Errorf("invalid value %d, must be negative", v)
		}
		val = uint64(v)
	case uint64:
		val = v
	default:
		return 0, fmt.Errorf("invalid value '%v'", v)
	}

	return val, nil
}

func readBigInt64LE(vm *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			vm.Interrupt(fmt.Errorf("the offset is not specified"))
			return goja.Undefined()
		}

		offset, err := convertOffset(call.Arguments[0])
		if err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		buffer, _ := BufferToBytes(call.This)
		if err := checkBufferOffsetAndLength(buffer, offset, 8); err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		value := binary.LittleEndian.Uint64(buffer[offset:])
		return vm.ToValue(int64(value))
	}
}

func writeBigInt64LE(vm *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			vm.Interrupt(fmt.Errorf("the value or offset is not specified"))
			return goja.Undefined()
		}

		var value int64
		var offset int

		switch v := call.Arguments[0].Export().(type) {
		case int:
			value = int64(v)
		case int32:
			value = int64(v)
		case int64:
			value = v
		default:
			vm.Interrupt(fmt.Errorf("invalid value '%v'", v))
			return goja.Undefined()
		}

		offset, err := convertOffset(call.Arguments[1])
		if err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		buffer, _ := BufferToBytes(call.This)
		if err := checkBufferOffsetAndLength(buffer, offset, 8); err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		binary.LittleEndian.PutUint64(buffer[offset:], uint64(value))

		return goja.Undefined()
	}
}

func readBigInt64BE(vm *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			vm.Interrupt(fmt.Errorf("the offset is not specified"))
			return goja.Undefined()
		}

		offset, err := convertOffset(call.Arguments[0])
		if err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		buffer, _ := BufferToBytes(call.This)
		if err := checkBufferOffsetAndLength(buffer, offset, 8); err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		value := binary.BigEndian.Uint64(buffer[offset:])
		return vm.ToValue(int64(value))
	}
}

func writeBigInt64BE(vm *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			vm.Interrupt(fmt.Errorf("the value or offset is not specified"))
			return goja.Undefined()
		}

		var value int64
		var offset int

		switch v := call.Arguments[0].Export().(type) {
		case int:
			value = int64(v)
		case int32:
			value = int64(v)
		case int64:
			value = v
		default:
			vm.Interrupt(fmt.Errorf("invalid value '%v'", v))
			return goja.Undefined()
		}

		offset, err := convertOffset(call.Arguments[1])
		if err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		buffer, _ := BufferToBytes(call.This)
		if err := checkBufferOffsetAndLength(buffer, offset, 8); err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		binary.BigEndian.PutUint64(buffer[offset:], uint64(value))

		return goja.Undefined()
	}
}

func readBigUInt64LE(vm *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			vm.Interrupt(fmt.Errorf("the offset is not specified"))
			return goja.Undefined()
		}

		offset, err := convertOffset(call.Arguments[0])
		if err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		buffer, _ := BufferToBytes(call.This)
		if err := checkBufferOffsetAndLength(buffer, offset, 8); err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		value := binary.LittleEndian.Uint64(buffer[offset:])
		return vm.ToValue(value)
	}
}

func writeBigUInt64LE(vm *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			vm.Interrupt(fmt.Errorf("the value or offset is not specified"))
			return goja.Undefined()
		}

		value, err := convertToUint64(call.Arguments[0])
		if err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		offset, err := convertOffset(call.Arguments[1])
		if err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		buffer, _ := BufferToBytes(call.This)
		if err := checkBufferOffsetAndLength(buffer, offset, 8); err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		binary.LittleEndian.PutUint64(buffer[offset:], value)

		return goja.Undefined()
	}
}

func readBigUInt64BE(vm *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			vm.Interrupt(fmt.Errorf("the offset is not specified"))
			return goja.Undefined()
		}

		offset, err := convertOffset(call.Arguments[0])
		if err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		buffer, _ := BufferToBytes(call.This)
		if err := checkBufferOffsetAndLength(buffer, offset, 8); err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		value := binary.BigEndian.Uint64(buffer[offset:])
		return vm.ToValue(value)
	}
}

func writeBigUInt64BE(vm *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			vm.Interrupt(fmt.Errorf("the value or offset is not specified"))
			return goja.Undefined()
		}

		value, err := convertToUint64(call.Arguments[0])
		if err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		offset, err := convertOffset(call.Arguments[1])
		if err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		buffer, _ := BufferToBytes(call.This)
		if err := checkBufferOffsetAndLength(buffer, offset, 8); err != nil {
			vm.Interrupt(err)
			return goja.Undefined()
		}

		binary.BigEndian.PutUint64(buffer[offset:], value)

		return goja.Undefined()
	}
}
