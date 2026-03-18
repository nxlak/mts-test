package models

// data from go.mod
type ModuleInfo struct {
	Name         string
	GoVersion    string
	Dependencies []Dependency
}

type Dependency struct {
	Path           string
	CurrentVersion string
	Indirect       bool
}

type Update struct {
	Dependency
	LatestVersion string
}

type Report struct {
	ModuleName string
	GoVersion  string
	Updates    []Update
}
