package bench_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	shamaton "github.com/shamaton/msgpack"
	"github.com/shamaton/msgpack_bench/protocmp"
	"github.com/shamaton/zeroformatter"
	"github.com/ugorji/go/codec"
	vmihailenco "github.com/vmihailenco/msgpack"
)

type Item struct {
	ID     int
	Name   string
	Effect float32
	Num    uint
}

type User struct {
	ID       int
	Name     string
	Level    uint
	Exp      uint64
	Type     bool
	EquipIDs []uint32
	Items    []Item
}

var user = User{
	ID:       12345,
	Name:     "しゃまとん",
	Level:    99,
	Exp:      math.MaxUint32 * 2,
	Type:     true,
	EquipIDs: []uint32{1, 100, 10000, 1000000, 100000000},
	Items:    []Item{},
}

var protouser = &protocmp.User{
	ID:       int32(user.ID),
	Name:     user.Name,
	Level:    uint32(user.Level),
	Exp:      user.Exp,
	Type:     user.Type,
	EquipIDs: user.EquipIDs,
	Items:    []*protocmp.Item{},
}

var (
	arrayMsgpackUser []byte
	mapMsgpackUser   []byte
	zeroFmtpackUser  []byte
	jsonPackUser     []byte
	gobPackUser      []byte
	protoPackUser    []byte
)

// for codec
var (
	mhUser = &codec.MsgpackHandle{}
)

func initUseCase() {
	// ugorji
	//mhUser.MapType = reflect.TypeOf(user)

	// item
	for i := 0; i < 100; i++ {
		name := "item" + fmt.Sprint(i)
		item := Item{
			ID:     i,
			Name:   name,
			Effect: float32(i*i) / 3.0,
			Num:    uint(i * i * i * i),
		}
		user.Items = append(user.Items, item)

		pItem := &protocmp.Item{
			ID:     int32(item.ID),
			Name:   item.Name,
			Effect: item.Effect,
			Num:    uint32(item.Num),
		}
		protouser.Items = append(protouser.Items, pItem)
	}

	d, err := shamaton.EncodeStructAsArray(user)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	arrayMsgpackUser = d
	d, err = shamaton.EncodeStructAsMap(user)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	mapMsgpackUser = d

	d, err = zeroformatter.Serialize(user)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	zeroFmtpackUser = d

	d, err = json.Marshal(user)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	jsonPackUser = d

	d, err = proto.Marshal(protouser)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	protoPackUser = d

	buf := bytes.NewBuffer(nil)
	err = gob.NewEncoder(buf).Encode(user)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	gobPackUser = buf.Bytes()
}

func BenchmarkUseCaseDecodeShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		err := shamaton.DecodeStructAsMap(mapMsgpackUser, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
func BenchmarkUseCaseDecodeVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		err := vmihailenco.Unmarshal(mapMsgpackUser, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseDecodeArrayShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		err := shamaton.DecodeStructAsArray(arrayMsgpackUser, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
func BenchmarkUseCaseDecodeArrayVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		err := vmihailenco.Unmarshal(arrayMsgpackUser, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseDecodeUgorji(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		dec := codec.NewDecoderBytes(mapMsgpackUser, mhUser)
		err := dec.Decode(&r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseDecodeZeroformatter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		err := zeroformatter.Deserialize(&r, zeroFmtpackUser)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseDecodeJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		err := json.Unmarshal(jsonPackUser, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseDecodeGob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		buf := bytes.NewBuffer(gobPackUser)
		err := gob.NewDecoder(buf).Decode(&r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseDecodeProtocolBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r protocmp.User
		err := proto.Unmarshal(protoPackUser, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

/////////////////////////////////////////////////////////////////

func BenchmarkUseCaseEncodeShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.EncodeStructAsMap(user)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseEncodeVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(user)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseEncodeArrayShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.EncodeStructAsArray(user)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseEncodeArrayVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {

		var buf bytes.Buffer
		enc := vmihailenco.NewEncoder(&buf)
		enc.UseArrayEncodedStructs(true)
		err := enc.Encode(user)

		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseEncodeUgorji(b *testing.B) {
	for i := 0; i < b.N; i++ {

		b := []byte{}
		enc := codec.NewEncoderBytes(&b, mhUser)
		err := enc.Encode(user)

		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseEncodeZeroformatter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := zeroformatter.Serialize(user)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseEncodeJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(user)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseEncodeGob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(nil)
		err := gob.NewEncoder(buf).Encode(user)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseEncodeProtocolBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(protouser)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
