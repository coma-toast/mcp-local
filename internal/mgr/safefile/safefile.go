// Package safefile writes files that users and other tools also edit: atomically (temp file in
// the target's directory, fsync, rename), through symlinks, keeping the file mode, with a
// timestamped backup of the previous contents under ~/.mcp-local/backups.
package safefile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// MaxBackups is how many timestamped backups are kept per path.
const MaxBackups = 5

// Now is swappable in tests.
var Now = time.Now

// StateDir is ~/.mcp-local (resolved per call so tests can redirect HOME).
func StateDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mcp-local")
}

func BackupDir() string { return filepath.Join(StateDir(), "backups") }

// ResolvePath follows a symlinked file so the link survives the atomic rename. A file that
// doesn't exist yet still gets its directory resolved, so the path is the same before and after
// the first write (e.g. macOS /var -> /private/var). A dangling symlink resolves to its target,
// so the first write creates the target instead of replacing the link.
func ResolvePath(path string) string {
	return resolve(path, 8)
}

func resolve(path string, hops int) string {
	if real, err := filepath.EvalSymlinks(path); err == nil {
		return real
	}
	if target, err := os.Readlink(path); err == nil && hops > 0 {
		if !filepath.IsAbs(target) {
			target = filepath.Join(filepath.Dir(path), target)
		}
		return resolve(target, hops-1)
	}
	if dir, err := filepath.EvalSymlinks(filepath.Dir(path)); err == nil {
		return filepath.Join(dir, filepath.Base(path))
	}
	return path
}

// BackupName is the file name a backup of path gets: the path with "/" replaced by "%".
func BackupName(path string) string {
	return strings.ReplaceAll(filepath.ToSlash(path), "/", "%")
}

// Backup copies the original bytes to ~/.mcp-local/backups/<YYYYMMDD-HHMMSS>/<BackupName(path)>
// and prunes all but the newest MaxBackups for that path.
func Backup(path string, original []byte) error {
	name := BackupName(path)
	dir := filepath.Join(BackupDir(), Now().Format("20060102-150405"))
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

// AtomicWrite writes via a temp file in the same dir (fsync + rename), preserving the original
// mode. path must already be resolved (see ResolvePath), or a symlink there is replaced.
func AtomicWrite(path string, data []byte) error {
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
	if err := rename(tmpName, path); err != nil {
		return fmt.Errorf("replace %s: %w", path, err)
	}
	if d, err := os.Open(dir); err == nil {
		_ = d.Sync()
		d.Close()
	}
	return nil
}

// rename is swappable in tests.
var rename = os.Rename

// Replace backs up the current contents of path (if any) and atomically writes data to the
// file path points at, so a symlink at path is kept.
func Replace(path string, data []byte) error {
	real := ResolvePath(path)
	if old, err := os.ReadFile(real); err == nil {
		if err := Backup(path, old); err != nil {
			return fmt.Errorf("backup %s: %w", path, err)
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	return AtomicWrite(real, data)
}
