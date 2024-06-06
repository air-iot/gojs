package gojs

import (
	"crypto/md5"
	"fmt"
	"github.com/air-iot/errors"
	"github.com/air-iot/gojs/api"
	log2 "github.com/air-iot/gojs/log"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"github.com/patrickmn/go-cache"
)

var localCache = cache.New(5*time.Minute, 10*time.Minute)
var registry *require.Registry
var programs []*goja.Program
var apilib *api.Lib
var logLib *log2.Log

type Callback func(vm *goja.Runtime) error

func init() {
	registry = require.NewRegistry()
	programs = make([]*goja.Program, 0)
	initPackages("packages/bcd.js")
	initPackages("packages/buffer.js")
	initPackages("packages/lodash.js")
	initPackages("packages/crypto-js.js")
	initPackages("packages/moment.js")
	initPackages("packages/xml-js.js")
	initPackages("packages/formulajs.js")
	initPackages("packages/iconv-lite.js")
	initPackages("packages/forge.js")
	//initPackages("packages/uuid.js")
	apilib = api.NewLib()
	logLib = log2.NewLogger()
}

func initPackages(packagePath string) {
	lodashBytes, err := F.ReadFile(packagePath)
	if err != nil {
		panic(fmt.Errorf("read %s err,%s", packagePath, err))
	}
	p, err := goja.Compile(packagePath, string(lodashBytes), false)
	if err != nil {
		panic(fmt.Errorf("compile %s err,%s", packagePath, err))
	}
	programs = append(programs, p)
}

var HandlerError = errors.New400Response(100040004, "脚本函数handler未找到")

type JSvm struct {
	lock    sync.Mutex
	VM      *goja.Runtime
	Handler goja.Callable
	Script  string
}

func NewJsVm(id, script string) (*JSvm, error) {
	jsVM, err := GetJsVm(id, script)
	if err != nil {
		return nil, err
	}
	return jsVM, nil
}

func (j *JSvm) SetObj(key string, obj interface{}) error {
	if err := j.VM.Set(key, obj); err != nil {
		return errors.Wrap400Err(err, 100040001)
	}
	return nil
}

func GetVm() (*goja.Runtime, error) {
	vm := goja.New()
	//registry := new(require.Registry) // this can be shared by multiple runtimes
	registry.Enable(vm)
	console.Enable(vm)
	obj := vm.GlobalObject()
	state := map[string]interface{}{}
	if err := obj.Set("_state", state); err != nil {
		return nil, errors.Wrap400Err(err, 100040001)
	}
	for _, program := range programs {
		_, err := vm.RunProgram(program)
		if err != nil {
			return nil, errors.Wrap400Err(err, 100040002)
		}
	}
	_ = vm.Set("_", vm.Get("lodash"))
	_ = vm.Set("CryptoJS", vm.Get("cryptoJs"))
	_ = vm.Set("Buffer", vm.Get("Buffer").(*goja.Object).Get("Buffer"))
	_ = vm.Set("formulajs", vm.Get("formulajsformulajs"))
	_ = vm.Set("iconv", vm.Get("iconvLite"))
	_ = vm.Set("apilib", apilib)
	_ = vm.Set(log2.Key, logLib)
	_ = AttachCrc(vm)
	return vm, nil
}

func GetVmCallback(cb Callback) (*goja.Runtime, error) {
	vm, err := GetVm()
	if err != nil {
		return nil, err
	}
	_ = AttachCrc(vm)
	if cb != nil {
		if err := cb(vm); err != nil {
			return nil, err
		}
	}
	return vm, nil
}

func GetJsVm(id, script string) (*JSvm, error) {
	jsVMI, ok := localCache.Get(id)
	var jsVM *JSvm
	if !ok {
		vm, err := GetVm()
		if err != nil {
			return nil, err
		}
		if _, err := vm.RunString(script); err != nil {
			return nil, errors.Wrap400Err(err, 100040003)
		}
		handler, ok := goja.AssertFunction(vm.Get("handler"))
		if !ok {
			return nil, HandlerError
		}
		jsVM = &JSvm{
			VM:      vm,
			Handler: handler,
			Script:  script,
		}
	} else {
		jsVM, _ = jsVMI.(*JSvm)
		if fmt.Sprintf("%x", md5.Sum([]byte(script))) != fmt.Sprintf("%x", md5.Sum([]byte(jsVM.Script))) {
			if _, err := jsVM.VM.RunString(script); err != nil {
				return nil, errors.Wrap400Err(err, 100040003)
			}
			handler, ok := goja.AssertFunction(jsVM.VM.Get("handler"))
			if !ok {
				return nil, HandlerError
			}
			jsVM.Script = script
			jsVM.Handler = handler
		}
	}
	localCache.Set(id, jsVM, cache.DefaultExpiration)
	return jsVM, nil
}

