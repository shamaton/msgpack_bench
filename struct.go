package bench

//go:generate go run github.com/shamaton/msgpackgen -strict
////go:generate msgpackgen -v
//go:generate go run github.com/tinylib/msgp -o msgp_gen.go

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
	Array  []int
	Map    map[string]uint
	Child  BenchChild
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
