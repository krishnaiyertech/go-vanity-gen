// SPDX-FileCopyrightText: Copyright 2025 Krishna Iyer <www.krishnaiyer.tech>
// SPDX-License-Identifier: Apache-2.0

// Package cmd provides command line options.
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"krishnaiyer.tech/golang/go-vanity-gen/pkg/generator"
)

// Config represents the configuration
type Config struct {
	In  string `name:"in" short:"i" description:"directory where input files. Must contain index.tmpl, project.tmpl and vanity.yml"`
	Out string `name:"out" short:"o" description:"directory where output files are generated. Default is ./gen"`
}

type fileData map[string][]byte

var (
	flags = pflag.NewFlagSet("go-vanity", pflag.ExitOnError)

	config = new(Config)

	addressRegex = regexp.MustCompile(`^([a-z-.0-9]+)(:[0-9]+)?$`)

	errTemplateNotDefined = fmt.Errorf("Template not defined")

	// Root is the root command.
	Root = &cobra.Command{
		Use:           "go-vanity",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "go-vanity generates vanity assets from templates",
		Long:          `go-vanity generates vanity assets from templates. Templates are usually simple html files that contain links to repositories.`,
		Run: func(cmd *cobra.Command, args []string) {
			baseCtx := context.Background()
			ctx, cancel := context.WithCancel(baseCtx)
			defer cancel()

			if config.Out == "" {
				config.Out = "./gen"
			}

			var (
				indexFile   = filepath.Join(config.In, "index.tmpl")
				projectFile = filepath.Join(config.In, "project.tmpl")
				vanityFile  = filepath.Join(config.In, "vanity.yml")
			)

			input := fileData{
				indexFile:   nil,
				projectFile: nil,
				vanityFile:  nil,
			}
			for name := range input {
				raw, err := os.ReadFile(name)
				if err != nil {
					log.Fatal(fmt.Errorf("Failed to read file %s: %v", name, err.Error()))
				}
				input[name] = raw
			}

			gen, err := generator.New(ctx, input[vanityFile])
			if err != nil {
				log.Fatal(err.Error())
			}

			index, err := gen.Index(ctx, string(input[indexFile]))
			if err != nil {
				log.Fatal(fmt.Errorf("Failed to generate index :%v", err.Error()))
			}
			err = os.WriteFile(filepath.Join(config.Out, "index.html"), index, 0755)
			if err != nil {
				log.Fatal(fmt.Errorf("Failed to write index at %s :%v", indexFile, err.Error()))
			}

			out, err := gen.Project(ctx, string(input[projectFile]))
			if err != nil {
				log.Fatal(fmt.Errorf("Failed to generate project files :%v", err.Error()))
			}

			for name, project := range out.Items() {
				basePath := fmt.Sprintf("%s%s", config.Out, name)
				paths := []string{basePath}
				paths = append(paths, project.PkgNames...)
				for _, path := range paths {
					if path != basePath {
						path = basePath + "/" + path
					}
					err := os.MkdirAll(path, 0755)
					if err != nil {
						log.Fatal(fmt.Errorf("Failed to create folder %s :%v", path, err.Error()))
					}
					err = os.WriteFile(path+"/index.html", project.Content, 0755)
					if err != nil {
						log.Fatal(fmt.Errorf("Failed to create file %s :%v", path+"/index.html", err.Error()))
					}
				}
			}
		},
	}
)

// Execute the root command
func Execute() {
	if err := Root.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	flags, err := genFlags(config)
	if err != nil {
		log.Fatal(err)
	}
	Root.PersistentFlags().AddFlagSet(flags)
	Root.AddCommand(VersionCommand(Root))
}
