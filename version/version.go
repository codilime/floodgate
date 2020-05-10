package version

import "fmt"

// GitCommit string
var GitCommit string

// Release string
var Release string

// GateVersion string
var GateVersion string

// BuiltDate string
var BuiltDate string

// GoVersion string
var GoVersion string

type version struct {
	major int
	minor int
	patch int
}

func (v version) BuildInfo() string {
	gateAPIInfo := "\nGate API version:\t" + GateVersion
	goInfo := "\nGo version:\t\t" + GoVersion
	gitCommitInfo := "\nGit commit:\t\t" + GitCommit
	builtDateInfo := "\nBuilt:\t\t\t" + BuiltDate + "\n"

	return "Version:\t\t" + v.Short() + gateAPIInfo + goInfo + gitCommitInfo + builtDateInfo

}

func (v version) Short() string {
	versionInfo := fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
	if Release != "" {
		versionInfo = versionInfo + "-" + Release
	}

	return versionInfo
}

var floodgateVersion = version{
	major: 0,
	minor: 2,
	patch: 0,
}

// Short get Floodgate version string
func Short() string {
	return floodgateVersion.Short()
}

// BuildInfo get Floodgate build info
func BuildInfo() string {
	return floodgateVersion.BuildInfo()
}
