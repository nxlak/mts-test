package render

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/nxlak/mts-test/internal/models"
)

type errWriter struct{}

func (errWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestRender_JSON(t *testing.T) {
	report := &models.Report{
		ModuleName: "github.com/nxlak/mts-test",
		GoVersion:  "1.24",
		Updates: []models.Update{
			{
				Dependency: models.Dependency{
					Path:           "github.com/stretchr/testify",
					CurrentVersion: "v1.8.0",
					Indirect:       false,
				},
				LatestVersion: "v1.9.0",
			},
			{
				Dependency: models.Dependency{
					Path:           "github.com/sirupsen/logrus",
					CurrentVersion: "v1.8.1",
					Indirect:       true,
				},
				LatestVersion: "v1.9.3",
			},
		},
	}

	var buf bytes.Buffer

	if err := Render(&buf, "json", report); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	var got jsonReport
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal json output: %v", err)
	}

	if got.ModuleName != report.ModuleName {
		t.Fatalf("module_name = %q; want %q", got.ModuleName, report.ModuleName)
	}
	if got.GoVersion != report.GoVersion {
		t.Fatalf("go_version = %q; want %q", got.GoVersion, report.GoVersion)
	}
	if len(got.Updates) != 2 {
		t.Fatalf("len(updates) = %d; want 2", len(got.Updates))
	}

	if got.Updates[0].Path != "github.com/stretchr/testify" {
		t.Errorf("updates[0].path = %q; want %q", got.Updates[0].Path, "github.com/stretchr/testify")
	}
	if got.Updates[0].CurrentVersion != "v1.8.0" {
		t.Errorf("updates[0].current_version = %q; want %q", got.Updates[0].CurrentVersion, "v1.8.0")
	}
	if got.Updates[0].LatestVersion != "v1.9.0" {
		t.Errorf("updates[0].latest_version = %q; want %q", got.Updates[0].LatestVersion, "v1.9.0")
	}
	if got.Updates[0].Indirect {
		t.Errorf("updates[0].indirect = %v; want false", got.Updates[0].Indirect)
	}

	if got.Updates[1].Path != "github.com/sirupsen/logrus" {
		t.Errorf("updates[1].path = %q; want %q", got.Updates[1].Path, "github.com/sirupsen/logrus")
	}
	if got.Updates[1].CurrentVersion != "v1.8.1" {
		t.Errorf("updates[1].current_version = %q; want %q", got.Updates[1].CurrentVersion, "v1.8.1")
	}
	if got.Updates[1].LatestVersion != "v1.9.3" {
		t.Errorf("updates[1].latest_version = %q; want %q", got.Updates[1].LatestVersion, "v1.9.3")
	}
	if !got.Updates[1].Indirect {
		t.Errorf("updates[1].indirect = %v; want true", got.Updates[1].Indirect)
	}
}

func TestRender_JSON_WriteError(t *testing.T) {
	report := &models.Report{
		ModuleName: "github.com/nxlak/mts-test",
		GoVersion:  "1.24",
	}

	err := Render(errWriter{}, "json", report)
	if err == nil {
		t.Fatal("Render() error = nil; want error")
	}
}

func TestRender_TextWithUpdates(t *testing.T) {
	report := &models.Report{
		ModuleName: "github.com/nxlak/mts-test",
		GoVersion:  "1.24",
		Updates: []models.Update{
			{
				Dependency: models.Dependency{
					Path:           "github.com/stretchr/testify",
					CurrentVersion: "v1.8.0",
					Indirect:       false,
				},
				LatestVersion: "v1.9.0",
			},
			{
				Dependency: models.Dependency{
					Path:           "github.com/sirupsen/logrus",
					CurrentVersion: "v1.8.1",
					Indirect:       true,
				},
				LatestVersion: "v1.9.3",
			},
		},
	}

	var buf bytes.Buffer

	if err := Render(&buf, "text", report); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	out := buf.String()

	checks := []string{
		"Module  : github.com/nxlak/mts-test",
		"Go      : 1.24",
		"Updates : 2",
		"MODULE",
		"CURRENT",
		"LATEST",
		"INDIRECT",
		"github.com/stretchr/testify",
		"v1.8.0",
		"v1.9.0",
		"github.com/sirupsen/logrus",
		"v1.8.1",
		"v1.9.3",
		"yes",
	}

	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Errorf("text output does not contain %q\noutput:\n%s", want, out)
		}
	}
}

func TestRender_TextNoUpdates(t *testing.T) {
	report := &models.Report{
		ModuleName: "github.com/nxlak/mts-test",
		GoVersion:  "1.24",
		Updates:    nil,
	}

	var buf bytes.Buffer

	if err := Render(&buf, "text", report); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	out := buf.String()

	checks := []string{
		"Module  : github.com/nxlak/mts-test",
		"Go      : 1.24",
		"Updates : 0",
		"All dependencies have actual versions.",
	}

	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Errorf("text output does not contain %q\noutput:\n%s", want, out)
		}
	}
}
