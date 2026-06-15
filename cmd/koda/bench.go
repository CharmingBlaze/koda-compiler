package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"koda/api"
	"koda/internal/nativebuild"
)

type benchOptions struct {
	count    int
	warmup   int
	noOpt    bool
	debug    bool
	src      string
	progArgs []string
}

func parseBenchArgs(args []string) (benchOptions, error) {
	o := benchOptions{count: 5, warmup: 1}
	before, after := splitAtDoubleDash(args)
	o.progArgs = after
	for i := 0; i < len(before); i++ {
		switch before[i] {
		case "--count":
			if i+1 >= len(before) {
				return o, fmt.Errorf("--count requires a number")
			}
			i++
			n, err := parsePositiveInt(before[i])
			if err != nil {
				return o, fmt.Errorf("--count: %w", err)
			}
			o.count = n
		case "--warmup":
			if i+1 >= len(before) {
				return o, fmt.Errorf("--warmup requires a number")
			}
			i++
			n, err := parsePositiveInt(before[i])
			if err != nil {
				return o, fmt.Errorf("--warmup: %w", err)
			}
			o.warmup = n
		case "--no-opt":
			o.noOpt = true
		case "--debug":
			o.debug = true
		default:
			if stringsHasPrefix(before[i], "-") {
				return o, fmt.Errorf("unknown flag: %s", before[i])
			}
			if o.src != "" {
				return o, fmt.Errorf("multiple source files")
			}
			o.src = before[i]
		}
	}
	if o.src == "" {
		return o, fmt.Errorf("usage: koda bench [--count N] [--warmup N] [--no-opt] [--debug] <file.koda> [-- <args...>]")
	}
	return o, nil
}

func stringsHasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func parsePositiveInt(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid number %q", s)
		}
		n = n*10 + int(c-'0')
		if n <= 0 {
			return 0, fmt.Errorf("must be positive")
		}
	}
	return n, nil
}

func runBench(args []string) error {
	o, err := parseBenchArgs(args)
	if err != nil {
		return err
	}
	opts := nativebuild.BuildOptions{NoOpt: o.noOpt || o.debug, Debug: o.debug}
	return withProject(o.src, func() error {
		return benchFile(o.src, opts, o.warmup, o.count, o.progArgs)
	})
}

func benchFile(path string, opts nativebuild.BuildOptions, warmup, count int, progArgs []string) error {
	timings := make([]time.Duration, 0, count)
	runs := warmup + count
	for i := 0; i < runs; i++ {
		start := time.Now()
		err := api.RunWithWritersOptsProgram(path, "", os.Stdout, os.Stderr, opts, progArgs)
		elapsed := time.Since(start)
		if err != nil {
			return err
		}
		if i >= warmup {
			timings = append(timings, elapsed)
			fmt.Fprintf(os.Stderr, "  run %d: %.3f ms\n", i-warmup+1, elapsed.Seconds()*1000)
		}
	}
	if len(timings) == 0 {
		return fmt.Errorf("no benchmark samples collected")
	}
	sorted := append([]time.Duration(nil), timings...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	var total time.Duration
	for _, t := range timings {
		total += t
	}
	avg := total / time.Duration(len(timings))
	p50 := percentile(sorted, 0.50)
	p95 := percentile(sorted, 0.95)
	p99 := percentile(sorted, 0.99)
	fmt.Printf("bench: %d runs (warmup %d)\n", count, warmup)
	fmt.Printf("  avg  %.3f ms\n", avg.Seconds()*1000)
	fmt.Printf("  p50  %.3f ms\n", p50.Seconds()*1000)
	fmt.Printf("  p95  %.3f ms\n", p95.Seconds()*1000)
	fmt.Printf("  p99  %.3f ms\n", p99.Seconds()*1000)
	fmt.Printf("  total %.3f ms\n", total.Seconds()*1000)
	return nil
}

func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	if p <= 0 {
		return sorted[0]
	}
	if p >= 1 {
		return sorted[len(sorted)-1]
	}
	idx := int(float64(len(sorted)-1) * p)
	return sorted[idx]
}
