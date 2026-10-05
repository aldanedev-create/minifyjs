// Command benchmark_report reads every JSON file in
// bench/results/history/ and prints the delta between the two most
// recent runs. Regressions larger than a configurable threshold are
// flagged.
//
// Usage:
//
//	go run ./tools/benchmark_report
//	go run ./tools/benchmark_report --threshold 0.10
//	go run ./tools/benchmark_report --dir bench/results/history
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// historyFile is the schema of one entry in bench/results/history/.
type historyFile struct {
	GeneratedAt    time.Time           `json:"generated_at"`
	MinifyJSVersion string             `json:"minifyjs_version"`
	EsbuildVersion string              `json:"esbuild_version"`
	Platform       string              `json:"platform"`
	PythonVersion  string              `json:"python"`
	Results        map[string]any      `json:"results"`
}

func main() {
	var (
		dir       = flag.String("dir", "", "history directory")
		threshold = flag.Float64("threshold", 0.05, "regression threshold (0.05 = 5%)")
	)
	flag.Parse()

	historyDir := *dir
	if historyDir == "" {
		candidates := []string{
			filepath.Join("bench", "results", "history"),
			filepath.Join("..", "bench", "results", "history"),
			filepath.Join("..", "..", "bench", "results", "history"),
		}
		for _, c := range candidates {
			if info, err := os.Stat(c); err == nil && info.IsDir() {
				historyDir = c
				break
			}
		}
	}
	if historyDir == "" {
		fmt.Fprintln(os.Stderr, "benchmark_report: could not find history directory")
		os.Exit(1)
	}

	files, err := loadHistory(historyDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "benchmark_report: %v\n", err)
		os.Exit(1)
	}
	if len(files) < 2 {
		fmt.Fprintf(os.Stderr, "benchmark_report: need at least 2 runs, found %d\n", len(files))
		os.Exit(1)
	}

	// Files are sorted oldest first.
	prev := files[len(files)-2]
	curr := files[len(files)-1]

	fmt.Printf("comparing runs\n")
	fmt.Printf("  older: %s (minifyjs %s)\n",
		prev.GeneratedAt.Format(time.RFC3339), prev.MinifyJSVersion)
	fmt.Printf("  newer: %s (minifyjs %s)\n\n",
		curr.GeneratedAt.Format(time.RFC3339), curr.MinifyJSVersion)

	regressions := compareSection("compression", prev.Results, curr.Results, *threshold, true)
	regressions += compareSection("throughput", prev.Results, curr.Results, *threshold, false)
	regressions += compareSection("memory", prev.Results, curr.Results, *threshold, true)

	if regressions > 0 {
		fmt.Printf("\n%d regression(s) above threshold %.1f%%\n",
			regressions, *threshold*100)
		os.Exit(1)
	}
	fmt.Println("\nno regressions above threshold")
}

func loadHistory(dir string) ([]historyFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []historyFile
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		var h historyFile
		if err := json.Unmarshal(data, &h); err != nil {
			return nil, fmt.Errorf("%s: %w", e.Name(), err)
		}
		files = append(files, h)
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].GeneratedAt.Before(files[j].GeneratedAt)
	})
	return files, nil
}

// compareSection walks a top-level results section (e.g.
// "compression") and prints the delta between prev and curr.
func compareSection(name string, prev, curr map[string]any, threshold float64, lowerIsBetter bool) int {
	prevSection, ok := prev[name].(map[string]any)
	if !ok {
		return 0
	}
	currSection, ok := curr[name].(map[string]any)
	if !ok {
		return 0
	}

	fmt.Printf("== %s\n", name)
	regressions := 0

	// Compare at the top level of the section. The exact shape of
	// each section is tool-specific, so we only walk one level deep
	// and look for float64 values.
	for key, prevVal := range prevSection {
		currVal, ok := currSection[key]
		if !ok {
			continue
		}
		prevF, ok1 := toFloat(prevVal)
		currF, ok2 := toFloat(currVal)
		if !ok1 || !ok2 || prevF == 0 {
			continue
		}
		delta := (currF - prevF) / prevF
		flagged := isRegression(delta, threshold, lowerIsBetter)
		marker := ""
		if flagged {
			marker = "  <-- REGRESSION"
			regressions++
		}
		fmt.Printf("  %-30s %10.3f -> %10.3f  (%+.1f%%)%s\n",
			key, prevF, currF, delta*100, marker)
	}
	fmt.Println()
	return regressions
}

func toFloat(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case int:
		return float64(x), true
	}
	return 0, false
}

func isRegression(delta float64, threshold float64, lowerIsBetter bool) bool {
	if lowerIsBetter {
		return delta > threshold
	}
	return -delta > threshold
}