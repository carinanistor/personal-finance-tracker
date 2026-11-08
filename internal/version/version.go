package version

import (
	"encoding/json"
	"io"
	"runtime/debug"
)

var (
	App     = "atad-project"
	Version = "0.0.0-dev"
	Commit  = ""
	BuiltAt = ""
)

type Info struct {
	App     string `json:"app"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
	BuiltAt string `json:"built_at"`
}

func Current() Info {
	commit := Commit
	if commit == "" {
		if bi, ok := debug.ReadBuildInfo(); ok {
			for _, s := range bi.Settings {
				if s.Key == "vcs.revision" && len(s.Value) >= 7 {
					commit = s.Value[:7]
				}
			}
		}
	}
	if commit == "" {
		commit = "unknown"
	}
	return Info{App: App, Version: Version, Commit: commit, BuiltAt: BuiltAt}
}

func WriteJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(Current())
}
