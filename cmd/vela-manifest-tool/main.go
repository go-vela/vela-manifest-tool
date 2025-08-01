// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-vela/vela-manifest-tool/version"
)

func main() {
	v := version.New()

	// serialize the version information as pretty JSON
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}

	// output the version information to stdout
	fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

	// create new CLI application
	app := &cli.Command{
		Name:      "vela-manifest-tool",
		Usage:     "Vela Manifest Tool plugin for building and publishing manifest lists/image indices",
		Copyright: "Copyright 2024 Target Brands, Inc. All rights reserved.",
		Authors: []any{
			&mail.Address{
				Name:    "Vela Admins",
				Address: "vela@target.com",
			},
		},
		Action:  run,
		Version: v.Semantic(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "log.level",
				Value: "info",
				Usage: "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_LOG_LEVEL"),
					cli.EnvVar("MANIFEST_TOOL_LOG_LEVEL"),
					cli.File("/vela/parameters/manifest_tool/log_level"),
					cli.File("/vela/secrets/manifest_tool/log_level"),
				),
			},

			// Registry Flags
			&cli.BoolFlag{
				Name:  "registry.dry_run",
				Usage: "enables building images without publishing to the registry",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_DRY_RUN"),
					cli.EnvVar("MANIFEST_TOOL_DRY_RUN"),
					cli.File("/vela/parameters/manifest_tool/dry_run"),
					cli.File("/vela/secrets/manifest_tool/dry_run"),
				),
			},
			&cli.StringFlag{
				Name:  "registry.name",
				Value: "index.docker.io",
				Usage: "Docker registry name to communicate with",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_REGISTRY"),
					cli.EnvVar("MANIFEST_TOOL_REGISTRY"),
					cli.File("/vela/parameters/manifest_tool/registry"),
					cli.File("/vela/secrets/manifest_tool/registry"),
				),
			},
			&cli.StringFlag{
				Name:  "registry.username",
				Usage: "user name for communication with the registry",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_USERNAME"),
					cli.EnvVar("MANIFEST_TOOL_USERNAME"),
					cli.EnvVar("DOCKER_USERNAME"),
					cli.File("/vela/parameters/manifest_tool/username"),
					cli.File("/vela/secrets/manifest_tool/username"),
					cli.File("/vela/secrets/managed-auth/username"),
				),
			},
			&cli.StringFlag{
				Name:  "registry.password",
				Usage: "password for communication with the registry",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_PASSWORD"),
					cli.EnvVar("MANIFEST_TOOL_PASSWORD"),
					cli.EnvVar("DOCKER_PASSWORD"),
					cli.File("/vela/parameters/manifest_tool/password"),
					cli.File("/vela/secrets/manifest_tool/password"),
					cli.File("/vela/secrets/managed-auth/password"),
				),
			},
			&cli.IntFlag{
				Name:  "registry.push_retry",
				Usage: "number of retries for pushing an image to a remote destination",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_PUSH_RETRY"),
					cli.EnvVar("MANIFEST_TOOL_PUSH_RETRY"),
					cli.File("/vela/parameters/manifest_tool/push_retry"),
					cli.File("/vela/secrets/manifest_tool/push_retry"),
				),
			},

			// Repo Flags
			&cli.StringFlag{
				Name:  "repo.name",
				Usage: "repository name for the image",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_REPO"),
					cli.EnvVar("MANIFEST_TOOL_REPO"),
					cli.File("/vela/parameters/manifest_tool/repo"),
					cli.File("/vela/secrets/manifest_tool/repo"),
				),
			},
			&cli.StringSliceFlag{
				Name:  "repo.tags",
				Value: []string{"latest"},
				Usage: "repository tags of the manifest list/image index",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_TAGS"),
					cli.EnvVar("MANIFEST_TOOL_TAGS"),
					cli.File("/vela/parameters/manifest_tool/tags"),
					cli.File("/vela/secrets/manifest_tool/tags"),
				),
			},
			&cli.StringSliceFlag{
				Name:  "repo.platforms",
				Value: []string{"linux/amd64", "linux/arm64/v8"},
				Usage: "docker platforms to include in the manifest list/image index",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_PLATFORMS"),
					cli.EnvVar("MANIFEST_TOOL_PLATFORMS"),
					cli.File("/vela/parameters/manifest_tool/platforms"),
					cli.File("/vela/secrets/manifest_tool/platforms"),
				),
			},
			&cli.StringFlag{
				Name:  "repo.component_template",
				Value: "{{.Repo}}:{{.Tag}}-{{.Os}}-{{.Arch}}{{if .Variant}}-{{.Variant}}{{end}}",
				Usage: "template used to render each component image",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_COMPONENT_TEMPLATE"),
					cli.EnvVar("MANIFEST_TOOL_COMPONENT_TEMPLATE"),
					cli.File("/vela/parameters/manifest_tool/component_template"),
					cli.File("/vela/secrets/manifest_tool/component_template"),
				),
			},
		},
	}

	err = app.Run(context.Background(), os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(_ context.Context, cmd *cli.Command) error {
	// set the log level for the plugin
	switch cmd.String("log.level") {
	case "t", "trace", "Trace", "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "w", "warn", "Warn", "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "e", "error", "Error", "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "f", "fatal", "Fatal", "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "p", "panic", "Panic", "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	case "i", "info", "Info", "INFO":
		fallthrough
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.WithFields(logrus.Fields{
		"code":     "https://github.com/go-vela/vela-manifest-tool",
		"docs":     "https://go-vela.github.io/docs/plugins/registry/pipeline/manifest-tool",
		"registry": "https://hub.docker.com/r/target/vela-manifest-tool",
	}).Info("Vela Manifest Tool Plugin")

	// create the plugin
	p := &Plugin{
		// build configuration
		// registry configuration
		Registry: &Registry{
			DryRun:    cmd.Bool("registry.dry_run"),
			Name:      cmd.String("registry.name"),
			Username:  cmd.String("registry.username"),
			Password:  cmd.String("registry.password"),
			PushRetry: cmd.Int("registry.push_retry"),
		},
		// repo configuration
		Repo: &Repo{
			Name:              cmd.String("repo.name"),
			Tags:              cmd.StringSlice("repo.tags"),
			Platforms:         cmd.StringSlice("repo.platforms"),
			ComponentTemplate: cmd.String("repo.component_template"),
		},
	}

	// validate the plugin
	err := p.Validate()
	if err != nil {
		return err
	}

	// execute the plugin
	return p.Exec()
}
