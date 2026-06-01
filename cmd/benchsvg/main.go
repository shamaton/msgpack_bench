package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"html"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const metricNSPerOp = "ns/op"

type benchmark struct {
	Name    string
	Section string
	Metrics map[string]float64
}

type sectionDef struct {
	Key   string
	Title string
}

var sections = []sectionDef{
	{Key: "compare-decode", Title: "Compare decode"},
	{Key: "compare-encode", Title: "Compare encode"},
	{Key: "usecase-decode", Title: "Use case decode"},
	{Key: "usecase-encode", Title: "Use case encode"},
	{Key: "primitive-decode", Title: "Primitive decode"},
	{Key: "primitive-encode", Title: "Primitive encode"},
	{Key: "stream", Title: "Stream"},
	{Key: "other", Title: "Other"},
}

var primitiveTypes = []string{
	"Bool",
	"Float",
	"Int",
	"String",
	"Time",
	"Byte",
	"Interface",
	"Array10000",
	"Map10000",
}

var primitiveImplementations = []string{
	"Shamaton",
	"Vmihailenco",
}

func main() {
	var inPath string
	var outDir string
	var metric string
	var title string

	flag.StringVar(&inPath, "in", "", "benchmark result file; stdin is used when empty")
	flag.StringVar(&outDir, "out-dir", "docs/benchmarks", "directory for generated SVG files")
	flag.StringVar(&metric, "metric", metricNSPerOp, "benchmark metric to chart")
	flag.StringVar(&title, "title", "msgpack benchmark", "title prefix for generated charts")
	flag.Parse()

	benches, err := readBenchmarks(inPath)
	if err != nil {
		exitf("read benchmark results: %v", err)
	}
	if len(benches) == 0 {
		exitf("read benchmark results: no benchmark rows found")
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		exitf("create output directory: %v", err)
	}

	generated := 0
	for _, section := range sections {
		rows := filterSection(benches, section.Key, metric)
		if len(rows) == 0 {
			continue
		}
		path := filepath.Join(outDir, section.Key+".svg")
		svg := renderSVG(title+": "+section.Title, metric, rows)
		if err := os.WriteFile(path, svg, 0o644); err != nil {
			exitf("write %s: %v", path, err)
		}
		fmt.Println(path)
		generated++
	}
	if generated == 0 {
		exitf("no benchmark rows had metric %q", metric)
	}
}

func readBenchmarks(path string) ([]benchmark, error) {
	var r io.Reader
	if path == "" {
		r = os.Stdin
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
	}
	return parseBenchmarks(r)
}

