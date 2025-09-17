// SPDX-FileCopyrightText: Copyright 2025 Krishna Iyer <www.krishnaiyer.tech>
// SPDX-License-Identifier: Apache-2.0

// Package generator provides functions to generate the static files.
package generator

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"

	"gopkg.in/yaml.v2"
)

// Path is the parsed configuration of vanity paths.
type Path struct {
	path     string
	repo     string
	display  string
	vcs      string
	packages []string
}

// Generator generates vanity assets.
type Generator struct {
	cfg   config
	paths []Path
	host  string
}

// OutItem is a single output item.
type OutItem struct {
	PkgNames []string
	Content  []byte
}

// Out is the raw output from the generator.
type Out struct {
	items map[string]OutItem
}

// Items returns the generated output items.
func (o *Out) Items() map[string]OutItem {
	return o.items
}

// config is the vanity config.
type config struct {
	Host  string `yaml:"host,omitempty"`
	Paths map[string]struct {
		Repo     string   `yaml:"repo,omitempty"`
		Display  string   `yaml:"display,omitempty"`
		VCS      string   `yaml:"vcs,omitempty"`
		Packages []string `yaml:"packages,omitempty"`
	} `yaml:"paths,omitempty"`
}

// New parses the provided vanity config and returns a new Generator.
func New(_ context.Context, vanity []byte) (*Generator, error) {
	var vanityCfg config
	if err := yaml.Unmarshal(vanity, &vanityCfg); err != nil {
		return nil, fmt.Errorf("could not parse vanity config: %w", err)
	}
	paths := make([]Path, 0)
	for path, config := range vanityCfg.Paths {
		project := Path{
			path:     strings.TrimSuffix(path, "/"),
			repo:     config.Repo,
			display:  config.Display,
			vcs:      config.VCS,
			packages: config.Packages,
		}
		switch {
		case config.Display != "":
		case strings.HasPrefix(config.Repo, "https://github.com/"):
			project.display = fmt.Sprintf(
				"%v %v/tree/master{/dir} %v/blob/master{/dir}/{file}#L{line}",
				config.Repo,
				config.Repo,
				config.Repo,
			)
		case strings.HasPrefix(config.Repo, "https://bitbucket.org"):
			project.display = fmt.Sprintf(
				"%v %v/src/default{/dir} %v/src/default{/dir}/{file}#{file}-{line}",
				config.Repo,
				config.Repo,
				config.Repo,
			)
		}
		switch {
		case config.VCS != "":
			if config.VCS != "bzr" && config.VCS != "git" && config.VCS != "hg" && config.VCS != "svn" {
				return nil, fmt.Errorf("configuration for %v: unknown VCS %s", path, config.VCS)
			}
		case strings.HasPrefix(config.Repo, "https://github.com/"):
			project.vcs = "git"
		default:
			return nil, fmt.Errorf("configuration for %v: cannot infer VCS from %s", path, config.Repo)
		}
		paths = append(paths, project)
	}
	return &Generator{
		cfg:   vanityCfg,
		paths: paths,
		host:  vanityCfg.Host,
	}, nil
}

type vanity struct {
	Path string
	Repo string
}

// Index generates the index.html at the root of the assets tree.
func (gen *Generator) Index(_ context.Context, input string) ([]byte, error) {
	index, err := template.New("index").Parse(input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}
	vanityPaths := make([]vanity, len(gen.paths))
	for i, h := range gen.paths {
		vanityPaths[i].Path = gen.host + h.path
		vanityPaths[i].Repo = h.repo
	}
	var buf bytes.Buffer
	if err := index.Execute(&buf, struct {
		Host   string
		Vanity []vanity
	}{
		Host:   gen.host,
		Vanity: vanityPaths,
	},
	); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.Bytes(), nil
}

// Project generates the index.html for a project path and returns an Out object.
func (gen *Generator) Project(_ context.Context, input string) (*Out, error) {
	out := &Out{
		items: make(map[string]OutItem, 0),
	}
	vanity, err := template.New("vanity").Parse(input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}
	for _, path := range gen.paths {
		var buf bytes.Buffer
		if err := vanity.Execute(&buf, struct {
			Import  string
			Subpath string
			Repo    string
			Display string
			VCS     string
			Host    string
		}{
			Import:  gen.host + path.path,
			Repo:    path.repo,
			Display: path.display,
			VCS:     path.vcs,
			Host:    gen.host,
		}); err != nil {
			return nil, fmt.Errorf("failed to execute template: %w", err)
		}
		out.items[path.path] = OutItem{
			PkgNames: path.packages,
			Content:  buf.Bytes(),
		}
	}
	return out, nil
}

// Paths returns the list of paths extracted from the vanity configuration.
func (gen *Generator) Paths() []Path {
	return gen.paths
}
