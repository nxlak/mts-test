package analyzer

import (
	"context"
	"fmt"

	"github.com/nxlak/mts-test/internal/checker"
	"github.com/nxlak/mts-test/internal/cloner"
	"github.com/nxlak/mts-test/internal/models"
	"github.com/nxlak/mts-test/internal/parser"
)

// all stages of analysis: cloning, parsing, and checking
type Analyzer interface {
	Analyze(ctx context.Context, repoURL string) (*models.Report, error)
}

type analyzer struct {
	cloner  cloner.Cloner
	parser  parser.Parser
	checker checker.Checker
}

func New(cloner cloner.Cloner, parser parser.Parser, checker checker.Checker) Analyzer {
	return &analyzer{
		cloner:  cloner,
		parser:  parser,
		checker: checker,
	}
}

func NewDefault() Analyzer {
	return New(
		cloner.NewCloner(),
		parser.NewParser(),
		checker.NewChecker(),
	)
}

func (a *analyzer) Analyze(ctx context.Context, repoURL string) (*models.Report, error) {
	dir, clean, err := a.cloner.Clone(ctx, repoURL)
	if err != nil {
		return nil, fmt.Errorf("analyzer: clone: %w", err)
	}
	defer clean()

	info, err := a.parser.Parse(dir)
	if err != nil {
		return nil, fmt.Errorf("analyzer: parse: %w", err)
	}

	updates, err := a.checker.FindUpdates(ctx, info.Dependencies)
	if err != nil {
		return nil, fmt.Errorf("analyzer: check updates: %w", err)
	}

	return &models.Report{
		ModuleName: info.Name,
		GoVersion:  info.GoVersion,
		Updates:    updates,
	}, nil
}
