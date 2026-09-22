// Package identity reports this build's identity (roadmap R-121). The
// source revision and commit time come from the build information the
// toolchain records in every binary built from a source repository, so
// two builds of one commit report one identity.
package identity

import (
	"errors"
	"os"
	"regexp"
	"runtime/debug"
	"time"
)

// Component is this component's directory name.
const Component = "gateway"

// ArtifactVariable names the environment variable a deployment sets to the
// published SHA-256 of the running artifact; a binary cannot contain its
// own digest.
const ArtifactVariable = "URSHANABI_ARTIFACT_SHA256"

// version is "development" unless a release build sets it with -ldflags -X.
var version = "development"

var sha256Hex = regexp.MustCompile(`^[0-9a-f]{64}$`)

// Build is one build's identity.
type Build struct {
	Component      string
	Version        string
	Revision       string
	CommittedAt    time.Time
	ArtifactSHA256 string
}

// Current returns the running binary's identity.
func Current() (Build, error) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return Build{}, errors.New("the binary carries no build information")
	}
	return FromBuildInfo(info, os.Getenv(ArtifactVariable))
}

// FromBuildInfo reads an identity from recorded build information. The
// artifact digest is kept only if it is a well-formed SHA-256.
func FromBuildInfo(info *debug.BuildInfo, artifact string) (Build, error) {
	settings := map[string]string{}
	for _, s := range info.Settings {
		settings[s.Key] = s.Value
	}
	revision := settings["vcs.revision"]
	if len(revision) != 40 {
		return Build{}, errors.New("the binary was not built from a source repository")
	}
	committed, err := time.Parse(time.RFC3339, settings["vcs.time"])
	if err != nil {
		return Build{}, errors.New("the binary records no commit time")
	}
	if !sha256Hex.MatchString(artifact) {
		artifact = ""
	}
	return Build{
		Component:      Component,
		Version:        version,
		Revision:       revision,
		CommittedAt:    committed.UTC(),
		ArtifactSHA256: artifact,
	}, nil
}
