package version

import "log"

var (
	GitCommit string
	Version   string
)

func DisplayVersion(appName string) {
	log.Printf("Name: %s\n", appName)
	log.Printf("Commit: %s\n", GitCommit)
	log.Printf("Version: %s\n", Version)
}
