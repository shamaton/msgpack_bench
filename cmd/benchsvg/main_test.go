package main

import (
	"strings"
	"testing"
)

func TestParseBenchmarks(t *testing.T) {
	input := strings.NewReader(`goos: darwin
goarch: arm64
BenchmarkCompareEncodeShamaton-10          12345      987.6 ns/op      128 B/op        4 allocs/op
BenchmarkUseCaseDecodeJson-10              9876      12.34 us/op      256 B/op        8 allocs/op
PASS
`)

	benches, err := parseBenchmarks(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(benches) != 2 {
		t.Fatalf("len = %d, want 2", len(benches))
	}
	if got, want := benches[0].Name, "CompareEncodeShamaton"; got != want {
		t.Fatalf("name = %q, want %q", got, want)
	}
	if got, want := benches[0].Section, "compare-encode"; got != want {
		t.Fatalf("section = %q, want %q", got, want)
	}
	if got, want := benches[0].Metrics[metricNSPerOp], 987.6; got != want {
		t.Fatalf("ns/op = %v, want %v", got, want)
	}
	if got, want := benches[0].Metrics["B/op"], 128.0; got != want {
		t.Fatalf("B/op = %v, want %v", got, want)
	}
	if got, want := benches[1].Metrics["us/op"], 12.34; got != want {
		t.Fatalf("us/op = %v, want %v", got, want)
	}
}

func TestFilterSectionSortsByMetric(t *testing.T) {
	benches := []benchmark{
		{Name: "CompareEncodeGob", Section: "compare-encode", Metrics: map[string]float64{metricNSPerOp: 40}},
		{Name: "CompareEncodeShamaton", Section: "compare-encode", Metrics: map[string]float64{metricNSPerOp: 10}},
		{Name: "UseCaseEncodeJson", Section: "usecase-encode", Metrics: map[string]float64{metricNSPerOp: 5}},
	}

	rows := filterSection(benches, "compare-encode", metricNSPerOp)
	if len(rows) != 2 {
		t.Fatalf("len = %d, want 2", len(rows))
	}
	if got, want := rows[0].Name, "CompareEncodeShamaton"; got != want {
		t.Fatalf("first row = %q, want %q", got, want)
	}
}

func TestFilterSectionSortsPrimitiveByTypeThenImplementation(t *testing.T) {
	benches := []benchmark{
		{Name: "MsgEncFloatVmihailenco", Section: "primitive-encode", Metrics: map[string]float64{metricNSPerOp: 1}},
		{Name: "MsgEncBoolVmihailenco", Section: "primitive-encode", Metrics: map[string]float64{metricNSPerOp: 1}},
		{Name: "MsgEncInterfaceShamaton", Section: "primitive-encode", Metrics: map[string]float64{metricNSPerOp: 1}},
		{Name: "MsgEncIntVmihailenco", Section: "primitive-encode", Metrics: map[string]float64{metricNSPerOp: 1}},
		{Name: "MsgEncFloatShamaton", Section: "primitive-encode", Metrics: map[string]float64{metricNSPerOp: 1}},
		{Name: "MsgEncArray10000Shamaton", Section: "primitive-encode", Metrics: map[string]float64{metricNSPerOp: 1}},
		{Name: "MsgEncBoolShamaton", Section: "primitive-encode", Metrics: map[string]float64{metricNSPerOp: 1}},
		{Name: "MsgEncIntShamaton", Section: "primitive-encode", Metrics: map[string]float64{metricNSPerOp: 1}},
	}

	rows := filterSection(benches, "primitive-encode", metricNSPerOp)
	got := make([]string, len(rows))
	for i := range rows {
		got[i] = rows[i].Name
	}
	want := []string{
		"MsgEncBoolShamaton",
		"MsgEncBoolVmihailenco",
		"MsgEncFloatShamaton",
		"MsgEncFloatVmihailenco",
		"MsgEncIntShamaton",
		"MsgEncIntVmihailenco",
		"MsgEncInterfaceShamaton",
		"MsgEncArray10000Shamaton",
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("row %d = %q, want %q; got all %v", i, got[i], want[i], got)
		}
	}
}

