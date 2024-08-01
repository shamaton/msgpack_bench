package bench

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	shamaton "github.com/shamaton/msgpack/v2"
	"github.com/shamaton/msgpack_bench/protocmp"
	shamatongen "github.com/shamaton/msgpackgen/msgpack"
	"github.com/shamaton/zeroformatter"
	"github.com/ugorji/go/codec"
	vmihailenco "github.com/vmihailenco/msgpack/v5"
	"google.golang.org/protobuf/proto"
)

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
}

func check() {
	var mp, arr, genmp, genarr, vmp, varr, tmp, c BenchMarkStruct
	shamaton.UnmarshalAsArray(arrayMsgpackBench, &arr)
	shamaton.UnmarshalAsMap(mapMsgpackBench, &mp)
	shamatongen.UnmarshalAsArray(arrayMsgpackBench, &genarr)
	shamatongen.UnmarshalAsMap(mapMsgpackBench, &genmp)
	vmihailenco.Unmarshal(arrayMsgpackBench, &varr)
	vmihailenco.Unmarshal(mapMsgpackBench, &vmp)
	tmp.UnmarshalMsg(mapMsgpackBench)
	codec.NewDecoderBytes(mapMsgpackBench, mhBench).Decode(&c)

	if !reflect.DeepEqual(mp, arr) {
		fmt.Println("not equal")
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, genarr) {
		fmt.Println("not equal")
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, genmp) {
		fmt.Println("not equal")
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, varr) {
		fmt.Println("not equal")
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, vmp) {
		fmt.Println("not equal")
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, tmp) {
		fmt.Println("not equal")
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, c) {
		fmt.Println("not equal")
		os.Exit(1)
	}
}

func BenchmarkCompareDecodeArrayShamatonGen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := shamatongen.UnmarshalAsArray(arrayMsgpackBench, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareDecodeTinylib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		_, err := r.UnmarshalMsg(mapMsgpackBench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareDecodeShamatonGen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := shamatongen.UnmarshalAsMap(mapMsgpackBench, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareDecodeArrayShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := shamaton.UnmarshalAsArray(arrayMsgpackBench, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareDecodeShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := shamaton.UnmarshalAsMap(mapMsgpackBench, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareDecodeUgorji(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		dec := codec.NewDecoderBytes(mapMsgpackBench, mhBench)
		err := dec.Decode(&r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareDecodeArrayVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := vmihailenco.Unmarshal(arrayMsgpackBench, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareDecodeVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := vmihailenco.Unmarshal(mapMsgpackBench, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareDecodeProtocolBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r protocmp.BenchMarkStruct
		err := proto.Unmarshal(protoPackBench, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareDecodeJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := json.Unmarshal(jsonPackBench, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareDecodeGob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		buf := bytes.NewBuffer(gobPackBench)
		err := gob.NewDecoder(buf).Decode(&r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareDecodeZeroformatter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := zeroformatter.Deserialize(&r, zeroFmtpackBench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

/////////////////////////////////////////////////////////////////

func BenchmarkCompareEncodeArrayShamatonGen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamatongen.MarshalAsArray(&bench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeTinylib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := bench.MarshalMsg(nil)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeShamatonGen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamatongen.MarshalAsMap(&bench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeArrayShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.MarshalAsArray(bench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.MarshalAsMap(bench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeUgorji(b *testing.B) {
	for i := 0; i < b.N; i++ {

		b := []byte{}
		enc := codec.NewEncoderBytes(&b, mhBench)
		err := enc.Encode(bench)

		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeArrayVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {

		var buf bytes.Buffer
		enc := vmihailenco.NewEncoder(&buf)
		enc.UseArrayEncodedStructs(true)
		err := enc.Encode(bench)

		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(bench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeProtocolBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(protobench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(bench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeGob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(nil)
		err := gob.NewEncoder(buf).Encode(bench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeZeroformatter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := zeroformatter.Serialize(bench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
