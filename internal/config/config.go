package config

import (
	"flag"
	"fmt"
	"os"
)

const usageText = `Go Dependency Update Checker

Usage:
  mts-test [flags] <repo-url>

Arguments:
  <repo-url>   URL of a git repository (e.g. https://github.com/foo/bar or github.com/foo/bar)

Flags:
`

type Config struct {
	RepoURL string
	Format  string
	AllDep  bool
}

func ParseConfig() Config {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usageText)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
	}

	format := flag.String("format", "text", "output format: text | json")
	all := flag.Bool("all", false, "report direct && indirect dependencies")

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	if *format != "text" && *format != "json" {
		fmt.Fprintf(os.Stderr, "error: unknown format %q; choose text or json\n", *format)
		os.Exit(2)
	}

	return Config{
		RepoURL: flag.Arg(0),
		Format:  *format,
		AllDep:  *all,
	}
}
