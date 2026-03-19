package render

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/nxlak/mts-test/internal/models"
)

const (
	colModule    = "MODULE"
	colCurrent   = "CURRENT"
	colLatest    = "LATEST"
	colIndirect  = "INDIRECT"
	colSeparator = "------\t-------\t------\t--------"
)

type jsonReport struct {
	ModuleName string       `json:"module_name"`
	GoVersion  string       `json:"go_version"`
	Updates    []jsonUpdate `json:"updates"`
}

type jsonUpdate struct {
	Path           string `json:"path"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	Indirect       bool   `json:"indirect"`
}

// outputs report in text / json format
func Render(w io.Writer, format string, report *models.Report) error {
	switch format {
	case "json":
		return renderJSON(w, report)
	default:
		return renderText(w, report)
	}
}

func renderText(w io.Writer, r *models.Report) error {
	fmt.Fprintf(w, "Module  : %s\n", r.ModuleName)
	fmt.Fprintf(w, "Go      : %s\n", r.GoVersion)
	fmt.Fprintf(w, "Updates : %d\n\n", len(r.Updates))

	if len(r.Updates) == 0 {
		_, err := fmt.Fprintln(w, "All dependencies have actual versions.")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	defer tw.Flush()

	fmt.Fprintln(tw, colModule+"\t"+colCurrent+"\t"+colLatest+"\t"+colIndirect)
	fmt.Fprintln(tw, colSeparator)
	for _, u := range r.Updates {
		indirect := ""
		if u.Indirect {
			indirect = "yes"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", u.Path, u.CurrentVersion, u.LatestVersion, indirect)
	}

	return nil
}

func renderJSON(w io.Writer, r *models.Report) error {
	out := jsonReport{
		ModuleName: r.ModuleName,
		GoVersion:  r.GoVersion,
	}
	for _, u := range r.Updates {
		out.Updates = append(out.Updates, jsonUpdate{
			Path:           u.Path,
			CurrentVersion: u.CurrentVersion,
			LatestVersion:  u.LatestVersion,
			Indirect:       u.Indirect,
		})
	}
	if out.Updates == nil {
		out.Updates = []jsonUpdate{}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
