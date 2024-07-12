// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"

	"github.com/go-vela/vela-manifest-tool/version"
)

func TestVersionCompatible(t *testing.T) {
	v := version.New()
	if v == nil {
		t.Error("version.New should return a value")
	}
}

func TestVersionSemver(t *testing.T) {
	version.Tag = "abcd"

	v := version.New()
	if v != nil {
		t.Errorf("version.New should return nil if a non-semver Tag (%q) is provided", version.Tag)
	}
}
