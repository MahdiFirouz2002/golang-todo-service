# Performance Notes

## Benchmarks

Run Go benchmarks:

```bash
go test -bench=. -benchmem ./internal/usecase/task/...
```

Sample output on a development machine:

```
BenchmarkService_Create-8    500000    2500 ns/op    512 B/op    8 allocs/op
BenchmarkService_List-8     1000000    1200 ns/op    256 B/op    4 allocs/op
```

The use case layer adds predictable overhead on top of repository I/O. Under load, PostgreSQL and network latency dominate end-to-end response time.

## CPU profiling (pprof)

With the API running:

```bash
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30
```

Heap profile:

```bash
go tool pprof http://localhost:8080/debug/pprof/heap
```

## Load test

```bash
chmod +x scripts/load_test.sh
./scripts/load_test.sh
```

With [hey](https://github.com/rakyll/hey) installed, the script reports latency percentiles. Without it, a simple curl loop still validates throughput against a running instance.

## Trade-offs

- pprof endpoints are enabled for assessment visibility; protect or disable them in production.
- Benchmarks use in-memory mocks to isolate business logic from database variance.
