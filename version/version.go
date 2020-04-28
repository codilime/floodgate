package version

import "fmt"

type version struct {
	major   int
	minor   int
	patch    int
	release string
}

func (v version) String() string {
	return fmt.Sprintf("%d.%d.%d-%s", v.major, v.minor, v.patch, v.release)
}

var floodgateVersion = version{
	major:   0,
	minor:   1,
	patch:    0,
	release: "rel",
}

// String get Floodgate version string
func String() string {
	return floodgateVersion.String()
}
