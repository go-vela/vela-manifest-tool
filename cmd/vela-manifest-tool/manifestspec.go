// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var allowedPlatforms = map[string]bool{
	"linux/amd64":    true,
	"linux/arm64":    true,
	"linux/arm64/v8": true,
	"linux/arm":      true,
	"linux/arm/v7":   true,
}

type Manifest struct {
	Spec     ManifestSpec
	Context  ComponentContext
	Template *template.Template
}

// ManifestSpec represents the structure of the manifest-tool yaml spec file.
type ManifestSpec struct {
	Image     string              // name of the image index including tag
	Manifests []ManifestComponent // list of component images to include in index
}

type ManifestPlatform struct {
	Os           string
	Architecture string
	Variant      string `yaml:",omitempty"`
}

type ManifestComponent struct {
	Image    string           // name of the component image to be referenced by the index
	Platform ManifestPlatform // The platform specification for the component image
}

type ComponentContext struct {
	Repo    string
	Tag     string
	Os      string
	Arch    string
	Variant string
}

func NewManifestSpec(reg *Registry, repo *Repo) ([]*ManifestSpec, error) {
	specs := []*ManifestSpec{}
	tmpl, err := template.New("component_template").Parse(repo.ComponentTemplate)

	if err != nil {
		return specs, err
	}

	if len(reg.Name) == 0 {
		return specs, fmt.Errorf("no registry name provided")
	}

	if len(repo.Name) == 0 {
		return specs, fmt.Errorf("no repository name provided")
	}

	for _, tag := range repo.Tags {
		ms := ManifestSpec{
			Image:     reg.Name + repo.Name + ":" + tag,
			Manifests: []ManifestComponent{},
		}

		for _, platform := range repo.Platforms {
			platformComp := strings.Split(platform, "/")
			if len(platformComp) < 2 {
				return nil, fmt.Errorf("malformed platform %s", platform)
			} else if len(platformComp) == 2 {
				// probably a better way to do this, just not sure how
				// else to make the variant below clean
				platformComp = append(platformComp, "")
			}

			ctx := ComponentContext{
				Repo:    repo.Name,
				Tag:     tag,
				Os:      platformComp[0],
				Arch:    platformComp[1],
				Variant: platformComp[2],
			}

			var compImgBuf bytes.Buffer
			err = tmpl.Execute(&compImgBuf, ctx)

			if err != nil {
				return specs, err
			}

			compImg := compImgBuf.String()
			comp := ManifestComponent{
				Image: fmt.Sprintf("%s%s", reg.Name, compImg),
				Platform: ManifestPlatform{
					Os:           ctx.Os,
					Architecture: ctx.Arch,
					Variant:      ctx.Variant,
				},
			}
			ms.Manifests = append(ms.Manifests, comp)
		}

		specs = append(specs, &ms)
	}

	return specs, nil
}

func (ms *ManifestSpec) Validate() error {
	logrus.Trace("validating manifest spec plugin configuration")

	// verify repo is provided
	if len(ms.Image) == 0 {
		return fmt.Errorf("no top-level image provided")
	}

	err := validateTagOfImage(ms.Image)
	if err != nil {
		return err
	}

	// check if tags are provided
	if len(ms.Manifests) > 0 {
		// check each tag value for valid docker tag syntax
		for _, compManifest := range ms.Manifests {
			err = validateTagOfImage(compManifest.Image)
			if err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("no component images provided")
	}

	return nil
}

func (ms *ManifestSpec) Render(wr io.Writer) error {
	yamlData, err := yaml.Marshal(ms)
	if err != nil {
		return err
	}

	_, err = wr.Write(yamlData)

	return err
}

func validateTagOfImage(fullImage string) error {
	topLevelImgParts := strings.Split(fullImage, ":")

	if len(topLevelImgParts) != 2 {
		return fmt.Errorf("%s not in image:tag format", fullImage)
	}

	if !tagRegexp.MatchString(topLevelImgParts[1]) {
		return fmt.Errorf(errTagValidation, topLevelImgParts[1])
	}

	return nil
}
