# gobin-info

`gobin-info` lists your locally installed Go binaries alongside info about their originating Git repository.

## Raison d'être

Most of your CLI tools might be installed via a package manager, like `apt` or `dnf` on Linux, [Homebrew](https://brew.sh/) on macOS, or [Scoop](https://scoop.sh/) on Windows. Then if you want to get the list of your installed tools, you can run `apt list --installed`, `brew list` or `scoop list` to list them, and if you want to know more about one of them you can run `apt show ...`, `brew info ...` or `scoop info ...`.

But what about the ones you installed via Go? You installed them via `go install ...` and they live in `$GOPATH/bin` or `$GOBIN` or maybe you move/symlink them to `/usr/local/bin` or so. But you don't immediately know the origin of the tools. For example if there's a binary called `arc`, it could be `github.com/mholt/archiver/v3/cmd/arc` or `github.com/evilsocket/arc/cmd/arc` for example. You could run `arc --help` and it might give a hint what exactly it is. Or you run `go version -m /path/to/arc` and check the `path` or `mod` line for the GitHub repo.

But that's cumbersome if you want to list the origins of all your installed Go binaries, for example when you're migrating from an old to a new laptop and want to install the same tools you had on your old one. You'd want to automate this a bit, perhaps by writing a script or Go CLI to iterate through all of them, in the various directories in which they could reside. Some projects also use vanity URLs like `gopkg.in` or their own instance of it, so you still don't know the original Git repo address and have to check where it redirects to.

=> `gobin-info` does all of that

## Installation

`go install github.com/philippgille/gobin-info@latest`

## Usage

You can run `gobin-info` in several modes:

- `gobin-info /path/to/dir` lists info about the Go binaries in a given directory (relative or absolute)
- `gobin-info -wd` lists info about the Go binaries in your *working directory*
- `gobin-info -gobin` lists info about the Go binaries in your *`$GOBIN`* directory
- `gobin-info -gopath` lists info about the Go binaries in your *`$GOPATH/bin`* directory
- 🚧 `gobin-info -path` lists info about the Go binaries in your *`$PATH`*

### Example

```text
$ gobin-info -gopath
Scanning /home/johndoe/go/bin
arc          v3.5.1                             https://github.com/mholt/archiver
dlv          v1.20.1                            https://github.com/go-delve/delve
fyne_demo    v2.3.0                             https://github.com/fyne-io/fyne
go-outline   v0.0.0-20210608161538-9736a4bde949 https://github.com/ramya-rao-a/go-outline
go1.17       v0.0.0-20220609182932-6cd2f0e318f7 ❓https://go.googlesource.com/dl❓
gomodifytags v1.16.0                            https://github.com/fatih/gomodifytags
gopkgs       v2.1.2                             https://github.com/uudashr/gopkgs
goplay       v1.0.0                             https://github.com/haya14busa/goplay
gopls        v0.11.0                            ❓https://go.googlesource.com/tools❓
gotests      v1.6.0                             https://github.com/cweill/gotests
impl         v1.1.0                             https://github.com/josharian/impl
lazygit      v0.32.2                            https://github.com/jesseduffield/lazygit
mage         (devel)                            https://github.com/magefile/mage
staticcheck  v0.3.3                             https://github.com/dominikh/go-tools
```
