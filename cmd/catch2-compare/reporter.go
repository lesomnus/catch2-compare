package main

import (
	"fmt"
	"io"
	"math"
	"sort"
	"time"

	"github.com/lesomnus/catch2-compare/internal/catch2"
)

type Reporter interface {
	Report(w io.Writer, target map[string]catch2.Report, source map[string]catch2.Report) error
}

type DiffReporter struct{}

func (r *DiffReporter) Report(w io.Writer, target map[string]catch2.Report, source map[string]catch2.Report) error {
	names := make([]string, 0, len(source))
	for name := range source {
		names = append(names, name)
	}
	sort.Strings(names)

	isFirst := true

	for _, name := range names {
		srcReport := source[name]
		tgtReport, ok := target[name]
		if !ok {
			tgtReport = catch2.Report{
				Name:      name,
				TestCases: make([]catch2.TestCase, 0),
			}
		}

		type tcPair struct {
			tgt *catch2.TestCase
			src *catch2.TestCase
		}

		tcNames := make(map[string]*tcPair)
		tcs := make([]*tcPair, 0)

		for i, tc := range srcReport.TestCases {
			pair := &tcPair{
				tgt: &catch2.TestCase{
					Name:             tc.Name,
					Filename:         tc.Filename,
					Line:             tc.Line,
					BenchmarkResults: make([]catch2.BenchmarkResult, 0),
				},
				src: &srcReport.TestCases[i],
			}
			tcNames[tc.Name] = pair
			tcs = append(tcs, pair)
		}
		for i, tc := range tgtReport.TestCases {
			pair, ok := tcNames[tc.Name]
			if !ok {
				continue
			}

			pair.tgt = &tgtReport.TestCases[i]
		}

		for _, pair := range tcs {
			if isFirst {
				isFirst = false
			} else {
				if _, err := io.WriteString(w, "\n"); err != nil {
					return err
				}
			}

			tcName := fmt.Sprintf("[%s] %s", name, pair.tgt.Name)
			pair.tgt.Name = tcName
			pair.src.Name = tcName

			if err := r.printTestCase(w, *pair.tgt, *pair.src); err != nil {
				return fmt.Errorf("failed to print test case %s: %w", tcName, err)
			}
		}
	}

	return nil
}

func (r *DiffReporter) printTestCase(w io.Writer, target catch2.TestCase, source catch2.TestCase) error {
	if target.Name != source.Name {
		panic(fmt.Sprintf("test case names do not match: %s != %s", target.Name, source.Name))
	}

	if _, err := fmt.Fprintln(w, "::: "+target.Name); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "+++ %s:%d\n", target.Filename, target.Line); err != nil {
		return err
	}
	if target.Filename != source.Filename || target.Line != source.Line {
		if _, err := fmt.Fprintf(w, "+++ %s:%d\n", source.Filename, source.Line); err != nil {
			return err
		}
	}

	type rstPair struct {
		tgt *catch2.BenchmarkResult
		src *catch2.BenchmarkResult
	}

	benchNames := make(map[string]*rstPair)
	benchResults := make([]*rstPair, 0)

	for i := 0; ; i++ {
		tgt_end := i >= len(target.BenchmarkResults)
		src_end := i >= len(source.BenchmarkResults)

		if tgt_end && src_end {
			break
		}

		if !tgt_end {
			rst := &target.BenchmarkResults[i]
			pair, ok := benchNames[rst.Name]
			if ok {
				pair.tgt = rst
			} else {
				pair = &rstPair{tgt: rst}
				benchNames[rst.Name] = pair
				benchResults = append(benchResults, pair)
			}
		}

		if !src_end {
			rst := &source.BenchmarkResults[i]
			pair, ok := benchNames[rst.Name]
			if ok {
				pair.src = rst
			} else {
				pair = &rstPair{src: rst}
				benchNames[rst.Name] = pair
				benchResults = append(benchResults, pair)
			}
		}
	}

	empty := catch2.BenchmarkResult{Name: ""}
	for _, rst := range benchResults {
		var name string
		if rst.tgt == nil {
			rst.tgt = &empty
			name = rst.src.Name
		} else if rst.src == nil {
			rst.src = &empty
			name = rst.tgt.Name
		} else {
			name = rst.tgt.Name
		}

		if err := r.printBenchmarkResult(w, *rst.tgt, *rst.src); err != nil {
			return fmt.Errorf("failed to print benchmark result \"%s\": %w", name, err)
		}
		if _, err := io.WriteString(w, "\n"); err != nil {
			return err
		}
	}

	return nil
}

func (r *DiffReporter) printBenchmarkResult(w io.Writer, target catch2.BenchmarkResult, source catch2.BenchmarkResult) error {
	tgt_empty := len(target.Name) == 0
	src_empty := len(source.Name) == 0
	has_empty := tgt_empty || src_empty
	if tgt_empty && src_empty {
		panic("both benchmarks have no name")
	}
	if !has_empty && target.Name != source.Name {
		panic(fmt.Sprintf("benchmark result names do not match: %s != %s", target.Name, source.Name))
	}

	name := target.Name
	if tgt_empty {
		name = source.Name
	}

	// Omit if test name too long.
	if len(name) > 24 {
		name = fmt.Sprintf("%s...", name[0:21])
	}

	tv := target.Mean.Value.Duration
	sv := source.Mean.Value.Duration

	if has_empty {
		if tgt_empty {
			if _, err := fmt.Fprintf(w, "  %-24s %12s %12v %10s", name, "-", r.truncateDuration(sv), "-"); err != nil {
				return err
			}
		} else if src_empty {
			if _, err := fmt.Fprintf(w, "  %-24s %12v %12s %10s", name, r.truncateDuration(tv), "-", "-"); err != nil {
				return err
			}
		}
	} else {
		rate := float64(tv-sv) / float64(tv) * 100

		sign := "≈"
		if rate > 1.0 {
			sign = "+"
		} else if rate < -1.0 {
			sign = "-"
		}

		if _, err := fmt.Fprintf(w, "%s %-24s %12v %12v %9.2f%%", sign, target.Name, r.truncateDuration(tv), r.truncateDuration(sv), rate); err != nil {
			return err
		}
	}

	return nil
}

func (r *DiffReporter) truncateDuration(d time.Duration) time.Duration {
	v := int(d)
	if d < 0 {
		v = -v
	}
	if v < 1e6 {
		return d
	}

	digits := 0
	w := v
	for w != 0 {
		w /= 10
		digits++
	}

	v = v - (v % int(math.Pow10(digits-6)))
	if d < 0 {
		v = -v
	}

	return time.Duration(v)
}
