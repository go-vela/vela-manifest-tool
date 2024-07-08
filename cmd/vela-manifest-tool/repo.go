// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type (
	// Repo represents the plugin configuration for repo information.
	Repo struct {
		Name              string   // name of the repository for the image
		Tags              []string // tags of the image for the repository
		Platforms         []string // platforms which should be included in the manifest
		ComponentTemplate string   // Template used to render each component image
	}
)

// Validate verifies the Repo is properly configured.
func (r *Repo) Validate() error {
	logrus.Trace("validating repo plugin configuration")

	// verify repo is provided
	if len(r.Name) == 0 {
		return fmt.Errorf("no repo name provided")
	}

	// check if tags are provided
	if len(r.Tags) > 0 {
		// check each tag value for valid docker tag syntax
		for _, tag := range r.Tags {
			if !tagRegexp.MatchString(tag) {
				return fmt.Errorf(errTagValidation, tag)
			}
		}
	} else {
		return fmt.Errorf("no tags provided")
	}

	if len(r.Platforms) > 0 {
		for _, platform := range r.Platforms {
			if _, ok := allowedPlatforms[platform]; !ok {
				return fmt.Errorf("unsupported platform %s requested", platform)
			}
		}
	} else {
		return fmt.Errorf("no platforms provided")
	}

	return nil
}
