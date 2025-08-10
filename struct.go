package bench

//go:generate go run github.com/shamaton/msgpackgen -strict
////go:generate msgpackgen -v
//go:generate go run github.com/tinylib/msgp -o msgp_gen.go

type BenchChild struct {
	Int    int    `msgpack:",omitempty"`
	String string `msgpack:",omitempty"`
}
type BenchMarkStruct struct {
	Int    int             `msgpack:",omitempty"`
	Uint   uint            `msgpack:",omitempty"`
	Float  float32         `msgpack:",omitempty"`
	Double float64         `msgpack:",omitempty"`
	Bool   bool            `msgpack:",omitempty"`
	String string          `msgpack:",omitempty"`
	Array  []int           `msgpack:",omitempty"`
	Map    map[string]uint `msgpack:",omitempty"`
	Child  BenchChild      `msgpack:",omitempty"`
}

type BenchChildOmitempty struct {
	Int    int    `msgpack:",omitempty"`
	String string `msgpack:",omitempty"`
}
type BenchMarkStructOmitempty struct {
	Int    int             `msgpack:",omitempty"`
	Uint   uint            `msgpack:",omitempty"`
	Float  float32         `msgpack:",omitempty"`
	Double float64         `msgpack:",omitempty"`
	Bool   bool            `msgpack:",omitempty"`
	String string          `msgpack:",omitempty"`
	Array  []int           `msgpack:",omitempty"`
	Map    map[string]uint `msgpack:",omitempty"`
	Child  BenchChild      `msgpack:",omitempty"`
}

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
