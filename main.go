package main

import (
	"bufio"
	"debug/buildinfo"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// Example: <meta name="go-import" content="fyne.io/fyne/v2 git https://github.com/fyne-io/fyne">
var vanityRegex = regexp.MustCompile(`< *meta name="go-import" content=".+ \w+ (https?://.+)" *\\?>`)

var defaultGetOwnerRepoPair = func(modulePath string) (string, string, error) {
	subs := strings.Split(modulePath, "/")
	if len(subs) < 3 {
		return "", "", fmt.Errorf("couldn't determine owner and repo name in module path '%s'", modulePath)
	}
	return subs[1], subs[2], nil
}

// For known Git providers, we don't need to check vanity URL redirects.
var knownGitProviders = map[string]Funcs{
	"github.com": {
		GetOwnerRepoPair: defaultGetOwnerRepoPair,
		GetRepoURL: func(owner, repo string) string {
			return fmt.Sprintf("https://github.com/%s/%s", owner, repo)
		}},
	"gitlab.com": {
		GetOwnerRepoPair: defaultGetOwnerRepoPair,
		GetRepoURL: func(owner, repo string) string {
			return fmt.Sprintf("https://gitlab.com/%s/%s", owner, repo)
		}},
	"bitbucket.org": {
		GetOwnerRepoPair: defaultGetOwnerRepoPair,
		GetRepoURL: func(owner, repo string) string {
			return fmt.Sprintf("https://bitbucket.com/%s/%s", owner, repo)
		}},
	"sr.ht": {
		GetOwnerRepoPair: defaultGetOwnerRepoPair,
		GetRepoURL: func(owner, repo string) string {
			return fmt.Sprintf("https://sr.ht/%s/%s", owner, repo)
		}},
	"cs.opensource.google": {
		GetOwnerRepoPair: defaultGetOwnerRepoPair,
		GetRepoURL: func(owner, repo string) string {
			return fmt.Sprintf("https://cs.opensource.google/%s/%s", owner, repo)
		}},
	"gitee.com": {
		GetOwnerRepoPair: defaultGetOwnerRepoPair,
		GetRepoURL: func(owner, repo string) string {
			return fmt.Sprintf("https://gitee.com/%s/%s", owner, repo)
		}},
	"codeberg.org": {
		GetOwnerRepoPair: defaultGetOwnerRepoPair,
		GetRepoURL: func(owner, repo string) string {
			return fmt.Sprintf("https://codeberg.org/%s/%s", owner, repo)
		}},
}

type Funcs struct {
	// Takes module path, returns owner and repo name
	GetOwnerRepoPair func(string) (string, string, error)
	// Takes Git owner and repo name, returns repo URL
	GetRepoURL func(string, string) string
}

type BinInfo struct {
	filename string // Without path, e.g. `arc`/`arc.exe`

	// 4 values from `go version -m`
	packagePath   string // e.g. `github.com/mholt/archiver/v3/cmd/arc`
	modulePath    string // e.g. `github.com/mholt/archiver/v3`
	moduleVersion string // Just the Git tag, not the full version as reported by `go version -m`. E.g. `v3.5.1`
	vcsRevision   string // e.g. `cc194d2e4af2dc09a812aa0ff61adc4813ea6c69`

	repoURL string // URL that can be visited in a browser, after vanity URL resolving. E.g. for binary `arc` installed from `github.com/mholt/archiver/v3/cmd/arc@latest` with v3.5.1 being latest, it's "https://github.com/mholt/archiver". We could make this version-specific, to "https://github.com/mholt/archiver/tree/v3.5.1".
}

var (
	wd     = flag.Bool("wd", false, "Scan current working directory")
	gobin  = flag.Bool("gobin", false, `Scan "$GOBIN" directory`)
	gopath = flag.Bool("gopath", false, `Scan "$GOPATH/bin" directory`)
)

func main() {
	// Precondition: CLI must be called with one argument.
	// os.Args always holds the name of the program as the first argument.
	if len(os.Args) != 2 {
		log.Fatalln("gobin-info requires exactly one argument - either a path to a file/directory, or a flag.")
	}

	flag.Parse()

	// Get path
	path, err := getPath()
	if err != nil {
		log.Fatalln("Couldn't get path:", err)
	}

	// Iterate trough files in the path
	log.Println("Scanning", path)
	binInfos, err := scanDir(path)
	if err != nil {
		log.Fatalln("Error scanning dir:", err)
	}

	// Print result
	printResult(binInfos)
}

func getPath() (string, error) {
	var path string
	var err error
	if *wd {
		path, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("couldn't get current working directory: %w", err)
		}
	} else if *gobin {
		path = os.Getenv("GOBIN")
		if path == "" {
			return "", errors.New("GOBIN environment variable is empty or not set")
		}
	} else if *gopath {
		env := os.Getenv("GOPATH")
		if env == "" {
			// When the env var is not set, Go's own behavior is to use $HOME/go.
			// See https://pkg.go.dev/cmd/go#hdr-GOPATH_environment_variable
			log.Println("GOPATH is not set, falling back to $HOME/go like Go does")
			env, err = os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("couldn't get user home directory: %w", err)
			}
			env = filepath.Join(env, "go")
		}
		// GOPATH can actually be multiple directories, separated by colon on Unix, semicolon on Windows.
		// The $GOPATH/bin is always in the *first* of those directories.
		// See https://pkg.go.dev/cmd/go#hdr-GOPATH_environment_variable and https://go.dev/doc/code#Command
		env = strings.Split(env, string(os.PathListSeparator))[0]
		path = filepath.Join(env, "bin")
	} else {
		// Can be path to a file or directory
		path = os.Args[1]
	}

	return path, nil
}

