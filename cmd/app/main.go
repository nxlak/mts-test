package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nxlak/mts-test/internal/analyzer"
	"github.com/nxlak/mts-test/internal/config"
	"github.com/nxlak/mts-test/internal/models"
	"github.com/nxlak/mts-test/internal/render"
)

func main() {
	cfg := config.ParseConfig()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	a := analyzer.NewDefault()

	report, err := a.Analyze(ctx, cfg.RepoURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if !cfg.AllDep {
		report.Updates = filterDirect(report.Updates)
	}

	if err := render.Render(cfg.Format, report); err != nil {
		fmt.Fprintf(os.Stderr, "error: render output: %v\n", err)
		os.Exit(1)
	}
}

func filterDirect(updates []models.Update) []models.Update {
	out := updates[:0]
	for _, u := range updates {
		if !u.Indirect {
			out = append(out, u)
		}
	}
	return out
}
