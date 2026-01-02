package version

import (
	"fmt"
	"runtime"
)

// Info contains version information
type Info struct {
	Version   string
	Commit    string
	BuildDate string
	GoVersion string
	Platform  string
}

// New creates a new version Info struct
func New(version, commit, buildDate string) *Info {
	return &Info{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns a formatted version string
func (v *Info) String() string {
	return fmt.Sprintf(`Version:   %s
Built:     %s
Commit:    %s
Go:        %s
Platform:  %s`,
		v.Version,
		v.BuildDate,
		v.Commit,
		v.GoVersion,
		v.Platform,
	)
}

// Short returns a short version string
func (v *Info) Short() string {
	return fmt.Sprintf("version %s\nBuilt: %s\nCommit: %s",
		v.Version,
		v.BuildDate,
		v.Commit,
	)
}
