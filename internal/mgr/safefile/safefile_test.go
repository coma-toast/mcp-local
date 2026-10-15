package safefile

import (
	"os"
	"path/filepath"
	"testing"
)

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestReplaceKeepsSymlinkModeAndBacksUp(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	real := filepath.Join(t.TempDir(), "shared.yaml")
	link := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(real, []byte("old\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(real, link); err != nil {
		t.Skip(err)
	}
	if err := Replace(link, []byte("new\n")); err != nil {
		t.Fatal(err)
	}
	if fi, _ := os.Lstat(link); fi.Mode()&os.ModeSymlink == 0 {
		t.Fatal("symlink replaced by a regular file")
	}
	if got := readFile(t, real); got != "new\n" {
		t.Errorf("target = %q", got)
	}
	if fi, _ := os.Stat(real); fi.Mode().Perm() != 0o600 {
		t.Errorf("mode = %v, want 0600", fi.Mode().Perm())
	}
	backups, _ := filepath.Glob(filepath.Join(BackupDir(), "*", BackupName(link)))
	if len(backups) != 1 || readFile(t, backups[0]) != "old\n" {
		t.Errorf("backups = %v", backups)
	}
}

func TestAtomicWriteRenamesCompleteTempFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "f.json")
	if err := os.WriteFile(path, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	orig := rename
	defer func() { rename = orig }()
	var tmpDir, tmpData, targetAtRename string
	rename = func(from, to string) error {
		tmpDir, tmpData, targetAtRename = filepath.Dir(from), readFile(t, from), readFile(t, to)
		return orig(from, to)
	}
	if err := AtomicWrite(path, []byte("new contents")); err != nil {
		t.Fatal(err)
	}
	if tmpDir != filepath.Dir(path) || tmpData != "new contents" || targetAtRename != "old" {
		t.Errorf("rename from dir %q with %q while target held %q", tmpDir, tmpData, targetAtRename)
	}
	if left, _ := filepath.Glob(filepath.Join(filepath.Dir(path), ".f.json.tmp-*")); len(left) != 0 {
		t.Errorf("temp files left: %v", left)
	}
}

func TestResolvePathDanglingSymlink(t *testing.T) {
	dir := t.TempDir()
	link := filepath.Join(dir, "config.yaml")
	if err := os.Symlink("synced/missing.yaml", link); err != nil {
		t.Skip(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "synced"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := AtomicWrite(ResolvePath(link), []byte("x")); err != nil {
		t.Fatal(err)
	}
	if fi, _ := os.Lstat(link); fi.Mode()&os.ModeSymlink == 0 {
		t.Fatal("dangling symlink replaced instead of creating its target")
	}
	if readFile(t, filepath.Join(dir, "synced", "missing.yaml")) != "x" {
		t.Error("target not created")
	}
}

func TestLedger(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	path := filepath.Join(t.TempDir(), "mcp.json")
	if Recorded(path, "mcpServers", "a") {
		t.Fatal("empty ledger")
	}
	if err := Record(path, "mcpServers", []string{"a", "b"}, nil); err != nil {
		t.Fatal(err)
	}
	if err := Record(path, "mcpServers", nil, []string{"a"}); err != nil {
		t.Fatal(err)
	}
	if Recorded(path, "mcpServers", "a") || !Recorded(path, "mcpServers", "b") || Recorded(path, "mcp", "b") {
		t.Error("ledger contents wrong")
	}
}
