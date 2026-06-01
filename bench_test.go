package bench

import (
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	shamaton "github.com/shamaton/msgpack/v3"
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
go get -u google.golang.org/protobuf/proto
go get -u github.com/shamaton/zeroformatter
go get -u github.com/ugorji/go/codec
go get -u github.com/vmihailenco/msgpack
*/

func init() {
	// RegisterGeneratedResolver()

	Array = make([]int, 10000)
	for i := 0; i < 10000; i++ {
		Array[i] = i * i
	}

	Map = make(map[string]int, 10000)
	for i := 0; i < 10000; i++ {
		Map[fmt.Sprint(i)+fmt.Sprint(i)] = i * i
	}

	dataInt = mustMarshal("int", Int)
	dataFloat = mustMarshal("float", Float)
	dataString = mustMarshal("string", String)
	dataBool = mustMarshal("bool", Bool)
	dataArray = mustMarshal("array", Array)
	dataMap = mustMarshal("map", Map)
	dataByte = mustMarshal("byte", Byte)
	dataInterfaces = mustMarshal("interfaces", Interfaces)
	dataTime = mustMarshal("time", Time)
}

func mustMarshal(name string, v any) []byte {
	d, err := shamaton.Marshal(v)
	if err != nil {
		fmt.Println("init err: ", name, err)
		os.Exit(1)
	}
	return d
}

func BenchmarkMsgEncIntShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Marshal(Int)
		if err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkMsgEncIntVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(Int)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgEncFloatShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Marshal(Float)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgEncFloatVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(Float)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgEncStringShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Marshal(String)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgEncStringVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(String)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgEncBoolShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Marshal(Bool)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgEncBoolVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(Bool)
		if err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkMsgEncArray10000Shamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Marshal(Array)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgEncArray10000Vmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(Array)
		if err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkMsgEncMap10000Shamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Marshal(Map)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgEncMap10000Vmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(Map)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgEncTimeShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Marshal(Time)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgEncTimeVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(Time)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgEncByteShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Marshal(Byte)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgEncByteVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(Byte)
		if err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkMsgEncInterfaceShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.Marshal(Interfaces)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgEncInterfaceVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(Interfaces)
		if err != nil {
			b.Fatal(err)
		}
	}
}

//////////////////////////

func BenchmarkMsgDecIntShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r int
		err := shamaton.Unmarshal(dataInt, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgDecIntVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r int
		err := vmihailenco.Unmarshal(dataInt, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgDecFloatShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r float64
		err := shamaton.Unmarshal(dataFloat, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgDecFloatVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r float64
		err := vmihailenco.Unmarshal(dataFloat, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkMsgDecStringShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r string
		err := shamaton.Unmarshal(dataString, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgDecStringVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r string
		err := vmihailenco.Unmarshal(dataString, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgDecBoolShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r bool
		err := shamaton.Unmarshal(dataBool, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgDecBoolVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r bool
		err := vmihailenco.Unmarshal(dataBool, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgDecArray10000Shamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r []int
		err := shamaton.Unmarshal(dataArray, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgDecArray10000Vmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r []int
		err := vmihailenco.Unmarshal(dataArray, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkMsgDecMap10000Shamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r map[string]int
		err := shamaton.Unmarshal(dataMap, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgDecMap10000Vmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r map[string]int
		err := vmihailenco.Unmarshal(dataMap, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkMsgDecTimeShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r time.Time
		err := shamaton.Unmarshal(dataTime, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgDecTimeVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r time.Time
		err := vmihailenco.Unmarshal(dataTime, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgDecByteShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r []byte
		err := shamaton.Unmarshal(dataByte, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgDecByteVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r []byte
		err := vmihailenco.Unmarshal(dataByte, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgDecInterfaceShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r []interface{}
		err := shamaton.Unmarshal(dataInterfaces, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgDecInterfaceVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r []interface{}
		err := vmihailenco.Unmarshal(dataInterfaces, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}
