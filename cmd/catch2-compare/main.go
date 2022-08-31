package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lesomnus/catch2-compare/internal/catch2"
)

func usage() {
	fmt.Printf("\nUsage: %s [OPTIONS] TARGET SOURCE\n", os.Args[0])
	flag.PrintDefaults()
}

func verifyArgs(opts Options) error {
	if opts.WorkingDirectory != "" {
		if wd, err := filepath.Abs(opts.WorkingDirectory); err != nil {
			return fmt.Errorf("failed to resolve absolute path from %s: %w", opts.WorkingDirectory, err)
		} else {
			opts.WorkingDirectory = wd
			fmt.Printf("wd: %v\n", wd)
		}
	}

	if flag.NArg() != 2 {
		return fmt.Errorf("expected 2 arguments but it was %d", flag.NArg())
	}

	return nil
}

type Options struct {
	WorkingDirectory string
}

func main() {
	opts := Options{
		WorkingDirectory: "",
	}

	flag.StringVar(&opts.WorkingDirectory, "working-dir", "", "Display the file path relative to this path")

	flag.Usage = usage
	flag.Parse()

	if err := verifyArgs(opts); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())

		flag.Usage()
		os.Exit(1)
	}

	tgt_path := flag.Arg(0)
	src_path := flag.Arg(1)

	tgt_reports := make(map[string]catch2.Report)
	src_reports := make(map[string]catch2.Report)

	if err := load(tgt_path, tgt_reports); err != nil {
		log.Fatalln("failed to load target report:", err)
	}
	if err := load(src_path, src_reports); err != nil {
		log.Fatalln("failed to load source report:", err)
	}

	if opts.WorkingDirectory != "" {
		// Make relative path.
		mkRel := func(reports map[string]catch2.Report) {
			for _, report := range reports {
				for i, tc := range report.TestCases {
					if rel, err := filepath.Rel(opts.WorkingDirectory, tc.Filename); err != nil {
						continue
					} else if strings.HasPrefix(rel, "..") {
						continue
					} else {
						report.TestCases[i].Filename = rel
					}
				}
			}
		}

		mkRel(tgt_reports)
		mkRel(src_reports)
	}

	r := DiffReporter{}
	r.Report(os.Stdout, tgt_reports, src_reports)
}
