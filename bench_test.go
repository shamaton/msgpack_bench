package bench_test

import (
	"bytes"
	"fmt"
	"testing"

	shamaton "github.com/shamaton/msgpack"
	vmihailenco "github.com/vmihailenco/msgpack"
)

type BenchChild struct {
	Int    int
	String string
}
type BenchMarkStruct struct {
	Int    int
	Uint   uint
	Float  float32
	Double float64
	Bool   bool
	String string
	Array  *[]int
	Map    map[string]uint
	Child  BenchChild
}

var v = BenchMarkStruct{
	Int:    -123,
	Uint:   456,
	Float:  1.234,
	Double: 6.789,
	Bool:   true,
	String: "this is text.",
	Array:  &[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
	Map:    map[string]uint{"this": 1, "is": 2, "map": 3},
	Child:  BenchChild{Int: 123456, String: "this is struct of child"},
}

// var v = 777
//var v = "thit is test"
//var v = []int{1, 2, 3, math.MinInt64}
// var v = []uint{1, 2, 3, 4, 5, 6, math.MaxUint64}
// var v = []string{"this", "is", "test"}
//var v = []interface{}{"aaa", math.MaxInt16, math.Pi, vv}
// var v = []byte("this is test byte array")

//var v = [4]string{"this", "is", "test"}

// var v = []float32{1.23, 4.56, math.MaxFloat32}
// var v = []float64{1.23, 4.56, math.MaxFloat64}
// var v = []bool{true, false, true}
// var v = []uint8{0x82, 0xa1, 0x41, 0x07, 0xa1, 0x42, 0xa1, 0x37}
// var v = [8]uint8{0x82, 0xa1, 0x41, 0x07, 0xa1, 0x42, 0xa1, 0x37}

// var v = map[string]BenchMarkStruct{"a": vv, "b": vv}
// var v = map[string]float32{"1": 2.34, "5": 6.78}
// var v = map[string]bool{"a": true, "b": false}
//var v = map[int]interface{}{1: 2, 3: "a", 4: []float32{1.23}}

//var v = time.Now()

var data []byte
var e2 error

var dataInt []byte

func init() {

	/*
		v = make([]int, 10000)
		for i := 0; i < 10000; i++ {
			v[i] = i
		}
	*/
	/*
		v = map[int]map[int]int{}
		for i := 0; i < 10000; i++ {
			v[i] = map[int]int{}
			for j := 0; j < 10; j++ {
				v[i][j] = i * j
			}
		}
	*/

	data, e2 = shamaton.Encode(v)
	if e2 != nil {
		fmt.Println("init err : ", e2)
	}

	dataInt, _ = shamaton.Encode(1234567)
}

func BenchmarkEncIntShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Encode(1234567)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkEncIntVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(1234567)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

//////////////////////////

func BenchmarkDecIntShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r int
		err := shamaton.Decode(dataInt, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkDecIntVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r int
		err := vmihailenco.Unmarshal(dataInt, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkDesShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := shamaton.Decode(data, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
func BenchmarkDesVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := vmihailenco.Unmarshal(data, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkArrayShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.EncodeStructAsArray(v)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkArrayVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {

		var buf bytes.Buffer
		enc := vmihailenco.NewEncoder(&buf).StructAsArray(true)
		err := enc.Encode(v)

		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMapShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.EncodeStructAsMap(v)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMapVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(v)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
