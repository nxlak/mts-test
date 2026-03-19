package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nxlak/mts-test/internal/analyzer"
	"github.com/nxlak/mts-test/internal/models"
	"github.com/nxlak/mts-test/internal/render"
	"github.com/spf13/cobra"
)

var (
	format string
	allDep bool
)

var rootCmd = &cobra.Command{
	Use:   "mts-test <repo-url>",
	Short: "Go Dependency Update Checker",
	Long:  "Check for available updates in Go project dependencies",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(args[0])
	},
}

func init() {
	rootCmd.Flags().StringVarP(&format, "format", "f", "text", "output format: text | json")
	rootCmd.Flags().BoolVarP(&allDep, "all", "a", false, "report direct && indirect dependencies")
}

func Execute() error {
	return rootCmd.Execute()
}

func run(repoURL string) error {
	if format != "text" && format != "json" {
		return fmt.Errorf("unknown format %q; choose text or json", format)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	a := analyzer.NewDefault()

	report, err := a.Analyze(ctx, repoURL)
	if err != nil {
		return err
	}

	if !allDep {
		report.Updates = filterDirect(report.Updates)
	}

	if err := render.Render(format, report); err != nil {
		return fmt.Errorf("render output: %w", err)
	}

	return nil
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
