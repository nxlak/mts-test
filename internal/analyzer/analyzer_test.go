package analyzer

import (
	"context"
	"errors"
	"testing"

	"github.com/nxlak/mts-test/internal/models"
)

type fakeCloner struct {
	dir string
	err error
}

func (f *fakeCloner) Clone(_ context.Context, _ string) (string, func(), error) {
	return f.dir, func() {}, f.err
}

type fakeParser struct {
	info *models.ModuleInfo
	err  error
}

func (f *fakeParser) Parse(_ string) (*models.ModuleInfo, error) {
	return f.info, f.err
}

type fakeChecker struct {
	updates []models.Update
	err     error
}

func (f *fakeChecker) FindUpdates(_ context.Context, _ []models.Dependency) ([]models.Update, error) {
	return f.updates, f.err
}

func TestAnalyzer_Success(t *testing.T) {
	a := NewAnalyzer(
		&fakeCloner{dir: "/tmp/fake"},
		&fakeParser{info: &models.ModuleInfo{
			Name:      "github.com/nxlak/test",
			GoVersion: "1.22",
			Dependencies: []models.Dependency{
				{Path: "github.com/stretchr/testify", CurrentVersion: "v1.8.0"},
			},
		}},
		&fakeChecker{updates: []models.Update{
			{
				Dependency:    models.Dependency{Path: "github.com/stretchr/testify", CurrentVersion: "v1.8.0"},
				LatestVersion: "v1.9.0",
			},
		}},
	)

	report, err := a.Analyze(context.Background(), "https://github.com/nxlak/test")
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if report.ModuleName != "github.com/nxlak/test" {
		t.Errorf("ModuleName = %q; want github.com/nxlak/test", report.ModuleName)
	}
	if report.GoVersion != "1.22" {
		t.Errorf("GoVersion = %q; want 1.22", report.GoVersion)
	}
	if len(report.Updates) != 1 {
		t.Fatalf("len(Updates) = %d; want 1", len(report.Updates))
	}
}

func TestAnalyzer_ClonerError(t *testing.T) {
	a := NewAnalyzer(
		&fakeCloner{err: errors.New("clone failed")},
		&fakeParser{},
		&fakeChecker{},
	)
	_, err := a.Analyze(context.Background(), "https://github.com/nxlak/test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAnalyzer_ParserError(t *testing.T) {
	a := NewAnalyzer(
		&fakeCloner{dir: "/tmp/fake"},
		&fakeParser{err: errors.New("no go.mod")},
		&fakeChecker{},
	)
	_, err := a.Analyze(context.Background(), "https://github.com/nxlak/invalid")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAnalyzer_CheckerError(t *testing.T) {
	a := NewAnalyzer(
		&fakeCloner{dir: "/tmp/fake"},
		&fakeParser{info: &models.ModuleInfo{Name: "github.com/nxlak/checker-test", GoVersion: "1.21"}},
		&fakeChecker{err: errors.New("checker down")},
	)
	_, err := a.Analyze(context.Background(), "https://github.com/nxlak/checker-test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
