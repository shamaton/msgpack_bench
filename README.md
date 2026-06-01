# msgpack_bench

Benchmarks for MessagePack libraries and related serialization formats.

## Update benchmark SVGs

Run the benchmark suite with the `purego` build tag and regenerate SVG charts:

```sh
scripts/update-bench-svg.sh
```

The script writes the raw benchmark output to `docs/benchmarks/latest.txt` and
generates SVG charts in `docs/benchmarks/`.

Generated chart files:

- `compare-decode.svg`
- `compare-encode.svg`
- `usecase-decode.svg`
- `usecase-encode.svg`
- `primitive-decode.svg`
- `primitive-encode.svg`
- `stream.svg`
