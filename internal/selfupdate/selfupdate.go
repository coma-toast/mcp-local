// Package selfupdate replaces the running mcp-local binary with a GitHub release build.
//
// Releases are published by .github/workflows/release.yml (GoReleaser) on every merge to
// main: one mcp-local_<version>_<os>_<arch>.tar.gz per platform plus checksums.txt.
// An archive is installed only when its SHA-256 matches checksums.txt.
package selfupdate

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	defaultAPIBase = "https://api.github.com"
	defaultRepo    = "coma-toast/mcp-local"
	binaryName     = "mcp-local"
	// maxDownload caps any single download; release archives are a few MB.
	maxDownload = 200 << 20
)

// Release is a published GitHub release.
type Release struct {
	Tag     string            // e.g. v0.1.9
	Version string            // Tag without the leading "v"
	Assets  map[string]string // asset name -> download URL
}

// Updater finds and installs releases. The zero value talks to github.com for this platform.
type Updater struct {
	APIBase string // GitHub API root; default https://api.github.com
	Repo    string // owner/name; default coma-toast/mcp-local
	Token   string // optional GitHub token, sent only to APIBase (raises the rate limit)
	GOOS    string // default runtime.GOOS
	GOARCH  string // default runtime.GOARCH
	Client  *http.Client
	// Verify runs the new binary before it replaces the old one; default runs "<path> --version".
	Verify func(path string) error
}

func (u *Updater) apiBase() string {
	if u.APIBase != "" {
		return strings.TrimRight(u.APIBase, "/")
	}
	return defaultAPIBase
}

func (u *Updater) repo() string {
	if u.Repo != "" {
		return u.Repo
	}
	return defaultRepo
}

func (u *Updater) client() *http.Client {
	if u.Client != nil {
		return u.Client
	}
	return &http.Client{Timeout: 2 * time.Minute}
}

func (u *Updater) platform() (string, string) {
	goos, goarch := u.GOOS, u.GOARCH
	if goos == "" {
		goos = runtime.GOOS
	}
	if goarch == "" {
		goarch = runtime.GOARCH
	}
	return goos, goarch
}

// Latest returns the newest published (non-draft, non-prerelease) release.
func (u *Updater) Latest(ctx context.Context) (*Release, error) {
	return u.release(ctx, "/releases/latest")
}

// ByTag returns the release for tag; a missing "v" prefix is added.
func (u *Updater) ByTag(ctx context.Context, tag string) (*Release, error) {
	if !strings.HasPrefix(tag, "v") {
		tag = "v" + tag
	}
	return u.release(ctx, "/releases/tags/"+tag)
}

func (u *Updater) release(ctx context.Context, path string) (*Release, error) {
	url := u.apiBase() + "/repos/" + u.repo() + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if u.Token != "" {
		req.Header.Set("Authorization", "Bearer "+u.Token)
	}
	resp, err := u.client().Do(req)
	if err != nil {
		return nil, fmt.Errorf("query %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no such release (%s)", url)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("query %s: %s: %s", url, resp.Status, strings.TrimSpace(string(body)))
	}
	var payload struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name string `json:"name"`
			URL  string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode release: %w", err)
	}
	rel := &Release{Tag: payload.TagName, Version: strings.TrimPrefix(payload.TagName, "v"), Assets: map[string]string{}}
	for _, a := range payload.Assets {
		rel.Assets[a.Name] = a.URL
	}
	return rel, nil
}

// ArchiveName is the GoReleaser archive for this platform, e.g. mcp-local_0.1.9_darwin_arm64.tar.gz.
func (u *Updater) ArchiveName(version string) string {
	goos, goarch := u.platform()
	return fmt.Sprintf("%s_%s_%s_%s.tar.gz", binaryName, version, goos, goarch)
}

