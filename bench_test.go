package bench

import (
	"fmt"
	"math"
	"testing"
	"time"

	shamaton "github.com/shamaton/msgpack"
	vmihailenco "github.com/vmihailenco/msgpack/v5"
)

var (
	Int        = int(1234567)
	Float      = float64(math.MaxFloat64)
	String     = "this_string_is_used_for_benchmark"
	Bool       = true
	Array      = []int{}
	Map        = map[string]int{}
	Byte       = []byte("this is test byte array")
	Interfaces = []interface{}{"aaa", uint64(math.MaxUint64), math.Pi, nil, true, []uint{1, 2, 3}, map[string]int{"a": 1, "b": 2}}
	Time       = time.Now()

	dataInt        []byte
	dataFloat      []byte
	dataString     []byte
	dataBool       []byte
	dataArray      []byte
	dataMap        []byte
	dataByte       []byte
	dataInterfaces []byte
	dataTime       []byte
)

/*
INSTALL PACKAGES
go get -u github.com/golang/protobuf/proto
go get -u github.com/shamaton/zeroformatter
go get -u github.com/ugorji/go/codec
go get -u github.com/vmihailenco/msgpack
*/

func init() {

	Array = make([]int, 10000)
	for i := 0; i < 10000; i++ {
		Array[i] = i * i
	}

	Map = make(map[string]int, 10000)
	for i := 0; i < 10000; i++ {
		Map[fmt.Sprint(i)+fmt.Sprint(i)] = i * i
	}

	dataInt, _ = shamaton.Encode(Int)
	dataFloat, _ = shamaton.Encode(Float)
	dataString, _ = shamaton.Encode(String)
	dataBool, _ = shamaton.Encode(Bool)
	dataArray, _ = shamaton.Encode(Array)
	dataMap, _ = shamaton.Encode(Map)
	dataByte, _ = shamaton.Encode(Byte)
	dataInterfaces, _ = shamaton.Encode(Interfaces)
	dataTime, _ = shamaton.Encode(Time)

	initCompare()
	initUseCase()
	RegisterGeneratedResolver()
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

func BenchmarkMsgEncBoolShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Encode(Bool)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgEncBoolVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(Bool)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
func BenchmarkMsgEncArrayShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Encode(Array)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgEncArrayVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(Array)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
func BenchmarkMsgEncMapShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Encode(Map)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgEncMapVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(Map)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgEncTimeShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Encode(Time)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgEncTimeVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(Time)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgEncByteShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Encode(Byte)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgEncByteVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(Byte)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
func BenchmarkMsgEncInterfaceShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Encode(Interfaces)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgEncInterfaceVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(Interfaces)
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

func BenchmarkMsgDecBoolShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r bool
		err := shamaton.Decode(dataBool, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgDecBoolVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r bool
		err := vmihailenco.Unmarshal(dataBool, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgDecArrayShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r []int
		err := shamaton.Decode(dataArray, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgDecArrayVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r []int
		err := vmihailenco.Unmarshal(dataArray, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
func BenchmarkMsgDecMapShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r map[string]int
		err := shamaton.Decode(dataMap, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgDecMapVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r map[string]int
		err := vmihailenco.Unmarshal(dataMap, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
func BenchmarkMsgDecTimeShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r time.Time
		err := shamaton.Decode(dataTime, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgDecTimeVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r time.Time
		err := vmihailenco.Unmarshal(dataTime, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgDecByteShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r []byte
		err := shamaton.Decode(dataByte, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgDecByteVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r []byte
		err := vmihailenco.Unmarshal(dataByte, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgDecInterfaceShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r []interface{}
		err := shamaton.Decode(dataInterfaces, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkMsgDecInterfaceVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r []interface{}
		err := vmihailenco.Unmarshal(dataInterfaces, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
