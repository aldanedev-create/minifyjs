// Command fixture_runner walks core/testdata/fixtures/ and runs each
// fixture's input.js through MinifyJS's public API, comparing the
// result against the fixture's expected.min.js.
//
// A fixture directory without an expected.min.js is reported as
// "skipped" rather than failed: the blueprint allows fixtures to be
// staged before their expected output has been authored.
//
// Usage:
//
//	go run ./tools/fixture_runner
//	go run ./tools/fixture_runner --json
//	go run ./tools/fixture_runner --verbose
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/minifyjs/minifyjs/core/api"
)

type result struct {
	Name     string `json:"name"`
	Status   string `json:"status"` // "pass", "fail", "skip"
	InputLen int    `json:"input_bytes"`
	WantLen  int    `json:"want_bytes,omitempty"`
	GotLen   int    `json:"got_bytes,omitempty"`
	Message  string `json:"message,omitempty"`
}

type report struct {
	GeneratedAt time.Time `json:"generated_at"`
	FixturesDir string    `json:"fixtures_dir"`
	Passed      int       `json:"passed"`
	Failed      int       `json:"failed"`
	Skipped     int       `json:"skipped"`
	Results     []result  `json:"results"`
}

func main() {
	var (
		asJSON  = flag.Bool("json", false, "print machine-readable JSON instead of a table")
		verbose = flag.Bool("verbose", false, "print per-fixture details even for passes")
		root    = flag.String("dir", "", "fixtures directory (default: auto-detect)")
	)
	flag.Parse()

	dir, err := findFixturesDir(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fixture_runner: %v\n", err)
		os.Exit(1)
	}

	rep, err := runFixtures(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fixture_runner: %v\n", err)
		os.Exit(1)
	}

	if *asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(rep); err != nil {
			fmt.Fprintf(os.Stderr, "fixture_runner: %v\n", err)
			os.Exit(1)
		}
	} else {
		printTable(rep, *verbose)
	}

	if rep.Failed > 0 {
		os.Exit(1)
	}
}

// findFixturesDir returns the fixtures directory, either the explicit
// --dir or the one relative to the current working directory.
func findFixturesDir(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	// The conventional location is core/testdata/fixtures. Try a few
	// common starting points so the tool works from either the repo
	// root or from tools/.
	candidates := []string{
		filepath.Join("core", "testdata", "fixtures"),
		filepath.Join("..", "core", "testdata", "fixtures"),
		filepath.Join("..", "..", "core", "testdata", "fixtures"),
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c, nil
		}
	}
	return "", fmt.Errorf("could not find fixtures directory; pass --dir")
}

func runFixtures(dir string) (report, error) {
	rep := report{
		GeneratedAt: time.Now().UTC(),
		FixturesDir: dir,
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return rep, err
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		r := runOne(filepath.Join(dir, e.Name()), e.Name())
		switch r.Status {
		case "pass":
			rep.Passed++
		case "fail":
			rep.Failed++
		case "skip":
			rep.Skipped++
		}
		rep.Results = append(rep.Results, r)
	}

	return rep, nil
}

func runOne(dir, name string) result {
	r := result{Name: name}

	inputPath := filepath.Join(dir, "input.js")
	expectedPath := filepath.Join(dir, "expected.min.js")

	input, err := os.ReadFile(inputPath)
	if err != nil {
		r.Status = "skip"
		r.Message = "no input.js"
		return r
	}
	r.InputLen = len(input)

	expected, err := os.ReadFile(expectedPath)
	if err != nil {
		r.Status = "skip"
		r.Message = "no expected.min.js"
		return r
	}
	r.WantLen = len(expected)

	got, apiErr := api.Minify(string(input), api.Options{
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
	})
	if apiErr != nil {
		r.Status = "fail"
		r.Message = apiErr.Error()
		return r
	}
	if got.HasErrors() {
		r.Status = "fail"
		var msgs []string
		for _, d := range got.Diagnostics {
			msgs = append(msgs, d.Message)
		}
		r.Message = strings.Join(msgs, "; ")
		return r
	}

	r.GotLen = len(got.Code)
	if got.Code != string(expected) {
		r.Status = "fail"
		r.Message = fmt.Sprintf(
			"output mismatch\n  want: %q\n  got:  %q",
			truncate(string(expected), 120),
			truncate(got.Code, 120),
		)
		return r
	}

	r.Status = "pass"
	return r
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func printTable(rep report, verbose bool) {
	fmt.Printf("fixture runner\n")
	fmt.Printf("  directory: %s\n", rep.FixturesDir)
	fmt.Printf("  passed:    %d\n", rep.Passed)
	fmt.Printf("  failed:    %d\n", rep.Failed)
	fmt.Printf("  skipped:   %d\n", rep.Skipped)
	fmt.Println()

	for _, r := range rep.Results {
		switch r.Status {
		case "pass":
			if verbose {
				fmt.Printf("  PASS  %-30s %d B -> %d B\n", r.Name, r.InputLen, r.GotLen)
			}
		case "skip":
			fmt.Printf("  SKIP  %-30s %s\n", r.Name, r.Message)
		case "fail":
			fmt.Printf("  FAIL  %-30s %s\n", r.Name, r.Message)
		}
	}
}