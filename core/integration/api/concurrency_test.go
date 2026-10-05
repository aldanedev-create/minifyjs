package api_test

import (
	"sync"
	"testing"

	"github.com/minifyjs/minifyjs/core/api"
)

func TestAPIConcurrentMinify(t *testing.T) {
	const n = 50
	var wg sync.WaitGroup
	errs := make(chan error, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r, err := api.Minify("const x = 1 + 2;", api.Options{
				MinifyWhitespace: true,
				MinifySyntax:     true,
			})
			if err != nil {
				errs <- err
				return
			}
			if r.HasErrors() {
				errs <- errFromDiags(r)
			}
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Errorf("concurrent run failed: %v", err)
	}
}

func TestAPIConcurrentDifferentInputs(t *testing.T) {
	const n = 30
	var wg sync.WaitGroup
	results := make([]api.Result, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			src := "function f" + itoa(i) + "(x) { return x + " + itoa(i) + "; }"
			r, err := api.Minify(src, api.Options{
				MinifyWhitespace: true,
			})
			if err != nil {
				t.Errorf("minify %d: %v", i, err)
				return
			}
			results[i] = r
		}(i)
	}
	wg.Wait()

	for i, r := range results {
		if r.HasErrors() {
			t.Errorf("result %d has errors: %v", i, r.Diagnostics)
		}
	}
}

func TestAPIConcurrentOptimize(t *testing.T) {
	const n = 20
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = api.Optimize("const x = 1 + 2 + 3;", api.Options{})
		}()
	}
	wg.Wait()
}

func TestAPIConcurrentWithTarget(t *testing.T) {
	const n = 10
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = api.Minify("const f = (x) => x;", api.Options{Target: "es5"})
		}()
	}
	wg.Wait()
}

func TestAPIConcurrentSyntaxError(t *testing.T) {
	const n = 20
	var wg sync.WaitGroup
	count := 0
	var mu sync.Mutex

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, err := api.Minify("function () { } }", api.Options{MinifyWhitespace: true})
			if err != nil {
				return
			}
			if r.HasErrors() {
				mu.Lock()
				count++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if count != n {
		t.Fatalf("expected all %d runs to report errors, got %d", n, count)
	}
}

func TestAPIConcurrentSameInput(t *testing.T) {
	const n = 100
	var wg sync.WaitGroup
	src := "function add(a, b) { return a + b; }"

	var first string
	var mu sync.Mutex

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, _ := api.Minify(src, api.Options{MinifyWhitespace: true})
			mu.Lock()
			if first == "" {
				first = r.Code
			} else if r.Code != first {
				t.Errorf("non-deterministic output: %q vs %q", r.Code, first)
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
}

func errFromDiags(r api.Result) error {
	if len(r.Diagnostics) == 0 {
		return nil
	}
	return &diagErr{msg: r.Diagnostics[0].Message}
}

type diagErr struct{ msg string }

func (e *diagErr) Error() string { return e.msg }

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}