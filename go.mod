module github.com/shamaton/msgpack_bench

go 1.26.3

require (
	github.com/go-json-experiment/json v0.0.0-20240524174822-2d9f40f7385b
	github.com/shamaton/msgpack/v3 v3.2.0
	github.com/shamaton/msgpackgen v1.1.1
	github.com/shamaton/zeroformatter v1.0.1
	github.com/tinylib/msgp v1.6.4
	github.com/ugorji/go/codec v1.3.2
	github.com/vmihailenco/msgpack/v5 v5.4.1
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/dave/jennifer v1.7.1 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/mod v0.18.0 // indirect
	golang.org/x/tools v0.22.0 // indirect
)

tool (
	github.com/shamaton/msgpackgen
	github.com/tinylib/msgp
)