// Install downloads rel's archive for this platform, checks it against checksums.txt, and
// atomically replaces the file at exe. exe is left untouched on any error.
func (u *Updater) Install(ctx context.Context, rel *Release, exe string) error {
	name := u.ArchiveName(rel.Version)
	archiveURL, ok := rel.Assets[name]
	if !ok {
		return fmt.Errorf("release %s has no %s (no build for this platform)", rel.Tag, name)
	}
	sumsURL, ok := rel.Assets["checksums.txt"]
	if !ok {
		return fmt.Errorf("release %s has no checksums.txt; refusing to install an unverified binary", rel.Tag)
	}
	sums, err := u.download(ctx, sumsURL)
	if err != nil {
		return err
	}
	want, err := checksumFor(sums, name)
	if err != nil {
		return err
	}
	archive, err := u.download(ctx, archiveURL)
	if err != nil {
		return err
	}
	got := sha256.Sum256(archive)
	if hex.EncodeToString(got[:]) != want {
		return fmt.Errorf("checksum mismatch for %s: got %x, checksums.txt says %s", name, got, want)
	}
	bin, err := extractBinary(archive)
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return u.replace(exe, bin)
}

// download fetches url without credentials: release assets are public, and the token is
// only for the API host.
func (u *Updater) download(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := u.client().Do(req)
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download %s: %s", url, resp.Status)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxDownload+1))
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", url, err)
	}
	if len(data) > maxDownload {
		return nil, fmt.Errorf("download %s: larger than %d bytes", url, maxDownload)
	}
	return data, nil
}

// checksumFor finds name's SHA-256 in a GoReleaser checksums.txt ("<hex>  <name>" lines).
func checksumFor(sums []byte, name string) (string, error) {
	sc := bufio.NewScanner(bytes.NewReader(sums))
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) == 2 && fields[1] == name {
			return strings.ToLower(fields[0]), nil
		}
	}
	return "", fmt.Errorf("checksums.txt has no entry for %s", name)
}

// extractBinary returns the mcp-local file from a .tar.gz archive.
func extractBinary(archive []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("no %s binary in archive", binaryName)
		}
		if err != nil {
			return nil, err
		}
		if hdr.Typeflag == tar.TypeReg && filepath.Base(hdr.Name) == binaryName {
			return io.ReadAll(io.LimitReader(tr, maxDownload))
		}
	}
}

// replace writes bin next to exe, checks it runs, and renames it over exe. Renaming keeps
// running copies (e.g. mcp-local bridge processes) on the old file until they restart.
func (u *Updater) replace(exe string, bin []byte) error {
	dir := filepath.Dir(exe)
	mode := os.FileMode(0o755)
	if fi, err := os.Stat(exe); err == nil {
		mode = fi.Mode().Perm() | 0o111
	}
	tmp, err := os.CreateTemp(dir, "."+binaryName+"-update-*")
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return fmt.Errorf("cannot write to %s (re-run with sudo, or install mcp-local to a directory you own): %w", dir, err)
		}
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath) // no-op after a successful rename
	if _, err := tmp.Write(bin); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpPath, mode); err != nil {
		return err
	}
	verify := u.Verify
	if verify == nil {
		verify = runVersion
	}
	if err := verify(tmpPath); err != nil {
		return fmt.Errorf("new binary failed its check, keeping the current one: %w", err)
	}
	return os.Rename(tmpPath, exe)
}

func runVersion(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, path, "--version").CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s --version: %w: %s", filepath.Base(path), err, strings.TrimSpace(string(out)))
	}
	return nil
}

// Executable is the real path of the running binary, with symlinks resolved so the
// target is replaced rather than the link.
func Executable() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(exe)
}

var versionRE = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(.*)$`)

// ErrDevBuild means the running binary has no release version to compare against.
var ErrDevBuild = errors.New("development build")

// Compare reports whether latest is newer than current. Versions are X.Y.Z with an optional
// suffix; a suffix on current (git describe output like 0.1.8-3-gabc123-dirty, from a
// source build) means "at or after X.Y.Z", so it is never offered the same release.
func Compare(current, latest string) (newer bool, err error) {
	c, err := parseVersion(current)
	if err != nil {
		return false, err
	}
	l, err := parseVersion(latest)
	if err != nil {
		return false, fmt.Errorf("release version: %w", err)
	}
	for i := range c {
		if l[i] != c[i] {
			return l[i] > c[i], nil
		}
	}
	return false, nil
}

func parseVersion(v string) ([3]int, error) {
	var out [3]int
	m := versionRE.FindStringSubmatch(strings.TrimSpace(v))
	if m == nil {
		return out, fmt.Errorf("%w (version %q)", ErrDevBuild, v)
	}
	for i := range out {
		out[i], _ = strconv.Atoi(m[i+1])
	}
	return out, nil
}