func TestFilterSectionSortsStreamByOperation(t *testing.T) {
	benches := []benchmark{
		{Name: "UnmarshalJsonStreamParallel", Section: "stream", Metrics: map[string]float64{metricNSPerOp: 1}},
		{Name: "MarshalMsgDirectParallel", Section: "stream", Metrics: map[string]float64{metricNSPerOp: 1}},
		{Name: "UnmarshalMsgStream", Section: "stream", Metrics: map[string]float64{metricNSPerOp: 1}},
		{Name: "MarshalMsgStreamParallel", Section: "stream", Metrics: map[string]float64{metricNSPerOp: 1}},
		{Name: "MarshalJsonStream", Section: "stream", Metrics: map[string]float64{metricNSPerOp: 1}},
		{Name: "UnmarshalMsgDirect", Section: "stream", Metrics: map[string]float64{metricNSPerOp: 1}},
		{Name: "MarshalMsgStream", Section: "stream", Metrics: map[string]float64{metricNSPerOp: 1}},
		{Name: "MarshalMsgDirect", Section: "stream", Metrics: map[string]float64{metricNSPerOp: 1}},
		{Name: "MarshalJsonDirect", Section: "stream", Metrics: map[string]float64{metricNSPerOp: 1}},
		{Name: "MarshalJsonDirectParallel", Section: "stream", Metrics: map[string]float64{metricNSPerOp: 1}},
		{Name: "MarshalJsonStreamParallel", Section: "stream", Metrics: map[string]float64{metricNSPerOp: 1}},
	}

	rows := filterSection(benches, "stream", metricNSPerOp)
	got := make([]string, len(rows))
	for i := range rows {
		got[i] = rows[i].Name
	}
	want := []string{
		"MarshalMsgDirect",
		"MarshalMsgStream",
		"MarshalJsonDirect",
		"MarshalJsonStream",
		"MarshalMsgDirectParallel",
		"MarshalMsgStreamParallel",
		"MarshalJsonDirectParallel",
		"MarshalJsonStreamParallel",
		"UnmarshalMsgDirect",
		"UnmarshalMsgStream",
		"UnmarshalJsonStreamParallel",
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("row %d = %q, want %q; got all %v", i, got[i], want[i], got)
		}
	}
}

func TestRenderSVGSupportsDarkTheme(t *testing.T) {
	rows := []benchmark{
		{Name: "CompareEncodeShamaton", Section: "compare-encode", Metrics: map[string]float64{metricNSPerOp: 10}},
	}

	svg := string(renderSVG("Benchmark", metricNSPerOp, rows))
	for _, want := range []string{
		`svg{color-scheme:light dark}`,
		`@media (prefers-color-scheme:dark)`,
		`<rect class="background" width="100%" height="100%"/>`,
		`.background{fill:#111827}`,
	} {
		if !strings.Contains(svg, want) {
			t.Fatalf("rendered SVG does not contain %q:\n%s", want, svg)
		}
	}
}

func TestRenderSVGIncludesBytesPerOpBarWhenAvailable(t *testing.T) {
	rows := []benchmark{
		{Name: "CompareEncodeShamaton", Section: "compare-encode", Metrics: map[string]float64{metricNSPerOp: 10, "B/op": 160}},
	}

	svg := string(renderSVG("Benchmark", metricNSPerOp, rows))
	for _, want := range []string{
		`10 ns/op`,
		`160 B/op`,
		`class="allocbar"`,
		`alloc bars: B/op`,
	} {
		if !strings.Contains(svg, want) {
			t.Fatalf("rendered SVG does not include %q:\n%s", want, svg)
		}
	}
	if strings.Contains(svg, `10 ns/op, 160 B/op`) {
		t.Fatalf("rendered SVG does not include B/op:\n%s", svg)
	}
}

func TestRenderSVGDoesNotDuplicateBytesPerOpMetric(t *testing.T) {
	rows := []benchmark{
		{Name: "CompareEncodeShamaton", Section: "compare-encode", Metrics: map[string]float64{"B/op": 160}},
	}

	svg := string(renderSVG("Benchmark", "B/op", rows))
	if strings.Contains(svg, `160 B/op, 160 B/op`) {
		t.Fatalf("rendered SVG duplicated B/op:\n%s", svg)
	}
	if strings.Contains(svg, `class="allocbar"`) {
		t.Fatalf("rendered SVG included allocation bar for B/op metric:\n%s", svg)
	}
	if !strings.Contains(svg, `160 B/op`) {
		t.Fatalf("rendered SVG does not include B/op:\n%s", svg)
	}
}

func TestClassifyStreamAndGeneratedBenchmarks(t *testing.T) {
	tests := map[string]string{
		"UnmarshalMsgStream":  "stream",
		"MarshalMsgDirect":    "stream",
		"MarshalJsonDirect":   "stream",
		"AppendMsgBenchChild": "tinylib-generated",
		"AppendMsgItem":       "tinylib-generated",
	}

	for name, want := range tests {
		if got := classify(name); got != want {
			t.Fatalf("classify(%q) = %q, want %q", name, got, want)
		}
	}
}
