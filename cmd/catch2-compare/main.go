package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/lesomnus/catch2-compare/internal/catch2"
)

func usage() {
	fmt.Printf("\nUsage: %s [OPTIONS] TARGET SOURCE\n", os.Args[0])
	flag.PrintDefaults()
}

func run() int {
	tgt_path := flag.Arg(0)
	src_path := flag.Arg(1)

	tgt_reports := make(map[string]catch2.XmlReport)
	src_reports := make(map[string]catch2.XmlReport)

	if err := load(tgt_path, tgt_reports); err != nil {
		log.Fatalln("failed to load target report:", err)
	}
	if err := load(src_path, src_reports); err != nil {
		log.Fatalln("failed to load source report:", err)
	}

	r := DiffReporter{}
	r.Report(os.Stdout, tgt_reports, src_reports)

	return 0
}

func verifyArgs() error {
	if flag.NArg() != 2 {
		return fmt.Errorf("expected 2 arguments but it was %d", flag.NArg())
	}

	return nil
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if err := verifyArgs(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())

		flag.Usage()
	}

	ec := run()
	if ec != 0 {
		os.Exit(ec)
	}
}
