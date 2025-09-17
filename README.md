# go-vanity-gen

Generates simple html pages for Golang's vanity redirection.

TODO: Add link to blog

## Installation

```bash
$ go get -u krishnaiyer.tech/golang/go-vanity-gen@<version>
```

## Options

```bash
go-vanity-gen generates vanity assets from templates. Templates are usually simple html files that contain links to repositories

Usage:
  go-vanity-gen [flags]
  go-vanity-gen [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Display version information

Flags:
  -h, --help         help for go-vanity-gen
  -i, --in string    directory where input files. Must contain index.tmpl, project.tmpl and vanity.yml
  -o, --out string   directory where output files are generated. Default is ./gen

Use "go-vanity-gen [command] --help" for more information about a command.
```

## Usage

1. Install

```bash
$ go install krishnaiyer.nl/golang/go-vanity-gen
```

2. Generate

```bash
$ go-vanity-gen -i sample -o gen
```

## Development

The following are the prerequisites;

- [Go](https://golang.org/) installed and configured.
- [optional] Custom HTML templates.

1. Clone this repository.
2. Initialize

```bash
$ make init
```

3. Build from source

```bash
$ GVG_PACKAGE=<your-path>/go-vanity-gen GVG_VERSION=<version> make build.local
```

4. Run tests

```bash
$ make test
```

5. Clean up

```bash
$ make clean
```

## Limitations

- Since we're building static assets, each package that need redirection needs an `index.html`. This does result in a lot of duplication. But since each file is very small, the cost of storing these files is much lower than running a server to serve these paths.
- Package paths need to be manually listed. In the root of each of your repositories, use the following command.

```bash
$ go list ./...
```

## Releases

Please check the [changelog](./CHANGELOG.md) for release history.

## License

This project is licensed as-is under the [Apache 2.0 license](./LICENSES/Apache-2.0.txt).
