// SPDX-License-Identifier: Apache-2.0

package main

import "testing"

func TestDocker_Repo_Validate(t *testing.T) {
	// setup types
	r := &Repo{
		Name:      "/target/vela-manifest-tool",
		Tags:      []string{"latest"},
		Platforms: []string{"linux/amd64", "linux/arm64/v8"},
	}

	err := r.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDocker_Repo_Validate_NoName(t *testing.T) {
	// setup types
	r := &Repo{
		Name:      "",
		Tags:      []string{"latest"},
		Platforms: []string{"linux/amd64", "linux/arm64/v8"},
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Repo_Validate_InvalidTags(t *testing.T) {
	// setup types
	r := &Repo{
		Name:      "/target/vela-manifest-tool",
		Tags:      []string{"!@#$%^&*()"},
		Platforms: []string{"linux/amd64", "linux/arm64/v8"},
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Repo_Validate_NoTags(t *testing.T) {
	// setup types
	r := &Repo{
		Name:      "/target/vela-manifest-tool",
		Tags:      []string{},
		Platforms: []string{"linux/amd64", "linux/arm64/v8"},
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Repo_Validate_NoPlatforms(t *testing.T) {
	r := &Repo{
		Name: "/target/vela-manifest-tool",
		Tags: []string{"latest"},
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Repo_InvalidPlatform(t *testing.T) {
	r := &Repo{
		Name:      "/target/vela-manifest-tool",
		Tags:      []string{"latest"},
		Platforms: []string{"windows/riscv64"},
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned an err")
	}
}
