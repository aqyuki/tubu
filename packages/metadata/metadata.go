// metadata package provides methods to get application metadata. e.g. version, build time, etc.
package metadata

import "sync"

var (
	once     sync.Once
	metadata *Metadata
)

var (
	Version    string = "development"
	GoVersion  string = "unknown"
	BuildDate  string = "unknown"
	CommitHash string = "unknown"
)

// Metadata struct provides methods to get application metadata.
type Metadata struct {
	Version    string
	GoVersion  string
	BuildDate  string
	CommitHash string
}

func GetMetadata() *Metadata {
	return metadata
}

// buildMetadata initialize metadata with version, go version, build date and commit hash.
func buildMetadata() *Metadata {
	return &Metadata{
		Version:    Version,
		GoVersion:  GoVersion,
		BuildDate:  BuildDate,
		CommitHash: CommitHash,
	}
}

// initialize metadata once when the package is imported
func init() {
	once.Do(func() {
		metadata = buildMetadata()
	})
}
