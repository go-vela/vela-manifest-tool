// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

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
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-manifest-tool"
	app.HelpName = "vela-manifest-tool"
	app.Usage = "Vela Manifest Tool plugin for building and publishing manifest lists/image indices"
	app.Copyright = "Copyright 2024 Target Brands, Inc. All rights reserved."
	app.Authors = []*cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Plugin Metadata

	app.Action = run
	app.Compiled = time.Now()
	app.Version = v.Semantic()

	// Plugin Flags

	app.Flags = []cli.Flag{

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_LOG_LEVEL", "MANIFEST_TOOL_LOG_LEVEL"},
			FilePath: "/vela/parameters/manifest_tool/log_level,/vela/secrets/manifest_tool/log_level",
			Name:     "log.level",
			Usage:    "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:    "info",
		},

		// Registry Flags
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_DRY_RUN", "MANIFEST_TOOL_DRY_RUN"},
			FilePath: "/vela/parameters/manifest_tool/dry_run,/vela/secrets/manifest_tool/dry_run",
			Name:     "registry.dry_run",
			Usage:    "enables building images without publishing to the registry",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_REGISTRY", "MANIFEST_TOOL_REGISTRY"},
			FilePath: "/vela/parameters/manifest_tool/registry,/vela/secrets/manifest_tool/registry",
			Name:     "registry.name",
			Usage:    "Docker registry name to communicate with",
			Value:    "index.docker.io",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_USERNAME", "MANIFEST_TOOL_USERNAME", "DOCKER_USERNAME"},
			FilePath: "/vela/parameters/manifest_tool/username,/vela/secrets/manifest_tool/username,/vela/secrets/managed-auth/username",
			Name:     "registry.username",
			Usage:    "user name for communication with the registry",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_PASSWORD", "MANIFEST_TOOL_PASSWORD", "DOCKER_PASSWORD"},
			FilePath: "/vela/parameters/manifest_tool/password,/vela/secrets/manifest_tool/password,/vela/secrets/managed-auth/password",
			Name:     "registry.password",
			Usage:    "password for communication with the registry",
		},
		&cli.IntFlag{
			EnvVars:  []string{"PARAMETER_PUSH_RETRY", "MANIFEST_TOOL_PUSH_RETRY"},
			FilePath: "/vela/parameters/manifest_tool/push_retry,/vela/secrets/manifest_tool/push_retry",
			Name:     "registry.push_retry",
			Usage:    "number of retries for pushing an image to a remote destination",
		},

		// Repo Flags
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_REPO", "MANIFEST_TOOL_REPO"},
			FilePath: "/vela/parameters/manifest_tool/repo,/vela/secrets/manifest_tool/repo",
			Name:     "repo.name",
			Usage:    "repository name for the image",
		},
		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_TAGS", "MANIFEST_TOOL_TAGS"},
			FilePath: "/vela/parameters/manifest_tool/tags,/vela/secrets/manifest_tool/tags",
			Name:     "repo.tags",
			Usage:    "repository tags of the manifest list/image index",
			Value:    cli.NewStringSlice("latest"),
		},
		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_PLATFORMS", "MANIFEST_TOOL_PLATFORMS"},
			FilePath: "/vela/parameters/manifest_tool/tags,/vela/secrets/manifest_tool/platforms",
			Name:     "repo.platforms",
			Usage:    "docker platforms to include in the manifest list/image index",
			Value:    cli.NewStringSlice("linux/amd64", "linux/arm64/v8"),
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_COMPONENT_TEMPLATE", "MANIFEST_TOOL_COMPONENT_TEMPLATE"},
			FilePath: "/vela/parameters/manifest_tool/component_template,/vela/secrets/manifest_tool/component_template",
			Name:     "repo.component_template",
			Usage:    "template used to render each component image",
			Value:    "{{.Repo}}:{{.Tag}}-{{.Os}}-{{.Arch}}{{if .Variant}}-{{.Variant}}{{end}}",
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(c *cli.Context) error {
	// set the log level for the plugin
	switch c.String("log.level") {
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
			DryRun:    c.Bool("registry.dry_run"),
			Name:      c.String("registry.name"),
			Username:  c.String("registry.username"),
			Password:  c.String("registry.password"),
			PushRetry: c.Int("registry.push_retry"),
		},
		// repo configuration
		Repo: &Repo{
			Name:              c.String("repo.name"),
			Tags:              c.StringSlice("repo.tags"),
			Platforms:         c.StringSlice("repo.platforms"),
			ComponentTemplate: c.String("repo.component_template"),
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
