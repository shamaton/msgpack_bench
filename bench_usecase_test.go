package bench

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math"
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

	d, err := shamaton.MarshalAsArray(user)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	arrayMsgpackUser = d
	d, err = shamaton.MarshalAsMap(user)
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

	{
		dd, err := shamatongen.MarshalAsArray(user)
		if err != nil {
			fmt.Println("init err : ", err)
			os.Exit(1)
		}
		if !reflect.DeepEqual(arrayMsgpackUser, dd) {
			fmt.Println("not equal as array")
			os.Exit(1)
		}
		dd, err = shamatongen.MarshalAsMap(user)
		if err != nil {
			fmt.Println("init err : ", err)
			os.Exit(1)
		}
		if !reflect.DeepEqual(mapMsgpackUser, dd) {
			fmt.Println("not equal as map")
			os.Exit(1)
		}
	}

	{
		var v User
		err := shamaton.UnmarshalAsMap(mapMsgpackUser, &v)
		if err != nil {
			fmt.Println("init err : ", err)
			os.Exit(1)
		}
		if err = checkUseCaseDecodeValue(v); err != nil {
			fmt.Println("init err : ", err)
			os.Exit(1)
		}
	}
	{
		var v User
		err := shamaton.UnmarshalAsArray(arrayMsgpackUser, &v)
		if err != nil {
			fmt.Println("init err : ", err)
			os.Exit(1)
		}
		if err = checkUseCaseDecodeValue(v); err != nil {
			fmt.Println("init err : ", err)
			os.Exit(1)
		}
	}
	{
		var v User
		err := shamatongen.UnmarshalAsMap(mapMsgpackUser, &v)
		if err != nil {
			fmt.Println("init err : ", err)
			os.Exit(1)
		}
		if err = checkUseCaseDecodeValue(v); err != nil {
			fmt.Println("init err : ", err)
			os.Exit(1)
		}
	}
	{
		var v User
		err := shamatongen.UnmarshalAsArray(arrayMsgpackUser, &v)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err = checkUseCaseDecodeValue(v); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func checkUseCaseDecodeValue(u User) error {
	if len(user.Items) != len(u.Items) {
		return fmt.Errorf("item length is different %d, %d", len(user.Items), len(u.Items))
	}

	if user.ID != u.ID {
		return fmt.Errorf("id is different %d, %d", user.ID, u.ID)
	}
	if user.Name != u.Name {
		return fmt.Errorf("name is different %s, %s", user.Name, u.Name)
	}
	if user.Type != u.Type {
		return fmt.Errorf("type is different %v, %v", user.Type, u.Type)
	}
	if user.Level != u.Level {
		return fmt.Errorf("level is different %d, %d", user.Level, u.Level)
	}
	if user.Exp != u.Exp {
		return fmt.Errorf("exp is different %d, %d", user.Exp, u.Exp)
	}

	for i := range user.EquipIDs {
		if user.EquipIDs[i] != u.EquipIDs[i] {
			return fmt.Errorf("equip id is different %d, %d, %d", i, user.EquipIDs[i], u.EquipIDs[i])
		}
	}

	for i := range user.Items {
		if user.Items[i].ID != u.Items[i].ID {
			return fmt.Errorf("item id is different, %d, %d, %d", i, user.Items[i].ID, u.Items[i].ID)
		}
		if user.Items[i].Name != u.Items[i].Name {
			return fmt.Errorf("item name is different %d, %s, %s", i, user.Items[i].Name, u.Items[i].Name)
		}
		if user.Items[i].Effect != u.Items[i].Effect {
			return fmt.Errorf("item effect is different %d, %f, %f", i, user.Items[i].Effect, u.Items[i].Effect)
		}
		if user.Items[i].Num != u.Items[i].Num {
			return fmt.Errorf("item num is different %d, %d, %d", i, user.Items[i].Num, u.Items[i].Num)
		}
	}
	return nil
}

func BenchmarkUseCaseDecodeShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		err := shamaton.UnmarshalAsMap(mapMsgpackUser, &r)
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

func BenchmarkUseCaseDecodeShamatonGen(b *testing.B) {

	for i := 0; i < b.N; i++ {
		var r User
		err := shamatongen.UnmarshalAsMap(mapMsgpackUser, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseDecodeArrayShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		err := shamaton.UnmarshalAsArray(arrayMsgpackUser, &r)
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

func BenchmarkUseCaseDecodeArrayShamatonGen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		err := shamatongen.UnmarshalAsArray(arrayMsgpackUser, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseDecodeTinylib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		_, err := r.UnmarshalMsg(mapMsgpackUser)
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
		_, err := shamaton.MarshalAsMap(user)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseEncodeShamatonGen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamatongen.MarshalAsMap(&user)
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
		_, err := shamaton.MarshalAsArray(user)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkUseCaseEncodeArrayShamatonGen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamatongen.MarshalAsArray(&user)
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

func BenchmarkUseCaseEncodeTinylib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := user.MarshalMsg(nil)
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