// scanDir scans a directory for executables to run scanFile on.
func scanDir(dir string) ([]BinInfo, error) {
	var binInfos []BinInfo

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.Type().IsRegular() || d.Type()&fs.ModeSymlink != 0 {
			info, err := d.Info()
			if err != nil {
				return err
			}
			binInfo, err := scanFile(path, info)
			if err != nil {
				return err
			}
			// scanFile returns (nil, nil) if the file is not an executable
			if binInfo == nil {
				return nil
			}
			binInfos = append(binInfos, *binInfo)
		}
		return nil
	})

	return binInfos, err
}

// scanFile scans file to try to return the Go binary info.
// If the file is not a Go binary, scanFile returns (nil, nil).
func scanFile(file string, info fs.FileInfo) (*BinInfo, error) {
	if info.Mode()&fs.ModeSymlink != 0 {
		// Accept file symlinks only.
		i, err := os.Stat(file)
		if err != nil || !i.Mode().IsRegular() {
			return nil, err
		}
		info = i
	}

	if !isExe(file, info) {
		return nil, nil
	}

	bi, err := buildinfo.ReadFile(file)
	if err != nil {
		return nil, err
	}

	binInfo := BinInfo{
		filename: filepath.Base(file),

		packagePath:   bi.Path,
		modulePath:    bi.Main.Path,
		moduleVersion: bi.Main.Version, // Most binaries have the proper mod info, but github.com/magefile/mage for example doesn't. It's installed via their bootstrap tool, and the mod info is just "(devel)"
		vcsRevision:   "?",

		repoURL: "?",
	}

	// Add revision and derived URL

	// Look for the revision and set if found
	for _, biSetting := range bi.Settings {
		if biSetting.Key == "vcs.revision" {
			binInfo.vcsRevision = biSetting.Value
			break
		}
	}
	// Derive URL, potentially from vanity URL
	gitProvider := strings.Split(binInfo.modulePath, "/")[0]
	if funcs, ok := knownGitProviders[gitProvider]; !ok {
		// Provider not known; assume it's a vanity URL
		resolvedURL := resolveVanityURL(binInfo.packagePath, gitProvider)
		if resolvedURL == "" {
			// It wasn't a vanity URL. Probably unknown Git provider.
			binInfo.repoURL = fallbackURL(binInfo.modulePath)
		} else {
			// We assume that the resolved URL is prefixed with a protocol
			if strings.HasPrefix(resolvedURL, "http") {
				_, resolvedURL, ok = strings.Cut(resolvedURL, "//") // trim protocol
			}
			// Resolved URL might be known or unknown Git provider.
			gitProvider = strings.Split(resolvedURL, "/")[0]
			if funcs, ok := knownGitProviders[gitProvider]; !ok {
				// Unknown Git provider
				binInfo.repoURL = fallbackURL(resolvedURL)
			} else {
				owner, repo, err := funcs.GetOwnerRepoPair(resolvedURL)
				if err != nil {
					return nil, err
				}
				binInfo.repoURL = funcs.GetRepoURL(owner, repo)
			}
		}
	} else {
		// Known Git provider
		owner, repo, err := funcs.GetOwnerRepoPair(binInfo.modulePath)
		if err != nil {
			return nil, err
		}
		binInfo.repoURL = funcs.GetRepoURL(owner, repo)
	}

	return &binInfo, nil
}

