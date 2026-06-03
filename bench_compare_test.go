package bench_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	shamaton "github.com/shamaton/msgpack/v3"
	. "github.com/shamaton/msgpack_bench"
	"github.com/shamaton/msgpack_bench/msgpackgen"
	"github.com/shamaton/msgpack_bench/protocmp"
	"github.com/shamaton/zeroformatter"
	"github.com/ugorji/go/codec"
	vmihailenco "github.com/vmihailenco/msgpack/v5"
	"google.golang.org/protobuf/proto"
)

func init() {
	initCompare()
}

var bench = BenchMarkStruct{
	Int:    -123,
	Uint:   456,
	Float:  1.234,
	Double: 6.789,
	Bool:   true,
	String: "this is text.",
	Array:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
	Map:    map[string]uint{"this": 1, "is": 2, "map": 3},
	Child:  BenchChild{Int: 123456, String: "this is struct of child"},
}

var protobench = &protocmp.BenchMarkStruct{
	Int:     int32(bench.Int),
	Uint:    uint32(bench.Uint),
	Float:   bench.Float,
	Double:  bench.Double,
	Bool:    bench.Bool,
	String_: bench.String,
	Array:   []int32{1, 2, 3, 4, 5, 6, 7, 8, 9},
	Map:     map[string]uint32{"this": 1, "is": 2, "map": 3},
	Child:   &protocmp.BenchChild{Int: 123456, String_: "this is struct of child"},
}

var (
	arrayMsgpackBench []byte
	mapMsgpackBench   []byte
	zeroFmtpackBench  []byte
	jsonPackBench     []byte
	gobPackBench      []byte
	protoPackBench    []byte
)

// for codec
var (
	mhBench = &codec.MsgpackHandle{}
)

func initCompare() {
	// ugorji
	//mhBench.MapType = reflect.TypeOf(bench)

	d, err := shamaton.MarshalAsArray(bench)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	arrayMsgpackBench = d
	d, err = shamaton.MarshalAsMap(bench)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	mapMsgpackBench = d

	d, err = zeroformatter.Serialize(bench)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	zeroFmtpackBench = d

	d, err = json.Marshal(bench)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	jsonPackBench = d

	d, err = proto.Marshal(protobench)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	protoPackBench = d

	buf := bytes.NewBuffer(nil)
	err = gob.NewEncoder(buf).Encode(bench)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	gobPackBench = buf.Bytes()

	// check
	check()
	checkCompareEncodeOutputs()
}

func check() {
	var mp, arr, genmp, genarr, vmp, varr, tmp, c, zero, jsn, gb BenchMarkStruct
	mustCompareCheck("shamaton array", shamaton.UnmarshalAsArray(arrayMsgpackBench, &arr))
	mustCompareCheck("shamaton map", shamaton.UnmarshalAsMap(mapMsgpackBench, &mp))
	mustCompareCheck("msgpackgen array", msgpackgen.UnmarshalAsArray(arrayMsgpackBench, &genarr))
	mustCompareCheck("msgpackgen map", msgpackgen.UnmarshalAsMap(mapMsgpackBench, &genmp))
	mustCompareCheck("vmihailenco array", vmihailenco.Unmarshal(arrayMsgpackBench, &varr))
	mustCompareCheck("vmihailenco map", vmihailenco.Unmarshal(mapMsgpackBench, &vmp))
	_, err := tmp.UnmarshalMsg(mapMsgpackBench)
	mustCompareCheck("tinylib map", err)
	mustCompareCheck("ugorji map", codec.NewDecoderBytes(mapMsgpackBench, mhBench).Decode(&c))
	mustCompareCheck("zeroformatter", zeroformatter.Deserialize(&zero, zeroFmtpackBench))
	mustCompareCheck("json", json.Unmarshal(jsonPackBench, &jsn))
	mustCompareCheck("gob", gob.NewDecoder(bytes.NewBuffer(gobPackBench)).Decode(&gb))
	var pb protocmp.BenchMarkStruct
	mustCompareCheck("proto", proto.Unmarshal(protoPackBench, &pb))

	if !reflect.DeepEqual(mp, arr) {
		fmt.Println("not equal: mp vs arr")
		fmt.Printf("mp: %+v\n", mp)
		fmt.Printf("arr: %+v\n", arr)
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, genarr) {
		fmt.Println("not equal: mp vs genarr")
		fmt.Printf("mp: %+v\n", mp)
		fmt.Printf("genarr: %+v\n", genarr)
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, genmp) {
		fmt.Println("not equal: mp vs genmp")
		fmt.Printf("mp: %+v\n", mp)
		fmt.Printf("genmp: %+v\n", genmp)
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, varr) {
		fmt.Println("not equal: mp vs varr")
		fmt.Printf("mp: %+v\n", mp)
		fmt.Printf("varr: %+v\n", varr)
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, vmp) {
		fmt.Println("not equal: mp vs vmp")
		fmt.Printf("mp: %+v\n", mp)
		fmt.Printf("vmp: %+v\n", vmp)
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, tmp) {
		fmt.Println("not equal: mp vs tmp")
		fmt.Printf("mp: %+v\n", mp)
		fmt.Printf("tmp: %+v\n", tmp)
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, c) {
		fmt.Println("not equal: mp vs c")
		fmt.Printf("mp: %+v\n", mp)
		fmt.Printf("c: %+v\n", c)
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, zero) {
		fmt.Println("not equal: mp vs zero")
		fmt.Printf("mp: %+v\n", mp)
		fmt.Printf("zero: %+v\n", zero)
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, jsn) {
		fmt.Println("not equal: mp vs json")
		fmt.Printf("mp: %+v\n", mp)
		fmt.Printf("json: %+v\n", jsn)
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, gb) {
		fmt.Println("not equal: mp vs gob")
		fmt.Printf("mp: %+v\n", mp)
		fmt.Printf("gob: %+v\n", gb)
		os.Exit(1)
	}
	if got := benchFromProto(&pb); !reflect.DeepEqual(mp, got) {
		fmt.Println("not equal: mp vs proto")
		fmt.Printf("mp: %+v\n", mp)
		fmt.Printf("proto: %+v\n", got)
		os.Exit(1)
	}
}

