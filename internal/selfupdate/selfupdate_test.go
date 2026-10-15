package selfupdate

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompare(t *testing.T) {
	cases := []struct {
		current, latest string
		newer           bool
	}{
		{"0.1.8", "0.1.9", true},
		{"0.1.9", "0.1.9", false},
		{"0.2.0", "0.1.9", false},
		{"v0.1.8", "v0.10.0", true},
		{"0.1.8-3-gabc1234", "0.1.8", false}, // source build ahead of the tag
		{"0.1.8-3-gabc1234-dirty", "0.1.9", true},
	}
	for _, c := range cases {
		got, err := Compare(c.current, c.latest)
		if err != nil || got != c.newer {
			t.Errorf("Compare(%q, %q) = %v, %v; want %v", c.current, c.latest, got, err, c.newer)
		}
	}
	if _, err := Compare("dev", "0.1.9"); !errors.Is(err, ErrDevBuild) {
		t.Errorf("dev build: err = %v, want ErrDevBuild", err)
	}
}

func tarGz(t *testing.T, name string, body []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for _, f := range []struct {
		name string
		body []byte
	}{{"README.md", []byte("readme")}, {name, body}} {
		if err := tw.WriteHeader(&tar.Header{Name: f.name, Mode: 0o755, Size: int64(len(f.body)), Typeflag: tar.TypeReg}); err != nil {
			t.Fatal(err)
		}
		tw.Write(f.body)
	}
	tw.Close()
	gz.Close()
	return buf.Bytes()
}

// fakeGitHub serves a v0.1.9 release with a linux/amd64 archive. corrupt serves a
// different archive than the one checksums.txt describes.
func fakeGitHub(t *testing.T, bin []byte, corrupt bool) (*httptest.Server, *string) {
	t.Helper()
	archive := tarGz(t, "mcp-local", bin)
	sum := sha256.Sum256(archive)
	if corrupt {
		archive = tarGz(t, "mcp-local", []byte("tampered"))
	}
	name := "mcp-local_0.1.9_linux_amd64.tar.gz"
	var auth string
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/coma-toast/mcp-local/releases/latest", "/repos/coma-toast/mcp-local/releases/tags/v0.1.9":
			auth = r.Header.Get("Authorization")
			fmt.Fprintf(w, `{"tag_name":"v0.1.9","assets":[{"name":%q,"browser_download_url":"%s/dl/%s"},{"name":"checksums.txt","browser_download_url":"%s/dl/checksums.txt"}]}`,
				name, srv.URL, name, srv.URL)
		case "/dl/" + name:
			if r.Header.Get("Authorization") != "" {
				t.Errorf("token sent to asset download")
			}
			w.Write(archive)
		case "/dl/checksums.txt":
			fmt.Fprintf(w, "%x  %s\n%x  other.tar.gz\n", sum, name, sha256.Sum256(nil))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv, &auth
}

func newUpdater(srv *httptest.Server) *Updater {
	return &Updater{APIBase: srv.URL, Token: "tok", GOOS: "linux", GOARCH: "amd64", Verify: func(string) error { return nil }}
}

func oldBinary(t *testing.T) string {
	t.Helper()
	exe := filepath.Join(t.TempDir(), "mcp-local")
	if err := os.WriteFile(exe, []byte("old"), 0o755); err != nil {
		t.Fatal(err)
	}
	return exe
}

func TestInstallReplacesBinary(t *testing.T) {
	srv, auth := fakeGitHub(t, []byte("new binary"), false)
	u := newUpdater(srv)
	rel, err := u.Latest(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if rel.Version != "0.1.9" || *auth != "Bearer tok" {
		t.Fatalf("release %+v, auth %q", rel, *auth)
	}
	exe := oldBinary(t)
	if err := u.Install(context.Background(), rel, exe); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(exe)
	if string(got) != "new binary" {
		t.Errorf("binary = %q", got)
	}
	if fi, _ := os.Stat(exe); fi.Mode().Perm()&0o111 == 0 {
		t.Errorf("new binary not executable: %v", fi.Mode())
	}
	if left, _ := filepath.Glob(filepath.Join(filepath.Dir(exe), ".mcp-local-update-*")); len(left) != 0 {
		t.Errorf("temp files left behind: %v", left)
	}
}

func TestInstallKeepsOldBinaryOnFailure(t *testing.T) {
	cases := map[string]func(*Updater, *httptest.Server){
		"checksum mismatch": func(*Updater, *httptest.Server) {},
		"no build for platform": func(u *Updater, _ *httptest.Server) {
			u.GOARCH = "riscv64"
		},
		"new binary fails to run": func(u *Updater, _ *httptest.Server) {
			u.Verify = func(string) error { return errors.New("exec format error") }
		},
	}
	for name, setup := range cases {
		t.Run(name, func(t *testing.T) {
			srv, _ := fakeGitHub(t, []byte("new binary"), name == "checksum mismatch")
			u := newUpdater(srv)
			setup(u, srv)
			rel, err := u.ByTag(context.Background(), "0.1.9")
			if err != nil {
				t.Fatal(err)
			}
			exe := oldBinary(t)
			if err := u.Install(context.Background(), rel, exe); err == nil {
				t.Fatal("expected an error")
			}
			if got, _ := os.ReadFile(exe); string(got) != "old" {
				t.Errorf("old binary changed to %q", got)
			}
			if left, _ := filepath.Glob(filepath.Join(filepath.Dir(exe), ".mcp-local-update-*")); len(left) != 0 {
				t.Errorf("temp files left behind: %v", left)
			}
		})
	}
}

func TestMissingReleaseIsClear(t *testing.T) {
	srv, _ := fakeGitHub(t, nil, false)
	_, err := newUpdater(srv).ByTag(context.Background(), "v9.9.9")
	if err == nil || !strings.Contains(err.Error(), "no such release") {
		t.Errorf("err = %v", err)
	}
}
