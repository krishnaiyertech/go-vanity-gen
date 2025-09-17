# SPDX-FileCopyrightText: Copyright 2025 Krishna Iyer <www.krishnaiyer.tech>
# SPDX-License-Identifier: Apache-2.0

GO_LINT=golangci-lint
GVG_VERSION=v0.0.1
GVG_GIT_COMMIT=$(shell git rev-parse --short HEAD)
GVG_DATE=$(shell date)
GVG_PACKAGE="krishnaiyer.tech/golang/go-vanity-gen"

.PHONY: init

init:
	@echo "Initialise repository..."
	@mkdir -p gen

test:
	go test ./... -cover

build.local:
	go build \
	-ldflags="-X '${GVG_PACKAGE}/cmd.version=${GVG_VERSION}' \
	-X '${GVG_PACKAGE}/cmd.gitCommit=${GVG_GIT_COMMIT}' \
	-X '${GVG_PACKAGE}/cmd.buildDate=${GVG_DATE}'" main.go

build.dist:
	goreleaser --snapshot --skip-publish --rm-dist

clean:
	@rm -rf dist
	@rm -rf gen
	@mkdir -p gen

lint:
	${GO_LINT} run