# gobin-info

`gobin-info` lists your locally installed Go binaries alongside info about their originating Git repository.

## Raison d'être

Most of your CLI tools might be installed via a package manager, like `apt` or `dnf` on Linux, [Homebrew](https://brew.sh/) on macOS, or [Scoop](https://scoop.sh/) on Windows. Then if you want to get the list of your installed tools, you can run `apt list --installed`, `brew list` or `scoop list` to list them, and if you want to know more about one of them you can run `apt show ...`, `brew info ...` or `scoop info ...`.

But what about the ones you installed via Go? You installed them via `go install ...` and they live in `$GOPATH/bin` or `$GOBIN` or maybe you move/symlink them to `/usr/local/bin` or so. But you don't immediately know the origin of the tools. For example if there's a binary called `arc`, it could be `github.com/mholt/archiver/v3/cmd/arc` or `github.com/evilsocket/arc/cmd/arc` for example. You could run `arc --help` and it might give a hint what exactly it is. Or you run `go version -m /path/to/arc` and check the `path` or `mod` line for the GitHub repo.

But that's cumbersome if you want to list the origins of all your installed Go binaries, for example when you're migrating from an old to a new laptop and want to install the same tools you had on your old one. You'd want to automate this a bit, perhaps by writing a script or Go CLI to iterate through all of them, in the various directories in which they could reside. Some projects also use vanity URLs like `gopkg.in` or their own instance of it, so you still don't know the original Git repo address and have to check where it redirects to.

=> `gobin-info` does all of that

## Installation

`go install github.com/philippgille/gobin-info/cmd/gobin-info@latest`

## Usage

You can run `gobin-info` in several modes:

1. `gobin-info /path/to/dir` lists info about the Go binaries in a given directory (relative or absolute)
2. `gobin-info -wd` lists info about the Go binaries in your *working directory*
3. `gobin-info -gobin` lists info about the Go binaries in your *`$GOBIN`* directory
4. `gobin-info -gopath` lists info about the Go binaries in your *`$GOPATH/bin`* directory
5. `gobin-info -path` lists info about the Go binaries in your *`$PATH`*
   - ⚠ Currently this includes the ones you install via package managers like Homebrew, but in the future we might add an option to try and exclude those

### Example

```text
$ gobin-info -gopath
TODO
```
