package version

import "fmt"

type version struct {
	major   int
	minor   int
	path    int
	release string
}

func (v version) String() string {
	return fmt.Sprintf("%d.%d.%d-%s", v.major, v.minor, v.path, v.release)
}

var floodgateVersion = version{
	major:   0,
	minor:   1,
	path:    0,
	release: "rel",
}

// String get Floodgate version string
func String() string {
	return floodgateVersion.String()
}
