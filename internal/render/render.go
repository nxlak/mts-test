package render

import (
	"encoding/json"
	"fmt"
	"os"
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

// outputs report in text / json format
func Render(format string, report *models.Report) error {
	switch format {
	case "json":
		return renderJSON(report)
	default:
		return renderText(report)
	}
}

func renderText(r *models.Report) error {
	fmt.Printf("Module  : %s\n", r.ModuleName)
	fmt.Printf("Go      : %s\n", r.GoVersion)
	fmt.Printf("Updates : %d\n\n", len(r.Updates))

	if len(r.Updates) == 0 {
		fmt.Println("All dependencies are up to date.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, colModule+"\t"+colCurrent+"\t"+colLatest+"\t"+colIndirect)
	fmt.Fprintln(w, colSeparator)
	for _, u := range r.Updates {
		indirect := ""
		if u.Indirect {
			indirect = "yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", u.Path, u.CurrentVersion, u.LatestVersion, indirect)
	}

	return nil
}

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

func renderJSON(r *models.Report) error {
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

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
