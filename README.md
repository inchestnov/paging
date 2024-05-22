# Compare standard (with offsets) and cursor paging

Current project demonstrated difference (in time) between standard and cursor paging.

# Local run

Assume docker host is running.
It is required because benchmarks use postgres in testcontainers.

```shell
# Clone project
git clone https://github.com/inchestnov/paging

# Go to directory with benchmark
cd paging/internal/pkg/paging

# Run benchmark
go test -bench=.
```

# Results

**goos: darwin <br> goarch: arm64**

| Paging type | Page size | Total elements |      Time (ns) |
|-------------|----------:|---------------:|---------------:|
| Standard    |       100 |            100 |        878 854 |
| Cursor      |       100 |            100 |        840 977 |
| Standard    |       100 |          1 000 |      9 628 672 |
| Cursor      |       100 |          1 000 |      9 049 814 |
| Standard    |       100 |         10 000 |    182 944 465 |
| Cursor      |       100 |         10 000 |    128 822 636 |
| Standard    |       500 |        100 000 |  2 715 294 917 |
| Cursor      |       500 |        100 000 |    929 520 834 |
| Standard    |       500 |      1 000 000 | 51 265 952 458 |
| Cursor      |       500 |      1 000 000 | 25 452 566 917 |