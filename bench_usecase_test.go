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

	shamaton "github.com/shamaton/msgpack/v3"
	"github.com/shamaton/msgpack_bench/protocmp"
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

	checkUseCaseEncodeOutputs()

	{
		dd, err :=  MarshalAsArray(user)
		if err != nil {
			fmt.Println("init err : ", err)
			os.Exit(1)
		}
		if !reflect.DeepEqual(arrayMsgpackUser, dd) {
			fmt.Println("not equal as array")
			os.Exit(1)
		}
		dd, err =  MarshalAsMap(user)
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
		err :=  UnmarshalAsMap(mapMsgpackUser, &v)
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
		err :=  UnmarshalAsArray(arrayMsgpackUser, &v)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err = checkUseCaseDecodeValue(v); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	{
		var v User
		err := vmihailenco.Unmarshal(mapMsgpackUser, &v)
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
		err := vmihailenco.Unmarshal(arrayMsgpackUser, &v)
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
		_, err := v.UnmarshalMsg(mapMsgpackUser)
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
		err := codec.NewDecoderBytes(mapMsgpackUser, mhUser).Decode(&v)
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
		err := zeroformatter.Deserialize(&v, zeroFmtpackUser)
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
		err := json.Unmarshal(jsonPackUser, &v)
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
		err := gob.NewDecoder(bytes.NewBuffer(gobPackUser)).Decode(&v)
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
		var v protocmp.User
		err := proto.Unmarshal(protoPackUser, &v)
		if err != nil {
			fmt.Println("init err : ", err)
			os.Exit(1)
		}
		if err = checkUseCaseDecodeValue(userFromProto(&v)); err != nil {
			fmt.Println("init err : ", err)
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

	if len(user.EquipIDs) != len(u.EquipIDs) {
		return fmt.Errorf("equip id length is different %d, %d", len(user.EquipIDs), len(u.EquipIDs))
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

func userFromProto(p *protocmp.User) User {
	items := make([]Item, len(p.Items))
	for i, item := range p.Items {
		items[i] = Item{
			ID:     int(item.ID),
			Name:   item.Name,
			Effect: item.Effect,
			Num:    uint(item.Num),
		}
	}
	return User{
		ID:       int(p.ID),
		Name:     p.Name,
		Level:    uint(p.Level),
		Exp:      p.Exp,
		Type:     p.Type,
		EquipIDs: append([]uint32(nil), p.EquipIDs...),
		Items:    items,
	}
}

func checkUseCaseEncodeOutputs() {
	var v User

	d, err :=  MarshalAsArray(&user)
	mustUseCaseCheck("msgpackgen encode array", err)
	mustUseCaseCheck("msgpackgen encoded array decode",  UnmarshalAsArray(d, &v))
	mustUseCaseValue("msgpackgen array", v)

	v = User{}
	d, err =  MarshalAsMap(&user)
	mustUseCaseCheck("msgpackgen encode map", err)
	mustUseCaseCheck("msgpackgen encoded map decode",  UnmarshalAsMap(d, &v))
	mustUseCaseValue("msgpackgen map", v)

	v = User{}
	d, err = shamaton.MarshalAsArray(user)
	mustUseCaseCheck("shamaton encode array", err)
	mustUseCaseCheck("shamaton encoded array decode", shamaton.UnmarshalAsArray(d, &v))
	mustUseCaseValue("shamaton array", v)

	v = User{}
	d, err = shamaton.MarshalAsMap(user)
	mustUseCaseCheck("shamaton encode map", err)
	mustUseCaseCheck("shamaton encoded map decode", shamaton.UnmarshalAsMap(d, &v))
	mustUseCaseValue("shamaton map", v)

	v = User{}
	d, err = user.MarshalMsg(nil)
	mustUseCaseCheck("tinylib encode", err)
	_, err = v.UnmarshalMsg(d)
	mustUseCaseCheck("tinylib encoded decode", err)
	mustUseCaseValue("tinylib", v)

	v = User{}
	buf := []byte{}
	mustUseCaseCheck("ugorji encode", codec.NewEncoderBytes(&buf, mhUser).Encode(user))
	mustUseCaseCheck("ugorji encoded decode", codec.NewDecoderBytes(buf, mhUser).Decode(&v))
	mustUseCaseValue("ugorji", v)

	v = User{}
	var bytesBuf bytes.Buffer
	enc := vmihailenco.NewEncoder(&bytesBuf)
	enc.UseArrayEncodedStructs(true)
	mustUseCaseCheck("vmihailenco encode array", enc.Encode(user))
	mustUseCaseCheck("vmihailenco encoded array decode", vmihailenco.Unmarshal(bytesBuf.Bytes(), &v))
	mustUseCaseValue("vmihailenco array", v)

	v = User{}
	d, err = vmihailenco.Marshal(user)
	mustUseCaseCheck("vmihailenco encode map", err)
	mustUseCaseCheck("vmihailenco encoded map decode", vmihailenco.Unmarshal(d, &v))
	mustUseCaseValue("vmihailenco map", v)

	var pb protocmp.User
	d, err = proto.Marshal(protouser)
	mustUseCaseCheck("proto encode", err)
	mustUseCaseCheck("proto encoded decode", proto.Unmarshal(d, &pb))
	mustUseCaseValue("proto", userFromProto(&pb))

	v = User{}
	d, err = json.Marshal(user)
	mustUseCaseCheck("json encode", err)
	mustUseCaseCheck("json encoded decode", json.Unmarshal(d, &v))
	mustUseCaseValue("json", v)

	v = User{}
	bytesBuf.Reset()
	mustUseCaseCheck("gob encode", gob.NewEncoder(&bytesBuf).Encode(user))
	mustUseCaseCheck("gob encoded decode", gob.NewDecoder(bytes.NewBuffer(bytesBuf.Bytes())).Decode(&v))
	mustUseCaseValue("gob", v)

	v = User{}
	d, err = zeroformatter.Serialize(user)
	mustUseCaseCheck("zeroformatter encode", err)
	mustUseCaseCheck("zeroformatter encoded decode", zeroformatter.Deserialize(&v, d))
	mustUseCaseValue("zeroformatter", v)
}

func mustUseCaseCheck(name string, err error) {
	if err != nil {
		fmt.Println("init err: ", name, err)
		os.Exit(1)
	}
}

func mustUseCaseValue(name string, got User) {
	if err := checkUseCaseDecodeValue(got); err != nil {
		fmt.Println("not equal: user vs ", name, err)
		fmt.Printf("user: %+v\n", user)
		fmt.Printf("%s: %+v\n", name, got)
		os.Exit(1)
	}
}

func BenchmarkUseCaseDecodeShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		err := shamaton.UnmarshalAsMap(mapMsgpackUser, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseDecodeVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		err := vmihailenco.Unmarshal(mapMsgpackUser, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseDecodeShamatonGenMap(b *testing.B) {

	for i := 0; i < b.N; i++ {
		var r User
		err := UnmarshalAsMap(mapMsgpackUser, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseDecodeShamatonGen(b *testing.B) {

	for i := 0; i < b.N; i++ {
		var r User
		err := Unmarshal(mapMsgpackUser, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseDecodeArrayShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		err := shamaton.UnmarshalAsArray(arrayMsgpackUser, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseDecodeArrayVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		err := vmihailenco.Unmarshal(arrayMsgpackUser, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseDecodeShamatonGenArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		err := UnmarshalAsArray(arrayMsgpackUser, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseDecodeTinylib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		_, err := r.UnmarshalMsg(mapMsgpackUser)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseDecodeUgorji(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		dec := codec.NewDecoderBytes(mapMsgpackUser, mhUser)
		err := dec.Decode(&r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseDecodeZeroformatter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		err := zeroformatter.Deserialize(&r, zeroFmtpackUser)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseDecodeJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		err := json.Unmarshal(jsonPackUser, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseDecodeGob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r User
		buf := bytes.NewBuffer(gobPackUser)
		err := gob.NewDecoder(buf).Decode(&r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseDecodeProtocolBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r protocmp.User
		err := proto.Unmarshal(protoPackUser, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/////////////////////////////////////////////////////////////////

func BenchmarkUseCaseEncodeShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.MarshalAsMap(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseEncodeShamatonGenMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := MarshalAsMap(&user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseEncodeShamatonGen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Marshal(&user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseEncodeVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseEncodeArrayShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.MarshalAsArray(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseEncodeShamatonGenArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := MarshalAsArray(&user)
		if err != nil {
			b.Fatal(err)
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
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseEncodeTinylib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := user.MarshalMsg(nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseEncodeUgorji(b *testing.B) {
	for i := 0; i < b.N; i++ {

		buf := []byte{}
		enc := codec.NewEncoderBytes(&buf, mhUser)
		err := enc.Encode(user)

		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseEncodeZeroformatter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := zeroformatter.Serialize(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseEncodeJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseEncodeGob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(nil)
		err := gob.NewEncoder(buf).Encode(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUseCaseEncodeProtocolBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(protouser)
		if err != nil {
			b.Fatal(err)
		}
	}
}
