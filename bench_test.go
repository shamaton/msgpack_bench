package bench_test

import (
	"fmt"
	"math"
	"testing"

	shamaton "github.com/shamaton/msgpack"
	vmihailenco "github.com/vmihailenco/msgpack"
)

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

var (
	Int    = int(1234567)
	Float  = float64(math.MaxFloat64)
	String = "this_string_is_used_for_benchmark"

	dataInt    []byte
	dataFloat  []byte
	dataString []byte
)

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

	dataInt, _ = shamaton.Encode(Int)
	dataFloat, _ = shamaton.Encode(Float)
	dataString, _ = shamaton.Encode(String)
}

func BenchmarkMsgEncIntShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Encode(Int)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
func BenchmarkMsgEncIntVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(Int)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgEncFloatShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Encode(Float)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgEncFloatVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(Float)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgEncStringShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Encode(String)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgEncStringVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(String)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

//////////////////////////

func BenchmarkMsgDecIntShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r int
		err := shamaton.Decode(dataInt, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgDecIntVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r int
		err := vmihailenco.Unmarshal(dataInt, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgDecFloatShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r float64
		err := shamaton.Decode(dataFloat, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgDecFloatVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r float64
		err := vmihailenco.Unmarshal(dataFloat, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
func BenchmarkMsgDecStringShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r string
		err := shamaton.Decode(dataString, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgDecStringVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r string
		err := vmihailenco.Unmarshal(dataString, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
