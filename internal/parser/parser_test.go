package parser

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleGoMod = `module github.com/nxlak/mts-test

go 1.22

require (
	github.com/stretchr/testify v1.8.0
	golang.org/x/mod v0.12.0 // indirect
)
`

func TestParser_Parse(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(sampleGoMod), 0o600); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	p := NewParser()
	info, err := p.Parse(dir)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if info.Name != "github.com/nxlak/mts-test" {
		t.Errorf("Name = %q; want %q", info.Name, "github.com/nxlak/mts-test")
	}
	if info.GoVersion != "1.22" {
		t.Errorf("GoVersion = %q; want %q", info.GoVersion, "1.22")
	}
	if len(info.Dependencies) != 2 {
		t.Fatalf("len(Dependencies) = %d; want 2", len(info.Dependencies))
	}

	direct := info.Dependencies[0]
	if direct.Path != "github.com/stretchr/testify" {
		t.Errorf("dep[0].Path = %q; want github.com/stretchr/testify", direct.Path)
	}
	if direct.CurrentVersion != "v1.8.0" {
		t.Errorf("dep[0].CurrentVersion = %q; want v1.8.0", direct.CurrentVersion)
	}
	if direct.Indirect {
		t.Errorf("dep[0].Indirect = true; want false")
	}

	indirect := info.Dependencies[1]
	if !indirect.Indirect {
		t.Errorf("dep[1].Indirect = false; want true")
	}
}

func TestParser_Parse_MissingGoMod(t *testing.T) {
	p := NewParser()
	_, err := p.Parse(t.TempDir())
	if err == nil {
		t.Fatal("expected error for missing go.mod, got nil")
	}
}

func TestParser_Parse_InvalidGoMod(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("not valid go.mod content !!!"), 0o600); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	p := NewParser()
	_, err := p.Parse(dir)
	if err == nil {
		t.Fatal("expected error for invalid go.mod, got nil")
	}
}