func Run(script string, values ...interface{}) (goja.Value, error) {
	id := fmt.Sprintf("%x", md5.Sum([]byte(script)))
	return RunByIdAndScript(id, script, values...)
}

func RunByIdAndScript(id, script string, values ...interface{}) (goja.Value, error) {
	jsVM, err := GetJsVm(id, script)
	if err != nil {
		return nil, err
	}
	vals := make([]goja.Value, len(values))
	if values != nil {
		for i, v := range values {
			gojaVal, ok := v.(goja.Value)
			if ok {
				vals[i] = gojaVal
			} else {
				vals[i] = jsVM.VM.ToValue(v)
			}
		}
	}
	jsVM.lock.Lock()
	defer jsVM.lock.Unlock()
	output, err := jsVM.Handler(goja.Undefined(), vals...)
	if err != nil {
		return nil, errors.Wrap400Err(err, 100040005)
	}
	return output, nil
}

func RunById(id string, values ...interface{}) (goja.Value, error) {
	jsVMI, ok := localCache.Get(id)
	if !ok {
		return nil, errors.New400Response(100040006, "未找到vm")
	}
	jsVM, _ := jsVMI.(*JSvm)
	localCache.Set(id, jsVM, cache.DefaultExpiration)
	vals := make([]goja.Value, len(values))
	if values != nil {
		for i, v := range values {
			gojaVal, ok := v.(goja.Value)
			if ok {
				vals[i] = gojaVal
			} else {
				vals[i] = jsVM.VM.ToValue(v)
			}
		}
	}
	jsVM.lock.Lock()
	defer jsVM.lock.Unlock()
	output, err := jsVM.Handler(goja.Undefined(), vals...)
	if err != nil {
		return nil, errors.Wrap400Err(err, 100040005)
	}
	return output, nil
}

func BufferToBytes(bufferVal goja.Value) ([]byte, error) {
	obj, ok := bufferVal.(*goja.Object)
	if !ok {
		return nil, errors.New400Response(100040007, "参数格式非Object")
	}
	buffer := obj.Get("buffer")
	if buffer == nil {
		return nil, errors.New400Response(100040008, "未找到buffer")
	}
	arrayBuffer, ok := buffer.Export().(goja.ArrayBuffer)
	if !ok {
		return nil, errors.New400Response(100040009, "参数格式非arrayBuffer")
	}

	dataBytes := arrayBuffer.Bytes()
	dataLength := int64(len(dataBytes))

	var offset, length int64 = -1, -1
	if v := obj.Get("offset"); v != nil {
		offset = v.ToInteger()
	}
	if v := obj.Get("length"); v != nil {
		length = v.ToInteger()
	}

	if offset < 0 || length < 0 || offset >= dataLength || offset+length > dataLength {
		return nil, errors.New400ErrResponse(100040010, map[string]interface{}{"offset": offset, "length": length}, "offset %d 或 length %d 超出 buffer 大小", offset, length)
	}

	return dataBytes[offset : offset+length], nil
}

func BytesToBuffer(vm *goja.Runtime, bs []byte) (goja.Value, error) {
	buffer := vm.Get("Buffer")
	if buffer == nil {
		return nil, errors.New400Response(100040011, "未加载Buffer")
	}
	bufferObj := buffer.ToObject(vm)
	from := bufferObj.Get("from")
	if from == nil {
		return nil, errors.New400Response(100040012, "未加载Buffer from方法")
	}
	fromFn, ok := from.Export().(func(goja.FunctionCall) goja.Value)
	if !ok {
		return nil, errors.New400Response(100040013, "Buffer from方法类型非FunctionCall")
	}
	buf := fromFn(goja.FunctionCall{
		This:      bufferObj,
		Arguments: []goja.Value{vm.ToValue(bs)},
	})
	return buf, nil
}