func parseBenchmarks(r io.Reader) ([]benchmark, error) {
	scanner := bufio.NewScanner(r)
	var benches []benchmark
	for scanner.Scan() {
		b, ok := parseBenchmarkLine(scanner.Text())
		if ok {
			benches = append(benches, b)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return benches, nil
}

func parseBenchmarkLine(line string) (benchmark, bool) {
	fields := strings.Fields(line)
	if len(fields) < 4 || !strings.HasPrefix(fields[0], "Benchmark") {
		return benchmark{}, false
	}
	name := trimBenchmarkName(fields[0])
	metrics := make(map[string]float64)
	for i := 2; i+1 < len(fields); i += 2 {
		value, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			continue
		}
		metrics[fields[i+1]] = value
	}
	if len(metrics) == 0 {
		return benchmark{}, false
	}
	return benchmark{
		Name:    name,
		Section: classify(name),
		Metrics: metrics,
	}, true
}

func trimBenchmarkName(field string) string {
	name := strings.TrimPrefix(field, "Benchmark")
	if i := strings.LastIndex(name, "-"); i >= 0 {
		if _, err := strconv.Atoi(name[i+1:]); err == nil {
			name = name[:i]
		}
	}
	return name
}

func classify(name string) string {
	switch {
	case strings.HasPrefix(name, "CompareDecode"):
		return "compare-decode"
	case strings.HasPrefix(name, "CompareEncode"):
		return "compare-encode"
	case strings.HasPrefix(name, "UseCaseDecode"):
		return "usecase-decode"
	case strings.HasPrefix(name, "UseCaseEncode"):
		return "usecase-encode"
	case strings.HasPrefix(name, "MsgDec"):
		return "primitive-decode"
	case strings.HasPrefix(name, "MsgEnc"):
		return "primitive-encode"
	case strings.Contains(name, "BenchChild") ||
		strings.Contains(name, "BenchMarkStruct") ||
		strings.Contains(name, "Item") ||
		strings.Contains(name, "User"):
		return "tinylib-generated"
	case strings.HasPrefix(name, "MarshalMsg") ||
		strings.HasPrefix(name, "MarshalJson") ||
		strings.HasPrefix(name, "UnmarshalMsg") ||
		strings.HasPrefix(name, "UnmarshalJson"):
		return "stream"
	default:
		return "other"
	}
}

func filterSection(benches []benchmark, section, metric string) []benchmark {
	var rows []benchmark
	for _, b := range benches {
		if b.Section != section {
			continue
		}
		if _, ok := b.Metrics[metric]; !ok {
			continue
		}
		rows = append(rows, b)
	}
	sortRows(rows, section, metric)
	return rows
}

func sortRows(rows []benchmark, section, metric string) {
	if section == "primitive-decode" || section == "primitive-encode" {
		sort.SliceStable(rows, func(i, j int) bool {
			ti, ii := primitiveSortKey(rows[i].Name)
			tj, ij := primitiveSortKey(rows[j].Name)
			if ti != tj {
				return ti < tj
			}
			if ii != ij {
				return ii < ij
			}
			return rows[i].Name < rows[j].Name
		})
		return
	}
	if section == "stream" {
		sort.SliceStable(rows, func(i, j int) bool {
			oi, pi, vi := streamSortKey(rows[i].Name)
			oj, pj, vj := streamSortKey(rows[j].Name)
			if oi != oj {
				return oi < oj
			}
			if pi != pj {
				return pi < pj
			}
			if vi != vj {
				return vi < vj
			}
			return rows[i].Name < rows[j].Name
		})
		return
	}

	sort.SliceStable(rows, func(i, j int) bool {
		vi := rows[i].Metrics[metric]
		vj := rows[j].Metrics[metric]
		if vi == vj {
			return rows[i].Name < rows[j].Name
		}
		return vi < vj
	})
}

func primitiveSortKey(name string) (int, int) {
	name = strings.TrimPrefix(name, "MsgEnc")
	name = strings.TrimPrefix(name, "MsgDec")

	implementation := ""
	for _, candidate := range primitiveImplementations {
		if strings.HasSuffix(name, candidate) {
			implementation = candidate
			name = strings.TrimSuffix(name, candidate)
			break
		}
	}

	return orderedPrimitiveTypeIndex(name), orderedValueIndex(implementation, primitiveImplementations)
}

func orderedPrimitiveTypeIndex(value string) int {
	for i, candidate := range primitiveTypes {
		if value == candidate {
			return i
		}
		if strings.HasPrefix(value, candidate) && isDigitSuffix(value[len(candidate):]) {
			return i
		}
	}
	return len(primitiveTypes)
}

func orderedValueIndex(value string, ordered []string) int {
	for i, candidate := range ordered {
		if value == candidate {
			return i
		}
	}
	return len(ordered)
}

func streamSortKey(name string) (int, int, int) {
	operation := 2
	switch {
	case strings.Contains(name, "Unmarshal"):
		operation = 1
	case strings.Contains(name, "Marshal"):
		operation = 0
	}

	parallel := 0
	if strings.Contains(name, "Parallel") {
		parallel = 1
	}

	return operation, parallel, streamVariantIndex(name)
}

func streamVariantIndex(name string) int {
	name = strings.TrimPrefix(name, "Marshal")
	name = strings.TrimPrefix(name, "Unmarshal")

	switch {
	case strings.HasPrefix(name, "MsgDirect"):
		return 0
	case strings.HasPrefix(name, "MsgStream"):
		return 1
	case strings.HasPrefix(name, "JsonDirect"):
		return 2
	case strings.HasPrefix(name, "JsonStream"):
		return 3
	default:
		return 4
	}
}

func isDigitSuffix(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func renderSVG(title, metric string, rows []benchmark) []byte {
	const (
		width       = 1120
		marginLeft  = 320
		marginRight = 170
		top         = 82
		rowHeight   = 28
		barHeight   = 16
		bottom      = 54
	)

	height := top + len(rows)*rowHeight + bottom
	plotWidth := width - marginLeft - marginRight
	minValue, maxValue := metricRange(rows, metric)
	useLog := shouldUseLogScale(minValue, maxValue)

	var buf bytes.Buffer
	fmt.Fprintf(&buf, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" role="img" aria-label="%s">`+"\n", width, height, width, height, esc(title))
	buf.WriteString(`<rect width="100%" height="100%" fill="#ffffff"/>` + "\n")
	buf.WriteString(`<style>
text{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;fill:#111827}
.subtitle{fill:#6b7280;font-size:13px}
.label{font-size:12px}
.value{fill:#374151;font-size:12px}
.axis{stroke:#d1d5db;stroke-width:1}
.grid{stroke:#e5e7eb;stroke-width:1}
</style>` + "\n")
	fmt.Fprintf(&buf, `<text x="24" y="32" font-size="22" font-weight="700">%s</text>`+"\n", esc(title))
	scaleLabel := "linear scale"
	if useLog {
		scaleLabel = "log scale"
	}
	fmt.Fprintf(&buf, `<text x="24" y="56" class="subtitle">metric: %s, lower is better, %s, generated %s</text>`+"\n", esc(metric), scaleLabel, time.Now().Format("2006-01-02"))

	axisX := marginLeft
	axisY1 := top - 16
	axisY2 := height - bottom + 8
	fmt.Fprintf(&buf, `<line x1="%d" y1="%d" x2="%d" y2="%d" class="axis"/>`+"\n", axisX, axisY1, axisX, axisY2)
	for i := 0; i <= 4; i++ {
		x := marginLeft + int(float64(plotWidth)*float64(i)/4)
		fmt.Fprintf(&buf, `<line x1="%d" y1="%d" x2="%d" y2="%d" class="grid"/>`+"\n", x, axisY1, x, axisY2)
	}

	for i, row := range rows {
		y := top + i*rowHeight
		value := row.Metrics[metric]
		barWidth := scaleValue(value, minValue, maxValue, float64(plotWidth), useLog)
		if barWidth < 2 {
			barWidth = 2
		}
		color := colorFor(row.Name)
		label := displayName(row.Name)
		fmt.Fprintf(&buf, `<text x="24" y="%d" class="label">%s</text>`+"\n", y+13, esc(label))
		fmt.Fprintf(&buf, `<rect x="%d" y="%d" width="%.1f" height="%d" rx="3" fill="%s"/>`+"\n", marginLeft, y, barWidth, barHeight, color)
		fmt.Fprintf(&buf, `<text x="%d" y="%d" class="value">%s</text>`+"\n", marginLeft+int(barWidth)+8, y+13, esc(formatMetric(value, metric)))
	}

	buf.WriteString("</svg>\n")
	return buf.Bytes()
}

func metricRange(rows []benchmark, metric string) (float64, float64) {
	minValue := math.Inf(1)
	maxValue := 0.0
	for _, row := range rows {
		value := row.Metrics[metric]
		if value <= 0 {
			continue
		}
		minValue = math.Min(minValue, value)
		maxValue = math.Max(maxValue, value)
	}
	if math.IsInf(minValue, 1) {
		return 0, 0
	}
	return minValue, maxValue
}

func shouldUseLogScale(minValue, maxValue float64) bool {
	return minValue > 0 && maxValue/minValue >= 50
}

func scaleValue(value, minValue, maxValue, width float64, useLog bool) float64 {
	if maxValue <= 0 {
		return 0
	}
	if maxValue == minValue {
		return width
	}
	if !useLog {
		return width * value / maxValue
	}
	if value <= 0 || minValue <= 0 {
		return 0
	}
	minLog := math.Log10(minValue)
	maxLog := math.Log10(maxValue)
	return width * (math.Log10(value) - minLog + 0.12) / (maxLog - minLog + 0.12)
}

func displayName(name string) string {
	replacer := strings.NewReplacer(
		"CompareDecode", "Compare decode ",
		"CompareEncode", "Compare encode ",
		"UseCaseDecode", "Use case decode ",
		"UseCaseEncode", "Use case encode ",
		"MsgDec", "Primitive decode ",
		"MsgEnc", "Primitive encode ",
	)
	return strings.TrimSpace(replacer.Replace(name))
}

func colorFor(name string) string {
	switch {
	case strings.Contains(name, "ShamatonGen"):
		return "#0f766e"
	case strings.Contains(name, "Shamaton"):
		return "#2563eb"
	case strings.Contains(name, "Vmihailenco"):
		return "#c2410c"
	case strings.Contains(name, "Json"):
		return "#ca8a04"
	case strings.Contains(name, "Gob"):
		return "#64748b"
	case strings.Contains(name, "ProtocolBuffer"):
		return "#7c3aed"
	case strings.Contains(name, "Zeroformatter"):
		return "#be123c"
	case strings.Contains(name, "Tinylib"):
		return "#15803d"
	case strings.Contains(name, "Ugorji"):
		return "#0891b2"
	default:
		return "#475569"
	}
}

func formatMetric(value float64, metric string) string {
	switch metric {
	case metricNSPerOp:
		return fmt.Sprintf("%s ns/op", formatNumber(value))
	case "B/op":
		if value >= 1024*1024 {
			return fmt.Sprintf("%.2f MiB/op", value/(1024*1024))
		}
		if value >= 1024 {
			return fmt.Sprintf("%.2f KiB/op", value/1024)
		}
		return fmt.Sprintf("%.0f B/op", value)
	case "allocs/op":
		return fmt.Sprintf("%.0f allocs/op", value)
	default:
		return fmt.Sprintf("%.2f %s", value, metric)
	}
}

func formatNumber(value float64) string {
	if value == math.Trunc(value) {
		return fmt.Sprintf("%.0f", value)
	}
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func esc(s string) string {
	return html.EscapeString(s)
}

func exitf(format string, args ...any) {
	err := fmt.Errorf(format, args...)
	if errors.Is(err, flag.ErrHelp) {
		os.Exit(0)
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
