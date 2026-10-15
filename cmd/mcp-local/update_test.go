package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/coma-toast/mcp-local/internal/selfupdate"
	"github.com/spf13/cobra"
)

// TestUpdateDecisions covers when update installs; the download itself is tested in selfupdate.
func TestUpdateDecisions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"tag_name":"v0.1.9","assets":[]}`)
	}))
	defer srv.Close()
	origUpdater, origTarget := newUpdater, updateTarget
	defer func() { newUpdater, updateTarget = origUpdater, origTarget }()
	newUpdater = func() *selfupdate.Updater {
		return &selfupdate.Updater{APIBase: srv.URL, GOOS: "linux", GOARCH: "amd64"}
	}

	cases := []struct {
		current     string
		check       bool
		force       bool
		wantInstall bool
		wantErr     string
		wantOut     string
	}{
		{current: "0.1.8", check: true, wantOut: "Update available"},
		{current: "0.1.9", wantOut: "Up to date"},
		{current: "0.1.9-2-gabc1234", wantOut: "Up to date"},
		{current: "dev", wantErr: "development build"},
		{current: "dev", check: true, wantOut: "development build"},
		{current: "0.1.8", wantInstall: true},
		{current: "0.1.9", force: true, wantInstall: true},
		{current: "dev", force: true, wantInstall: true},
	}
	for _, c := range cases {
		installed := false
		updateTarget = func() (string, error) { installed = true; return "", fmt.Errorf("stop before install") }
		cmd := &cobra.Command{}
		cmd.SetContext(context.Background())
		var out bytes.Buffer
		err := runUpdate(cmd, &out, c.current, "", c.check, c.force)
		name := fmt.Sprintf("current=%s check=%v force=%v", c.current, c.check, c.force)
		if installed != c.wantInstall {
			t.Errorf("%s: install attempted = %v, want %v", name, installed, c.wantInstall)
		}
		if c.wantErr != "" && (err == nil || !strings.Contains(err.Error(), c.wantErr)) {
			t.Errorf("%s: err = %v, want %q", name, err, c.wantErr)
		}
		if !c.wantInstall && c.wantErr == "" && err != nil {
			t.Errorf("%s: unexpected err %v", name, err)
		}
		if !strings.Contains(out.String(), c.wantOut) {
			t.Errorf("%s: output %q missing %q", name, out.String(), c.wantOut)
		}
	}
}
