package jsonagent

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/coma-toast/mcp-local/internal/mgr/safefile"
)

// MaxBackups is how many timestamped backups are kept per config path.
const MaxBackups = safefile.MaxBackups

func StateDir() string  { return safefile.StateDir() }
func BackupDir() string { return safefile.BackupDir() }

func resolvePath(path string) string { return safefile.ResolvePath(path) }
func backupName(path string) string  { return safefile.BackupName(path) }

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
	return safefile.AtomicWrite(resolvePath(ownershipPath()), append(b, '\n'))
}
