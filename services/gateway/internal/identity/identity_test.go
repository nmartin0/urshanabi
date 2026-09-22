package identity

import (
	"runtime/debug"
	"strings"
	"testing"
	"time"
)

func info(revision, committed string) *debug.BuildInfo {
	return &debug.BuildInfo{Settings: []debug.BuildSetting{
		{Key: "vcs.revision", Value: revision},
		{Key: "vcs.time", Value: committed},
	}}
}

func TestTheIdentityDescribesTheSource(t *testing.T) {
	rev := strings.Repeat("a", 40)
	b, err := FromBuildInfo(info(rev, "2026-09-21T12:00:00Z"), "")
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)
	if b.Component != "gateway" || b.Version != "development" || b.Revision != rev || !b.CommittedAt.Equal(want) {
		t.Fatalf("unexpected identity: %+v", b)
	}
	if b.CommittedAt.Location() != time.UTC {
		t.Fatalf("the commit time is not UTC: %v", b.CommittedAt.Location())
	}
}

func TestABinaryOutsideARepositoryHasNoIdentity(t *testing.T) {
	if _, err := FromBuildInfo(info("", ""), ""); err == nil {
		t.Fatal("an identity without a revision was accepted")
	}
}

func TestOnlyAWellFormedDigestIsKept(t *testing.T) {
	rev := strings.Repeat("a", 40)
	for digest, want := range map[string]string{
		strings.Repeat("b", 64): strings.Repeat("b", 64),
		strings.Repeat("B", 64): "",
		"not-a-digest":          "",
	} {
		b, err := FromBuildInfo(info(rev, "2026-09-21T12:00:00Z"), digest)
		if err != nil || b.ArtifactSHA256 != want {
			t.Errorf("digest %q: got %q, %v", digest, b.ArtifactSHA256, err)
		}
	}
}