func mustCompareCheck(name string, err error) {
	if err != nil {
		fmt.Println("init err: ", name, err)
		os.Exit(1)
	}
}

func mustCompareValue(name string, got BenchMarkStruct) {
	if !reflect.DeepEqual(bench, got) {
		fmt.Println("not equal: bench vs ", name)
		fmt.Printf("bench: %+v\n", bench)
		fmt.Printf("%s: %+v\n", name, got)
		os.Exit(1)
	}
}

func benchFromProto(p *protocmp.BenchMarkStruct) BenchMarkStruct {
	array := make([]int, len(p.Array))
	for i, v := range p.Array {
		array[i] = int(v)
	}
	m := make(map[string]uint, len(p.Map))
	for k, v := range p.Map {
		m[k] = uint(v)
	}
	var child BenchChild
	if p.Child != nil {
		child = BenchChild{
			Int:    int(p.Child.Int),
			String: p.Child.String_,
		}
	}
	return BenchMarkStruct{
		Int:    int(p.Int),
		Uint:   uint(p.Uint),
		Float:  p.Float,
		Double: p.Double,
		Bool:   p.Bool,
		String: p.String_,
		Array:  array,
		Map:    m,
		Child:  child,
	}
}

func checkCompareEncodeOutputs() {
	var v BenchMarkStruct

	d, err := msgpackgen.MarshalAsArray(&bench)
	mustCompareCheck("msgpackgen encode array", err)
	mustCompareCheck("msgpackgen encoded array decode", msgpackgen.UnmarshalAsArray(d, &v))
	mustCompareValue("msgpackgen array", v)

	v = BenchMarkStruct{}
	d, err = msgpackgen.MarshalAsMap(&bench)
	mustCompareCheck("msgpackgen encode map", err)
	mustCompareCheck("msgpackgen encoded map decode", msgpackgen.UnmarshalAsMap(d, &v))
	mustCompareValue("msgpackgen map", v)

	v = BenchMarkStruct{}
	d, err = shamaton.MarshalAsArray(bench)
	mustCompareCheck("shamaton encode array", err)
	mustCompareCheck("shamaton encoded array decode", shamaton.UnmarshalAsArray(d, &v))
	mustCompareValue("shamaton array", v)

	v = BenchMarkStruct{}
	d, err = shamaton.MarshalAsMap(bench)
	mustCompareCheck("shamaton encode map", err)
	mustCompareCheck("shamaton encoded map decode", shamaton.UnmarshalAsMap(d, &v))
	mustCompareValue("shamaton map", v)

	v = BenchMarkStruct{}
	d, err = bench.MarshalMsg(nil)
	mustCompareCheck("tinylib encode", err)
	_, err = v.UnmarshalMsg(d)
	mustCompareCheck("tinylib encoded decode", err)
	mustCompareValue("tinylib", v)

	v = BenchMarkStruct{}
	buf := []byte{}
	mustCompareCheck("ugorji encode", codec.NewEncoderBytes(&buf, mhBench).Encode(bench))
	mustCompareCheck("ugorji encoded decode", codec.NewDecoderBytes(buf, mhBench).Decode(&v))
	mustCompareValue("ugorji", v)

	v = BenchMarkStruct{}
	var bytesBuf bytes.Buffer
	enc := vmihailenco.NewEncoder(&bytesBuf)
	enc.UseArrayEncodedStructs(true)
	mustCompareCheck("vmihailenco encode array", enc.Encode(bench))
	mustCompareCheck("vmihailenco encoded array decode", vmihailenco.Unmarshal(bytesBuf.Bytes(), &v))
	mustCompareValue("vmihailenco array", v)

	v = BenchMarkStruct{}
	d, err = vmihailenco.Marshal(bench)
	mustCompareCheck("vmihailenco encode map", err)
	mustCompareCheck("vmihailenco encoded map decode", vmihailenco.Unmarshal(d, &v))
	mustCompareValue("vmihailenco map", v)

	var pb protocmp.BenchMarkStruct
	d, err = proto.Marshal(protobench)
	mustCompareCheck("proto encode", err)
	mustCompareCheck("proto encoded decode", proto.Unmarshal(d, &pb))
	mustCompareValue("proto", benchFromProto(&pb))

	v = BenchMarkStruct{}
	d, err = json.Marshal(bench)
	mustCompareCheck("json encode", err)
	mustCompareCheck("json encoded decode", json.Unmarshal(d, &v))
	mustCompareValue("json", v)

	v = BenchMarkStruct{}
	bytesBuf.Reset()
	mustCompareCheck("gob encode", gob.NewEncoder(&bytesBuf).Encode(bench))
	mustCompareCheck("gob encoded decode", gob.NewDecoder(bytes.NewBuffer(bytesBuf.Bytes())).Decode(&v))
	mustCompareValue("gob", v)

	v = BenchMarkStruct{}
	d, err = zeroformatter.Serialize(bench)
	mustCompareCheck("zeroformatter encode", err)
	mustCompareCheck("zeroformatter encoded decode", zeroformatter.Deserialize(&v, d))
	mustCompareValue("zeroformatter", v)
}

func BenchmarkCompareDecodeShamatonGenArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := msgpackgen.UnmarshalAsArray(arrayMsgpackBench, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareDecodeTinylib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		_, err := r.UnmarshalMsg(mapMsgpackBench)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareDecodeShamatonGenMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := msgpackgen.UnmarshalAsMap(mapMsgpackBench, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareDecodeShamatonGen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := msgpackgen.Unmarshal(mapMsgpackBench, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareDecodeShamatonArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := shamaton.UnmarshalAsArray(arrayMsgpackBench, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareDecodeShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := shamaton.UnmarshalAsMap(mapMsgpackBench, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareDecodeUgorji(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		dec := codec.NewDecoderBytes(mapMsgpackBench, mhBench)
		err := dec.Decode(&r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareDecodeVmihailencoArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := vmihailenco.Unmarshal(arrayMsgpackBench, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareDecodeVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := vmihailenco.Unmarshal(mapMsgpackBench, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareDecodeProtocolBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r protocmp.BenchMarkStruct
		err := proto.Unmarshal(protoPackBench, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareDecodeJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := json.Unmarshal(jsonPackBench, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareDecodeGob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		buf := bytes.NewBuffer(gobPackBench)
		err := gob.NewDecoder(buf).Decode(&r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareDecodeZeroformatter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := zeroformatter.Deserialize(&r, zeroFmtpackBench)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/////////////////////////////////////////////////////////////////

func BenchmarkCompareEncodeShamatonGenArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := msgpackgen.MarshalAsArray(&bench)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareEncodeTinylib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := bench.MarshalMsg(nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareEncodeShamatonGenMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := msgpackgen.MarshalAsMap(&bench)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareEncodeShamatonGen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := msgpackgen.Marshal(&bench)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareEncodeShamatonArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.MarshalAsArray(bench)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareEncodeShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.MarshalAsMap(bench)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareEncodeUgorji(b *testing.B) {
	for i := 0; i < b.N; i++ {

		buf := []byte{}
		enc := codec.NewEncoderBytes(&buf, mhBench)
		err := enc.Encode(bench)

		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareEncodeVmihailencoArray(b *testing.B) {
	for i := 0; i < b.N; i++ {

		var buf bytes.Buffer
		enc := vmihailenco.NewEncoder(&buf)
		enc.UseArrayEncodedStructs(true)
		err := enc.Encode(bench)

		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareEncodeVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(bench)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareEncodeProtocolBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(protobench)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareEncodeJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(bench)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareEncodeGob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(nil)
		err := gob.NewEncoder(buf).Encode(bench)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompareEncodeZeroformatter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := zeroformatter.Serialize(bench)
		if err != nil {
			b.Fatal(err)
		}
	}
}
