package jsonagent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// MaxBackups is how many timestamped backups are kept per config path.
const MaxBackups = 5

var now = time.Now

// StateDir is ~/.mcp-local (resolved per call so tests can redirect HOME).
func StateDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mcp-local")
}

func BackupDir() string { return filepath.Join(StateDir(), "backups") }

// resolvePath follows a symlinked config so the link survives the atomic rename.
func resolvePath(path string) string {
	if real, err := filepath.EvalSymlinks(path); err == nil {
		return real
	}
	return path
}

func backupName(path string) string {
	return strings.ReplaceAll(filepath.ToSlash(path), "/", "%")
}

// backup copies the original bytes to ~/.mcp-local/backups/<YYYYMMDD-HHMMSS>/<path with / → %>
// and prunes all but the newest MaxBackups for that path.
func backup(path string, original []byte) error {
	name := backupName(path)
	dir := filepath.Join(BackupDir(), now().Format("20060102-150405"))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	dst := filepath.Join(dir, name)
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		if err := os.WriteFile(dst, original, 0o600); err != nil {
			return err
		}
	}
	return pruneBackups(name)
}

func pruneBackups(name string) error {
	entries, err := os.ReadDir(BackupDir())
	if err != nil {
		return err
	}
	var dirs []string
	for _, e := range entries {
		if e.IsDir() {
			if _, err := os.Stat(filepath.Join(BackupDir(), e.Name(), name)); err == nil {
				dirs = append(dirs, e.Name())
			}
		}
	}
	sort.Strings(dirs)
	for len(dirs) > MaxBackups {
		dir := filepath.Join(BackupDir(), dirs[0])
		if err := os.Remove(filepath.Join(dir, name)); err != nil {
			return err
		}
		_ = os.Remove(dir) // only succeeds once the timestamp dir is empty
		dirs = dirs[1:]
	}
	return nil
}

// atomicWrite writes via a temp file in the same dir (fsync + rename), preserving the original mode.
func atomicWrite(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	mode := os.FileMode(0o644)
	if fi, err := os.Stat(path); err == nil {
		mode = fi.Mode().Perm()
	}
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Chmod(mode); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("replace %s: %w", path, err)
	}
	if d, err := os.Open(dir); err == nil {
		_ = d.Sync()
		d.Close()
	}
	return nil
}

// Ownership: blocks mcp-local created (vs. ones the user already had) are recorded so that
// removing the last entry deletes a block we added but leaves a pre-existing one as {}.
func ownershipPath() string { return filepath.Join(StateDir(), "created-blocks.json") }

func loadOwnership() map[string][]string {
	m := map[string][]string{}
	if raw, err := os.ReadFile(ownershipPath()); err == nil {
		_ = json.Unmarshal(raw, &m)
	}
	return m
}

func createdBlock(path, blockKey string) bool {
	for _, k := range loadOwnership()[path] {
		if k == blockKey {
			return true
		}
	}
	return false
}

func setCreatedBlock(path, blockKey string, created bool) error {
	m := loadOwnership()
	var keys []string
	for _, k := range m[path] {
		if k != blockKey {
			keys = append(keys, k)
		}
	}
	if created {
		keys = append(keys, blockKey)
	}
	if len(keys) == 0 {
		delete(m, path)
	} else {
		m[path] = keys
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return atomicWrite(ownershipPath(), append(b, '\n'))
}
