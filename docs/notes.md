# Notes

## TODO:

- Add support for more Git providers
  - SourceForge: Only found mirrors from elsewhere. How does a natively hosted project look like?
  - codegiant.io: Didn't find public repos
  - launchpad.net: No typical owner/project separation on the main project site? E.g. <https://launchpad.net/apparmor>. Only the tree has it: <https://code.launchpad.net/~apparmor-dev/apparmor/+git/apparmor/+ref/master>
- Go 1.20 changed `go version -m` a bit (e.g. working with non-executable files). Take inspiration from that.
- Improve repo examples in notes.md to contain package and module
- Parallelize iterating through files (mostly useful for parallelizing HTTP requests for vanity URL resolving)
- Make vanity resolving optional (some people might want to prevent HTTP requests)
- Add tests (see below "To test" section)
- `-v` flag for a more verbose mode (version-specific tree URL, package path for easier reinstall)

## Repo examples

- github.com
  - Package: github.com/mholt/archiver/v3/cmd/arc
  - Module: github.com/mholt/archiver/v3
  - URL: <https://github.com/mholt/archiver>
  - With tree: <https://github.com/mholt/archiver/tree/v3.5.1>
- gitlab.com
  - Package: ?
  - Module: ?
  - URL: <https://gitlab.com/gitlab-org/gitlab-runner>
  - With tree: <https://gitlab.com/gitlab-org/gitlab-runner/-/tree/v15.5.2>
- bitbucket.org
  - Package: ?
  - Module: ?
  - URL: <https://bitbucket.org/pcas/golua/> (but always redirects to tree specific URL on main branch)
  - With tree: <https://bitbucket.org/pcas/golua/src/v0.1.6/>
- sr.ht
  - Package: ?
  - Module: ?
  - URL: <https://sr.ht/~hedy/gelim/>
  - With tree: <https://git.sr.ht/~hedy/gelim/tree/v0.9.3>
- cs.opensource.google
  - Package: ?
  - Module: ?
  - URL: <https://cs.opensource.google/gvisor/gvisor>
  - With tree: <https://cs.opensource.google/gvisor/gvisor/+/refs/tags/release-20230123.0:>
- gitee.com
  - Package: ?
  - Module: ?
  - URL: <https://gitee.com/mirrors/gohugo>
  - With tree: <https://gitee.com/mirrors/gohugo/tree/v0.74.3/>
- codeberg.org
  - Package: codeberg.org/anaseto/goal
  - Module: codeberg.org/anaseto/goal
  - URL: <https://codeberg.org/anaseto/goal>
  - With tree: <https://codeberg.org/anaseto/goal/src/tag/v0.1.0>

## To test

- golang.org/dl/go1.17
- golang.org/x/tools/gopls
- honnef.co/go/tools/cmd/staticcheck
- gvisor (opensource.google)
  - maybe mention as example vs github.com/mjwhitta/runsc in README
- module with v2 suffix
- module that's doesn't have go.mod in repo root
- various Git providers
