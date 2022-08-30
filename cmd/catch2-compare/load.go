package main

import (
	"encoding/xml"
	"os"
	"path/filepath"

	"github.com/lesomnus/catch2-compare/internal/catch2"
)

func load(path string, reports map[string]catch2.Report) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				// No recursive traverse.
				continue
			}

			if err := load(filepath.Join(path, entry.Name()), reports); err != nil {
				return err
			}
		}

		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var report catch2.Report
	if err := xml.Unmarshal(data, &report); err != nil {
		return err
	}

	reports[report.Name] = report
	return nil
}
