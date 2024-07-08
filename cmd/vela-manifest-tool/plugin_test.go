// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"testing"
)

func TestPluginScenarios(t *testing.T) {
	testCases := []struct {
		name  string
		valid bool
		p     *Plugin
	}{
		{
			name:  "all populated",
			valid: true,
			p:     makeDefaultPlugin(),
		},
		{
			name:  "invalid template",
			valid: false,
			p:     trP(func(p *Plugin) *Plugin { p.Repo.ComponentTemplate = "{{"; return p }),
		},
		{
			name:  "only one platform component",
			valid: false,
			p:     trP(func(p *Plugin) *Plugin { p.Repo.Platforms = []string{"linux"}; return p }),
		},
		{
			name:  "no image provided",
			valid: false,
			p:     trP(func(p *Plugin) *Plugin { p.Repo.Name = ""; return p }),
		},
		{
			name:  "no tags provided",
			valid: false,
			p:     trP(func(p *Plugin) *Plugin { p.Repo.Tags = []string{}; return p }),
		},
		{
			name:  "no registry name",
			valid: false,
			p:     trP(func(p *Plugin) *Plugin { p.Registry.Name = ""; return p }),
		},
		{
			name:  "invalid template variable",
			valid: false,
			p: trP(func(p *Plugin) *Plugin {
				p.Repo.ComponentTemplate = "{{.Res}}"
				return p
			}),
		},
	}
	for _, tc := range testCases {
		err := tc.p.Validate()
		if err != nil && tc.valid {
			t.Errorf("%s: expected valid plugin, but error was %v", tc.name, err)
		} else if err == nil && !tc.valid {
			t.Errorf("%s: expected invalid plugin, but error was nil", tc.name)
		} else {
			fmt.Printf("%s: Completed successfully: %v\n", tc.name, err)
		}
	}
}

func makeDefaultPlugin() *Plugin {
	return &Plugin{
		// build configuration
		// registry configuration
		Registry: &Registry{
			DryRun:    true,
			Name:      "registry.example.com",
			Username:  "docker_user",
			Password:  "docker_pass",
			PushRetry: 1,
		},
		// repo configuration
		Repo: &Repo{
			Name:              "project/image",
			Tags:              []string{"latest", "v0.0.0"},
			Platforms:         []string{"linux/amd64", "linux/arm64/v8"},
			ComponentTemplate: "{{.Repo}}:{{.Tag}}-{{.Os}}-{{.Arch}}{{if .Variant}}-{{.Variant}}{{end}}",
		},
	}
}

// Translate Plugin.
func trP(t func(*Plugin) *Plugin) *Plugin {
	return t(makeDefaultPlugin())
}