// isExe reports whether the file should be considered executable.
func isExe(file string, info fs.FileInfo) bool {
	if runtime.GOOS == "windows" {
		return strings.HasSuffix(strings.ToLower(file), ".exe")
	}
	return info.Mode().IsRegular() && info.Mode()&0111 != 0
}

// TODO: Check if we can reuse the code Go uses for this:
// https://github.com/golang/go/blob/1102616/src/cmd/go/internal/get/discovery.go#L31
func resolveVanityURL(packagePath, gitProvider string) string {
	// We check the package path instead of module path because it's what's registered at the vanity URL service (the go install command calls the package URL)
	res, err := http.Get("https://" + packagePath)
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	scanner := bufio.NewScanner(res.Body)
	for ok := scanner.Scan(); ok; ok = scanner.Scan() {
		line := scanner.Text()
		match := vanityRegex.FindStringSubmatch(line)
		if match == nil {
			// Relevant meta tags are defined inside the head, so if we're already at the end of head we can stop.
			if strings.Contains(line, "</head>") {
				return ""
			}
			// Otherwise continue to next line
			continue
		}
		// We found a match. The redirect is defined as capture group in the regex.
		redir := match[1]
		_, err := url.Parse(redir)
		if err != nil {
			return ""
		}
		return redir
	}
	return ""
}

func fallbackURL(modulePath string) string {
	// In our known examples the URLs are always The host and then owner and repo as separate path elements.
	// Let's apply this, but warn that it might be wrong
	subs := strings.Split(modulePath, "/")
	if len(subs) < 3 {
		return "❓https://" + modulePath + "❓"
	}
	return "❓https://" + subs[0] + "/" + subs[1] + "/" + subs[2] + "❓"
}

func printResult(binInfos []BinInfo) {
	var maxFileNameLen int
	var maxVersionLen int
	for _, binInfo := range binInfos {
		if len(binInfo.filename) > maxFileNameLen {
			maxFileNameLen = len(binInfo.filename)
		}
		if len(binInfo.moduleVersion) > maxVersionLen {
			maxVersionLen = len(binInfo.moduleVersion)
		}
	}
	for _, binInfo := range binInfos {
		filenameWithPadding := binInfo.filename + strings.Repeat(" ", maxFileNameLen-len(binInfo.filename))
		versionWithPadding := binInfo.moduleVersion + strings.Repeat(" ", maxVersionLen-len(binInfo.moduleVersion))

		fmt.Printf("%s %s %s\n", filenameWithPadding, versionWithPadding, binInfo.repoURL)
	}
}
