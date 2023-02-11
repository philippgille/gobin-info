# gobin-info

`gobin-info` lists your locally installed Go binaries alongside their version and original Git repository.

It's kind of like a convenience wrapper around `go version -m ...` with some niceties on top, like vanity URL resolving.

## Installation

`go install github.com/philippgille/gobin-info@latest`

## Usage

You can run `gobin-info` in several modes:

- `gobin-info /path/to/dir` lists info about the Go binaries in a given directory (relative or absolute)
- `gobin-info -wd` lists info about the Go binaries in your *working directory*
- `gobin-info -gobin` lists info about the Go binaries in your *`$GOBIN`* directory
- `gobin-info -gopath` lists info about the Go binaries in your *`$GOPATH/bin`* directory
- 🚧 `gobin-info -path` lists info about the Go binaries in your *`$PATH`* (not implemented yet)

It prints a `(❓)` after the URL in case the URL couldn't be reliably determined.

> Note: `gobin-info` doesn't recurse into subdirectories. This might be added with an optional flag in the future.

### Example

```text
$ gobin-info -gopath
Scanning /home/johndoe/go/bin
arc         v3.5.1  https://github.com/mholt/archiver
gopls       v0.11.0 https://go.googlesource.com/tools
mage        (devel) https://github.com/magefile/mage
staticcheck v0.3.3  https://github.com/dominikh/go-tools
```

## Raison d'être

Most of your CLI tools were probably installed with a package manager like `apt` or `dnf` on Linux, [Homebrew](https://brew.sh/) on macOS, or [Scoop](https://scoop.sh/) on Windows. Then if you want to get the list of your installed tools, you can run `apt list --installed`, `brew list` or `scoop list` to list them, and if you want to know more about one of them you can run `apt show ...`, `brew info ...` or `scoop info ...`.

But what about the ones you installed with Go? You installed them with `go install ...` and they live in `$GOPATH/bin` or `$GOBIN` or maybe you move/symlink them to `/usr/local/bin` or so.

- Now you don't immediately know the origin of the tools. For example if there's a binary called `arc`, is it `github.com/mholt/archiver/v3/cmd/arc` or `github.com/evilsocket/arc/cmd/arc`?
- You could run `arc --help` and it might give a hint what exactly it is, but it's not reliable
- Or you run `go version -m /path/to/arc` and among the dozens of output lines you check the `path` or `mod`
  - But their values are not `https://`-prefixed, so you can't click them in your terminal and have to copy paste them into your browser
  - Then for example `arc` has the module path `github.com/mholt/archiver/v3`, which leads to a `404 Not Found` error on GitHub because of the `v3`
  - And for `staticcheck` the module path is `honnef.co/go/tools`, which is a vanity URL that doesn't point to the original Git repository (<https://github.com/dominikh/go-tools>) and the browser also doesn't redirect to it

`gobin-info` makes all of this much easier.
