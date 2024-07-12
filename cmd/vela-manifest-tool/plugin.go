// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/spf13/afero"

	"github.com/sirupsen/logrus"
)

var (
	appFS = afero.NewOsFs()

	// regular expression to validate docker tags
	// refs:
	//  - https://docs.docker.com/engine/reference/commandline/tag/#extended-description
	//  - https://github.com/distribution/distribution/blob/01f589cf8726565aa3c5c053be12873bafedbedc/reference/regexp.go#L41
	tagRegexp = regexp.MustCompile(`^[\w][\w.-]{0,127}$`)
)

// errTagValidation defines the error message
// when the provided tag is not allowed.
const errTagValidation = "tag '%s' not allowed - see https://docs.docker.com/engine/reference/commandline/tag/#extended-description"

// Plugin represents the configuration loaded for the plugin.
type Plugin struct {
	Registry      *Registry       // registry arguments loaded for the plugin
	Repo          *Repo           // repo arguments loaded for the plugin
	manifestSpecs []*ManifestSpec // Parsed specs, populated as side effect of validate
}

// Formats and outputs the command necessary for manifest-tool to build
// and publish a Docker Manifest List or OCI Image Index.
func (p *Plugin) Command(specFile string) *exec.Cmd {
	logrus.Debug("creating manifest-tool command from plugin configuration")

	// variable to store flags for command
	flags := []string{
		"push",
		"from-spec",
		specFile,
	}

	return exec.Command(manifestToolBin, flags...)
}

// Exec formats and runs the commands for building and publishing a Docker image.
func (p *Plugin) Exec() error {
	logrus.Debug("running plugin with provided configuration")

	if len(p.manifestSpecs) == 0 {
		return errors.New("no manifest specs")
	}

	// create registry file for authentication
	err := p.Registry.Write()
	if err != nil {
		return err
	}

	// output the manifest-tool version for troubleshooting
	err = execCmd(versionCmd())
	if err != nil {
		return err
	}

	manifestSpecs, err := NewManifestSpec(p.Registry, p.Repo)
	if err != nil {
		return err
	}

	a := &afero.Afero{
		Fs: appFS,
	}

	err = a.Mkdir("/root/specs", 0755)
	if err != nil {
		return err
	}

	for i, spec := range manifestSpecs {
		var data bytes.Buffer

		fmt.Printf("Processing manifest list/image index %s\n", spec.Image)

		err = spec.Render(&data)
		if err != nil {
			return err
		}

		fmt.Printf("Rendered spec file:\n%s\n", data.String())
		specFilename := fmt.Sprintf("/root/specs/spec_%d.yml", i)

		err = a.WriteFile(specFilename, data.Bytes(), 0644)
		if err != nil {
			return err
		}

		cmd := p.Command(specFilename)
		// If a dry run, return without executing the cmd
		if p.Registry.DryRun {
			fmt.Println("Not pushing manifest list/image index as dry_run is true")
		} else {
			// run manifest-tool command from plugin configuration
			err = execCmd(cmd)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Validate verifies the Plugin is properly configured.
func (p *Plugin) Validate() error {
	logrus.Debug("validating plugin configuration")

	var err error

	// validate registry configuration
	err = p.Registry.Validate()
	if err != nil {
		return err
	}

	// validate repo configuration
	err = p.Repo.Validate()
	if err != nil {
		return err
	}

	manifestSpecs, err := NewManifestSpec(p.Registry, p.Repo)
	for _, ms := range manifestSpecs {
		err = ms.Validate()
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	p.manifestSpecs = manifestSpecs

	return nil
}
