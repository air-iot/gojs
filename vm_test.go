package gojs

import (
	"encoding/json"
	"github.com/air-iot/errors"
	"github.com/dop251/goja"
	"math"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestRun(t *testing.T) {
	type args struct {
		script string
		values []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				script: `
	function handler(j) {
	console.log(cryptoJs.SHA256("Info"))
	console.log(xmlJs)
	console.log(new Buffer("1123").toString("hex"));
	console.log(j);
	console.log(moment());
	return j;
}`,
				values: []interface{}{1},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "test_formulajs",
			args: args{
				script: `
	function handler(j) {
	console.log(formulajs)
	console.log(formulajs.SUM([1, 2, 3]))
	return formulajs.SUM([1, 2, 3]);
}`,
				values: []interface{}{1},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Run(tt.args.script, tt.args.values...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got.Export())
			//if !reflect.DeepEqual(got.Export(), tt.want) {
			//	t.Errorf("Run() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func Benchmark_Run1(b *testing.B) {
	var f = func(j int) float64 {
		return math.Max(10, float64(j))
	}
	for i := 0; i < b.N; i++ {
		f(i)
	}
}

func Benchmark_Run(b *testing.B) {
	for i := 0; i < b.N; i++ {
		js := `function handler(j) {
  return _.max([1,20,j]);
}`
		res, err := Run(js, i)
		if err != nil {
			b.Fatal(err)
		}
		_ = res
		//b.Log(1, res.Export())
	}
}

func TestRunByIdAndScript(t *testing.T) {
	type args struct {
		id     string
		script string
		values []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				id: "test1",
				script: `
	function handler(j) {
	console.log(cryptoJs.SHA256("Info"))
	console.log(xmlJs)
	console.log(new Buffer("1123").toString("hex"));
	console.log(j);
	console.log(moment());
	return j;
}`,
				values: []interface{}{1},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RunByIdAndScript(tt.args.id, tt.args.script, tt.args.values...)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunByIdAndScript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got.Export())
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("RunByIdAndScript() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func Test_RunByIdAndScript(t *testing.T) {
	js := `function handler(j) {
	_state["a"] = 1
  return _.max([1,20,j]);
}`
	res, err := RunByIdAndScript("1", js, 1)
	if err != nil {
		t.Fatal(err)
	}
	_ = res
	t.Log(1, res.Export())
	js = `function handler(j) {
	console.log(_state["a"])
  return _.max([1,3,20,j]);
}`
	res, err = RunByIdAndScript("1", js, 1)
	if err != nil {
		t.Fatal(err)
	}
	_ = res
	t.Log(2, res.Export())

	res, err = RunByIdAndScript("1", js, 2)
	if err != nil {
		t.Fatal(err)
	}
	_ = res
	t.Log(3, res.Export())
}

func Benchmark_RunByIdAndScript(b *testing.B) {
	for i := 0; i < b.N; i++ {
		js := `function handler(j) {
  return _.max([1,20,j]);
}`

		res, err := RunByIdAndScript("1", js, i)
		if err != nil {
			b.Fatal(err)
		}
		_ = res
		//b.Log(1, res.Export())
	}
}

func TestGojaBuffer(t *testing.T) {
	js1 := `function handler() {
		return {"messageType":1,"messageType1":"a","data":new Uint8Array(Buffer.from('{"type":"query","data":{"all":true}}')).buffer};
	}`
	//type Val struct {
	//	MessageType int         `json:"messageType"`
	//	Result        interface{} `json:"data"`
	//}
	val, err := Run(js1)
	if err != nil {
		t.Fatal(err)
	}
	val1 := val.(*goja.Object)
	messageType := val1.Get("messageType")
	t.Log(reflect.TypeOf(messageType), messageType.ToInteger())
	messageType1 := val1.Get("messageType1")
	t.Log(reflect.TypeOf(messageType1), messageType1.String())
	data := val1.Get("data").Export()
	t.Log(reflect.TypeOf(data), data)

	//b, err := val1.MarshalJSON()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//var tmp Val
	//if err := json.Unmarshal(b, &tmp); err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(tmp.Result)
}

func TestGojaBuffer1(t *testing.T) {
	//js1 := `function handler() {
	//	return new Uint8Array(Buffer.from('{"type":"query","data":{"all":true}}')).buffer;
	//}`

	js2 := `function handler(val) {
		console.log(val.toString("hex"))
		return val;
	}`
	js1 := `function handler() {
		return Buffer.from('{"type":"query","data":{"all":true}}');
	}`

	//type Val struct {
	//	MessageType int         `json:"messageType"`
	//	Result        interface{} `json:"data"`
	//}
	_, _ = GetJsVm("test1", js1)
	val, err := RunById("test1")
	if err != nil {
		t.Fatal(err)
	}
	bs, err := BufferToBytes(val)
	if err != nil {
		t.Fatal(err)
	}

	vm2, _ := GetJsVm("test2", js2)
	bsVal, err := BytesToBuffer(vm2.VM, bs)
	if err != nil {
		t.Fatal(err)
	}

	_, err = RunById("test2", bsVal)
	if err != nil {
		t.Fatal(err)
	}
	//val1 := val.Export().(goja.ArrayBuffer)
	//t.Log(reflect.TypeOf(val1), val1.Bytes(), string(val1.Bytes()))
}

func TestCo(t *testing.T) {
	js1 := `function handler() {
		console.log(1)
		return Buffer.from('{"type":"query","data":{"all":true}}');
	}`
	RunByIdAndScript("test1", js1)
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			RunByIdAndScript("test1", js1)
		}()
	}
	wg.Wait()
}

func TestIconv(t *testing.T) {
	got, err := RunByIdAndScript("iconv", `function handler() {
		const buffer = iconv.encode("123.ntp", "UTF-16LE");
		// 49,0,50,0,51,0,46,0,110,0,116,0,112,0
		console.log("iconv.encode:", new Uint8Array(buffer));
		return iconv.decode(buffer, "UTF-16LE");
	}`)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Compare(got.String(), "123.ntp") != 0 {
		t.Errorf("iconv.decode, expected: 123.ntp, got: %s", got.String())
	}
}

func Test_1(t *testing.T) {
	js1 := `function handler() {
		console.log(1)
		return {"type":"query","data":{"all":true}};
	}`
	val, err := Run(js1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reflect.TypeOf(val.Export()))
}

func Test_2(t *testing.T) {
	js1 := `function handler() {
		console.log(new Date())
		apilib.SleepMill(10000)
		console.log(new Date())
		return {};
	}`
	val, err := Run(js1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)
}

func Test_Buffer(t *testing.T) {
	js1 := `function handler() {
const buf = Buffer.from([0x62, 0x75, 0x66, 0x66, 0x65, 0x72]);
console.log(buf.toString()); // 'buffer'
		return {};
	}`
	val, err := Run(js1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)
}

func Test_xmljs(t *testing.T) {
	js1 := `function handler() {
let xml =
'<?xml version="1.0" encoding="utf-8"?>' +
'<note importance="high" logged="true">' +
'    <title>Happy</title>' +
'    <todo>Work</todo>' +
'    <todo>Play</todo>' +
'</note>';
console.log(xmlJs.xml2json(xml))
		return {};
	}`
	val, err := Run(js1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)
}

func Test_formulajs(t *testing.T) {
	js1 := `function handler() {
console.log(formulajs.SUM([1, 2, 3]))
		return {};
	}`
	val, err := Run(js1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)
}

func Test_lodash(t *testing.T) {
	js1 := `function handler() {
let result = _.max([1,20])
console.log(result)
		return {};
	}`
	val, err := Run(js1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)
}

func Test_reduce(t *testing.T) {
	js1 := `function handler(){
let arr = [{"name":"temperature","value":26.3},{"name":"humidity","value":65}];

let obj = arr.reduce((acc, cur) => {
  acc[cur.name] = cur.value;
  return acc;
}, {});
return obj;
}`
	val, err := Run(js1)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(val)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}

func Test_err(t *testing.T) {
	err := HandlerError
	t.Log(errors.Is(err, HandlerError))
}

func Test_Crc(t *testing.T) {
	js1 := `function handler(){
	const buf = Buffer.from([0x01, 0x03, 0x00, 0x12, 0x00, 0x10, 0xE4, 0x03]);
	const got = crc.checksumModbus(buf.slice(0, buf.length-2));
	const expected = buf.slice(buf.length - 2).readUInt16LE();
console.log(got, expected)
return got === expected;
}`
	val, err := Run(js1)
	if err != nil {
		t.Fatal(err)
	}

	if !val.Export().(bool) {
		t.Fatal("crc error, expected: true, got: ", val.Export().(bool))
	}
}

func Test_throw(t *testing.T) {
	js1 := `function handler() {
  try {
    let a = JSON.parse('{"a":1}')
    if (a.a === 1) {
      throw new Error(1)
    }
    return a.a * 2;
  } catch (err) {
    return 1
  }
}
`
	val, err := Run(js1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val.Export())
}
