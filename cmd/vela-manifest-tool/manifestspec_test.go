// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
}

func TestManifestSpec_New_Validate(t *testing.T) {
	man := defaultFixture(t)

	assertImageMatch(t, "index.docker.io/octocat/hello-world:latest", man.Image)
	assertImageMatch(t, "index.docker.io/octocat/hello-world:latest-linux-amd64",
		man.Manifests[0].Image)
	assertImageMatch(t, "index.docker.io/octocat/hello-world:latest-linux-arm64-v8",
		man.Manifests[1].Image)

	var data bytes.Buffer

	err := man.Render(&data)
	if err != nil {
		t.Errorf("Error encountered during render: %v", err)
	}

	expected := "image: index.docker.io/octocat/hello-world:latest\n" +
		"manifests:\n" +
		"    - image: index.docker.io/octocat/hello-world:latest-linux-amd64\n" +
		"      platform:\n" +
		"        os: linux\n" +
		"        architecture: amd64\n" +
		"    - image: index.docker.io/octocat/hello-world:latest-linux-arm64-v8\n" +
		"      platform:\n" +
		"        os: linux\n" +
		"        architecture: arm64\n" +
		"        variant: v8\n"

	if data.String() != expected {
		t.Errorf("failed yaml rendering.\nexpected:\n%sactual:\n%s", expected, data.String())
	}
}

func TestManifestSpec_Validations(t *testing.T) {
	testCases := []struct {
		name  string
		valid bool
		avail bool
		ms    *ManifestSpec
	}{
		{
			name:  "missing image",
			valid: false,
			avail: true,
			ms:    trMS(t, func(ms *ManifestSpec) *ManifestSpec { ms.Image = ""; return ms }),
		},
		{
			name:  "missing repo",
			valid: false,
			avail: true,
			ms:    trMS(t, func(ms *ManifestSpec) *ManifestSpec { ms.Image = ""; return ms }),
		},
		{
			name:  "invalid image",
			valid: false,
			avail: false,
			ms: firstMS(NewManifestSpec(trReg(func(r *Registry) *Registry {
				r.Name = ""
				return r
			}), defaultRepo())),
		},
		{
			name:  "invalid reg",
			valid: false,
			avail: false,
			ms: firstMS(NewManifestSpec(defaultRegistry(), trRep(func(r *Repo) *Repo {
				r.Name = ""
				return r
			}))),
		},
		{
			name:  "invalid platform",
			valid: false,
			avail: false,
			ms: firstMS(NewManifestSpec(defaultRegistry(), trRep(func(r *Repo) *Repo {
				r.Platforms = []string{"linux"}
				return r
			}))),
		},
		{
			name:  "no platforms",
			valid: false,
			avail: true,
			ms: firstMS(NewManifestSpec(defaultRegistry(), trRep(func(r *Repo) *Repo {
				r.Platforms = []string{}
				return r
			}))),
		},
		{
			name:  "no tags",
			valid: false,
			avail: false,
			ms: firstMS(NewManifestSpec(defaultRegistry(), trRep(func(r *Repo) *Repo {
				r.Tags = []string{}
				return r
			}))),
		},
		{
			name:  "incomplete template",
			valid: false,
			avail: true,
			ms: firstMS(NewManifestSpec(defaultRegistry(), trRep(func(r *Repo) *Repo {
				r.ComponentTemplate = "{{.Repo}}"
				return r
			}))),
		},
		{
			name:  "invalid tag",
			valid: false,
			avail: true,
			ms: firstMS(NewManifestSpec(defaultRegistry(), trRep(func(r *Repo) *Repo {
				r.Tags = []string{"invalid|tag"}
				return r
			}))),
		},
	}
	for _, tc := range testCases {
		if tc.avail && tc.ms == nil {
			t.Errorf("%s: expected manifest spec to be available, but none returned", tc.name)
		} else if !tc.avail && tc.ms != nil {
			t.Errorf("%s: expected no manifest specs, but at least one returned", tc.name)
		} else if tc.avail && tc.ms != nil {
			err := tc.ms.Validate()
			if err != nil && tc.valid {
				t.Errorf("%s: expected valid ManifestSpec, but got %v", tc.name, err)
			} else if err == nil && !tc.valid {
				t.Errorf("%s: expected invalid ManifestSpec, but got nil", tc.name)
			}
		}
	}
}

func firstMS(ms []*ManifestSpec, _ error) *ManifestSpec {
	if len(ms) > 0 {
		return ms[0]
	}

	return nil
}

func defaultRegistry() *Registry {
	return &Registry{
		Name:      "index.docker.io",
		Username:  "test",
		Password:  "pass",
		PushRetry: 1,
		DryRun:    true,
	}
}

func defaultRepo() *Repo {
	return &Repo{
		Name:              "/octocat/hello-world",
		Tags:              []string{"latest"},
		Platforms:         []string{"linux/amd64", "linux/arm64/v8"},
		ComponentTemplate: "{{.Repo}}:{{.Tag}}-{{.Os}}-{{.Arch}}{{if .Variant}}-{{.Variant}}{{end}}",
	}
}

func defaultFixture(t *testing.T) *ManifestSpec {
	ms, err := NewManifestSpec(defaultRegistry(), defaultRepo())
	if err != nil {
		t.Fatalf("error encountered: %v", err)
	}

	if len(ms) != 1 {
		t.Fatalf("should only have returned a single manifest spec")
	}

	return ms[0]
}

// Translate ManifestSpec.
func trMS(t *testing.T, f func(*ManifestSpec) *ManifestSpec) *ManifestSpec {
	return f(defaultFixture(t))
}

func trReg(f func(r *Registry) *Registry) *Registry {
	return f(defaultRegistry())
}

func trRep(f func(r *Repo) *Repo) *Repo {
	return f(defaultRepo())
}

func assertImageMatch(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("image mismatch\nexpected: %s !=\nactual:   %s", expected, actual)
	}
}
