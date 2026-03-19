package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nxlak/mts-test/internal/models"
	"golang.org/x/mod/modfile"
)

// extracts info from go.mod
type Parser interface {
	Parse(dir string) (*models.ModuleInfo, error)
}

type parser struct{}

func NewParser() Parser {
	return &parser{}
}

func (p *parser) Parse(dir string) (*models.ModuleInfo, error) {
	goModPath := filepath.Join(dir, "go.mod")

	data, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, fmt.Errorf("gomod: read %s: %w", goModPath, err)
	}

	f, err := modfile.Parse(goModPath, data, nil)
	if err != nil {
		return nil, fmt.Errorf("gomod: parse %s: %w", goModPath, err)
	}

	if f.Module == nil || f.Module.Mod.Path == "" {
		return nil, fmt.Errorf("gomod: %s: missing module directive", goModPath)
	}

	info := &models.ModuleInfo{
		Name:      f.Module.Mod.Path,
		GoVersion: goVersion(f),
	}

	for _, req := range f.Require {
		info.Dependencies = append(info.Dependencies, models.Dependency{
			Path:           req.Mod.Path,
			CurrentVersion: req.Mod.Version,
			Indirect:       req.Indirect,
		})
	}

	return info, nil
}

func goVersion(f *modfile.File) string {
	if f.Go != nil {
		return f.Go.Version
	}
	return "unknown"
}
