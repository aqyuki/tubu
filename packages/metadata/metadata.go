// metadata package provides methods to get application metadata. e.g. version, build time, etc.
package metadata

import "sync"

var (
	once     sync.Once
	metadata *Metadata
)

var (
	Version    string
	GoVersion  string
	BuildDate  string
	CommitHash string
)

// Metadata struct provides methods to get application metadata.
type Metadata struct {
	version    string
	goVersion  string
	buildDate  string
	commitHash string
}

func GetMetadata() *Metadata {
	return metadata
}

func (m *Metadata) Version() string {
	return m.version
}

func (m *Metadata) GoVersion() string {
	return m.goVersion
}

func (m *Metadata) BuildDate() string {
	return m.buildDate
}

func (m *Metadata) CommitHash() string {
	return m.commitHash
}

// buildMetadata initialize metadata with version, go version, build date and commit hash.
func buildMetadata() *Metadata {
	return &Metadata{
		version:    Version,
		goVersion:  GoVersion,
		buildDate:  BuildDate,
		commitHash: CommitHash,
	}
}

// initialize metadata once when the package is imported
func init() {
	once.Do(func() {
		metadata = buildMetadata()
	})
}
