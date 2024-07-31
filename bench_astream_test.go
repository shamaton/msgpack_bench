package bench

import (
	"bytes"
	"fmt"
	"github.com/go-json-experiment/json"
	"github.com/shamaton/msgpack/v2"
	"io"
	"reflect"
	"testing"
	"time"
)

var mtt []byte
var jtt []byte

func init() {
	{
		tt, err := msgpack.Marshal(hoge)
		if err != nil {
			panic(err)
		}
		var fuga Hoge
		err = msgpack.UnmarshalRead(bytes.NewReader(tt), &fuga)
		if err != nil {
			panic(err)
		}
		fmt.Println(hoge)
		fmt.Println(fuga)
		if !reflect.DeepEqual(hoge, fuga) {
			panic("fuga is not equal to hoge")
		}
		mtt = make([]byte, len(tt))
		copy(mtt, tt)

		fmt.Printf("% 02x\n", tt)
	}
	{
		buf := new(bytes.Buffer)
		err := msgpack.MarshalWrite(buf, hoge)
		if err != nil {
			panic(err)
		}
		fmt.Printf("% 02x\n", buf.Bytes())
		var fuga Hoge
		err = msgpack.UnmarshalRead(bytes.NewReader(buf.Bytes()), &fuga)
		if err != nil {
			panic(err)
		}
		fmt.Println(hoge)
		fmt.Println(fuga)
		if !reflect.DeepEqual(hoge, fuga) {
			panic("fuga is not equal to hoge")
		}
	}

	b, err := json.Marshal(hoge)
	if err != nil {
		panic(err)
	}
	jtt = b
}

type Hoge struct {
	Int     int
	String  string
	Bool    bool
	Uints   []uint
	Strings []string
	Time    time.Time
}

var hoge Hoge = Hoge{
	Int:     777,
	String:  "string string string",
	Bool:    true,
	Uints:   []uint{1, 2, 3, 4, 5, 10000},
	Strings: []string{"s", "h", "a", "maton"},
	Time:    time.Unix(12345, 0),
}

func BenchmarkAStreamUnmarshal(b *testing.B) {
	ttt, err := msgpack.Marshal(hoge)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewReader(ttt)
	b.ResetTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var fuga Hoge
		err = msgpack.UnmarshalRead(buf, &fuga)
		if err != nil {
			panic(err)
		}
		_, _ = buf.Seek(0, 0)
	}
	b.StopTimer()
	b.SetBytes(int64(len(ttt)))
}

func BenchmarkANoStreamUnmarshal(b *testing.B) {
	ttt, err := msgpack.Marshal(hoge)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewReader(ttt)
	b.ResetTimer()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		data, err := io.ReadAll(buf)
		if err != nil {
			panic(err)
		}
		var fuga Hoge
		err = msgpack.Unmarshal(data, &fuga)
		if err != nil {
			panic(err)
		}
		_, _ = buf.Seek(0, 0)
	}
	b.StopTimer()
	b.SetBytes(int64(len(ttt)))
}

func BenchmarkAStreamJsonUnmarshal(b *testing.B) {
	ttt, err := json.Marshal(hoge)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewReader(ttt)
	b.ResetTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var fuga Hoge
		err = json.UnmarshalRead(buf, &fuga)
		if err != nil {
			panic(err)
		}
		_, _ = buf.Seek(0, 0)
	}
	b.StopTimer()
	b.SetBytes(int64(len(ttt)))
}

func BenchmarkAStreamUnmarshal_Parallel(b *testing.B) {
	ttt, err := msgpack.Marshal(hoge)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var fuga Hoge
			buf := bytes.NewReader(ttt)
			err := msgpack.UnmarshalRead(buf, &fuga)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.StopTimer()
	b.SetBytes(int64(len(ttt)))
}

func BenchmarkANoStreamUnmarshal_Parallel(b *testing.B) {
	ttt, err := msgpack.Marshal(hoge)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var fuga Hoge
			buf := bytes.NewReader(ttt)
			data, err := io.ReadAll(buf)
			//data := make([]byte, len(ttt))
			//copy(data, ttt)
			err = msgpack.Unmarshal(data, &fuga)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.StopTimer()
	b.SetBytes(int64(len(ttt)))
}

func BenchmarkAStreamJsonUnmarshal_Parallel(b *testing.B) {
	ttt, err := json.Marshal(hoge)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var fuga Hoge
			buf := bytes.NewReader(ttt)
			err := json.UnmarshalRead(buf, &fuga)

			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.StopTimer()
	b.SetBytes(int64(len(ttt)))
}

func BenchmarkAStreamMarshal(b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.Buffer{}
		err := msgpack.MarshalWrite(&buf, hoge)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	b.SetBytes(int64(len(mtt)))
}

func BenchmarkANoStreamMarshal(b *testing.B) {
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf := bytes.Buffer{}
		_, err := msgpack.Marshal(hoge)
		if err != nil {
			b.Fatal(err)
		}
		_ = buf
	}
	b.StopTimer()
	b.SetBytes(int64(len(mtt)))
}

func BenchmarkAStreamJsonMarshal(b *testing.B) {
	data := make([]byte, 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		//buf := bytes.Buffer{}
		buf := bytes.NewBuffer(data)
		err := json.MarshalWrite(buf, hoge)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	b.SetBytes(int64(len(jtt)))
}

func BenchmarkAStreamMarshal_Parallel(b *testing.B) {
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := bytes.Buffer{}
			err := msgpack.MarshalWrite(&buf, hoge)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.StopTimer()
	b.SetBytes(int64(len(mtt)))
}

func BenchmarkANoStreamMarshal_Parallel(b *testing.B) {
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := bytes.Buffer{}
			_, err := msgpack.Marshal(hoge)
			if err != nil {
				b.Fatal(err)
			}
			_ = buf
		}
	})
	b.StopTimer()
	b.SetBytes(int64(len(mtt)))
}

func BenchmarkAStreamJsonMarshal_Parallel(b *testing.B) {
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := bytes.Buffer{}
			err := json.MarshalWrite(&buf, hoge)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.StopTimer()
	b.SetBytes(int64(len(jtt)))
}
