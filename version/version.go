package version

import "fmt"

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func FullVersion() string {
	return fmt.Sprintf("%s (%s @ %s)", Version, Commit, BuildDate)
}
