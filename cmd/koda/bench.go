package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"koda/api"
	"koda/internal/nativebuild"
)

type benchOptions struct {
	count   int
	warmup  int
	noOpt   bool
	debug   bool
	src     string
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
			if strings.HasPrefix(before[i], "-") {
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
	var total time.Duration
	runs := warmup + count
	for i := 0; i < runs; i++ {
		start := time.Now()
		err := api.RunWithWritersOptsProgram(path, "", os.Stdout, os.Stderr, opts, progArgs)
		elapsed := time.Since(start)
		if err != nil {
			return err
		}
		if i >= warmup {
			total += elapsed
			fmt.Fprintf(os.Stderr, "  run %d: %.3f ms\n", i-warmup+1, elapsed.Seconds()*1000)
		}
	}
	avg := total / time.Duration(count)
	fmt.Printf("bench: %d runs (warmup %d), avg %.3f ms, total %.3f ms\n",
		count, warmup, avg.Seconds()*1000, total.Seconds()*1000)
	return nil
}
