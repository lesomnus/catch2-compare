package main_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	main "github.com/lesomnus/catch2-compare/cmd/catch2-compare"
	"github.com/lesomnus/catch2-compare/internal/catch2"
)

func TestDiffReporter(t *testing.T) {
	tgt := map[string]catch2.Report{
		"test 1": {
			Name: "test 1",
			TestCases: []catch2.TestCase{
				{
					Name:     "test case A",
					Filename: "/path/to/test-A.cpp",
					Line:     42,
					BenchmarkResults: []catch2.BenchmarkResult{
						{
							Name: "benchmark a",
							Mean: catch2.Measurement{
								Value: catch2.Duration{Duration: 4 * time.Millisecond},
							},
						},
						{
							Name: "benchmark b",
							Mean: catch2.Measurement{
								Value: catch2.Duration{Duration: 3 * time.Millisecond},
							},
						},
						{
							Name: "benchmark c",
							Mean: catch2.Measurement{
								Value: catch2.Duration{Duration: 7 * time.Second},
							},
						},
						{
							Name: "benchmark e",
							Mean: catch2.Measurement{
								Value: catch2.Duration{Duration: 10 * time.Microsecond},
							},
						},
						{
							Name: "benchmark g",
							Mean: catch2.Measurement{
								Value: catch2.Duration{Duration: 999},
							},
						},
					},
				},
				{
					Name:     "test case B",
					Filename: "/path/to/test-B.cpp",
					Line:     74,
					BenchmarkResults: []catch2.BenchmarkResult{
						{
							Name: "benchmark a",
							Mean: catch2.Measurement{
								Value: catch2.Duration{Duration: 9 * time.Millisecond},
							},
						},
					},
				},
			},
		},
		"test 2": {
			Name: "test 2",
			TestCases: []catch2.TestCase{
				{
					Name:     "test case A",
					Filename: "/path/to/test-A.cpp",
					Line:     314,
					BenchmarkResults: []catch2.BenchmarkResult{
						{
							Name: "benchmark a",
							Mean: catch2.Measurement{
								Value: catch2.Duration{Duration: 123 * time.Millisecond},
							},
						},
					},
				},
			},
		},
	}

	src := map[string]catch2.Report{
		"test 1": {
			Name: "test 1",
			TestCases: []catch2.TestCase{
				{
					Name:     "test case A",
					Filename: "/path/to/test-A.cpp",
					Line:     42,
					BenchmarkResults: []catch2.BenchmarkResult{
						{
							Name: "benchmark a",
							Mean: catch2.Measurement{
								Value: catch2.Duration{Duration: 2 * time.Millisecond},
							},
						},
						{
							Name: "benchmark b",
							Mean: catch2.Measurement{
								Value: catch2.Duration{Duration: 6 * time.Millisecond},
							},
						},
						{
							Name: "benchmark d",
							Mean: catch2.Measurement{
								Value: catch2.Duration{Duration: time.Minute + 7*time.Second},
							},
						},
						{
							Name: "benchmark f",
							Mean: catch2.Measurement{
								Value: catch2.Duration{Duration: 7 * time.Hour},
							},
						},
						{
							Name: "benchmark e",
							Mean: catch2.Measurement{
								Value: catch2.Duration{Duration: 2 * time.Microsecond},
							},
						},
						{
							Name: "benchmark g",
							Mean: catch2.Measurement{
								Value: catch2.Duration{Duration: 1000},
							},
						},
					},
				},
				{
					Name:     "test case C",
					Filename: "/path/to/test-C.cpp",
					Line:     69,
					BenchmarkResults: []catch2.BenchmarkResult{
						{
							Name: "benchmark a",
							Mean: catch2.Measurement{
								Value: catch2.Duration{Duration: 13 * time.Nanosecond},
							},
						},
					},
				},
			},
		},
		"test 3": {
			Name: "test 3",
			TestCases: []catch2.TestCase{
				{
					Name:     "test case A",
					Filename: "/path/to/test-A.cpp",
					Line:     12,
					BenchmarkResults: []catch2.BenchmarkResult{
						{
							Name: "benchmark a",
							Mean: catch2.Measurement{
								Value: catch2.Duration{Duration: 11*time.Second + 100*time.Nanosecond},
							},
						},
						{
							Name: "Very long long benchmark name",
							Mean: catch2.Measurement{
								Value: catch2.Duration{Duration: 1*time.Hour + 13*time.Second + 321*time.Nanosecond},
							},
						},
					},
				},
			},
		},
	}

	expected := `
@@ [test 1] test case A @@
# /path/to/test-A.cpp:42
+ benchmark a                       4ms          2ms     50.00%
- benchmark b                       3ms          6ms   -100.00%
  benchmark c                        7s            -          -
  benchmark d                         -         1m7s          -
+ benchmark e                      10µs          2µs     80.00%
  benchmark f                         -       7h0m0s          -
≈ benchmark g                     999ns          1µs     -0.10%

@@ [test 1] test case C @@
# /path/to/test-C.cpp:69
  benchmark a                         -         13ns          -

@@ [test 3] test case A @@
# /path/to/test-A.cpp:12
  benchmark a                         -          11s          -
  Very long long benchm...            -      1h0m13s          -
`

	b := new(strings.Builder)
	b.WriteString("\n")

	r := main.DiffReporter{}
	r.Report(b, tgt, src)

	require.Equal(t, expected, b.String())
}
