# gen

A tool to generate license and .gitignore files locally.
This means that the binay is 'huuuuuuge', if this is unacceptable, look at this nice projects:

- <https://choosealicense.com/> (managed by GitHub)
- <https://github.com/tcnksm/license> (Uses GitHub Api, written in Go)
- <https://github.com/licenses/lice> (Uses template files which can enable 'add your license as you go' model, written in Python)

The files used here:

- <https://github.com/github/choosealicense.com>
- <https://github.com/github/gitignore>

## Installation

Using go install

```
go install codeberg.org/gekkowrld/gen@latest
```

Build from source:

- Clone the repository
- Run `go build -ldflags "-s -w"` to build the executable.
- Copy the binary to $PATH.

Run `gen completion $(basename ${0})` (your shell basically) to get the autocompletion text.

## Available licenses

Use `gen license --all` to view all available licenses and their descriptions.
This also displays their name which can be used to generate the license inside the brackets.

`gpl3` - GNU General Public License v3.0 
`mit` -  MIT License

Run `gen license --help` to view helpful information like flags and their defaults.

## Available .gitignore

Use `gen gitignore --all` to view all available gitignore templates.

## Help text

```txt
Generate .gitignore and license files locally

A near bad clone of:
https://github.com/generate/generate-license
https://github.com/generate/generate-gitignore

Usage:
  gen [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  license     Generate a license file

Flags:
  -h, --help      help for gen
  -v, --version   version for gen

Use "gen [command] --help" for more information about a command.
```

## Author

[Gekko Wrld](https://codeberg.org/gekkowrld/)
