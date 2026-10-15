package safefile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

// The ledger records which entry names mcp-local wrote into which file (and block, or
// "scope"), so removal can leave alone entries the user wrote that happen to share a name.
// Keys are resolved paths.
func ledgerPath() string { return filepath.Join(StateDir(), "written-entries.json") }

type ledger map[string]map[string][]string

func loadLedger() ledger {
	l := ledger{}
	if raw, err := os.ReadFile(ledgerPath()); err == nil {
		_ = json.Unmarshal(raw, &l)
	}
	return l
}

// Recorded reports whether mcp-local wrote name into scope of the file at path.
func Recorded(path, scope, name string) bool {
	for _, n := range loadLedger()[ResolvePath(path)][scope] {
		if n == name {
			return true
		}
	}
	return false
}

// Record adds and removes names for scope of the file at path. It writes only on change.
func Record(path, scope string, add, remove []string) error {
	if len(add) == 0 && len(remove) == 0 {
		return nil
	}
	l := loadLedger()
	key := ResolvePath(path)
	set := map[string]bool{}
	for _, n := range l[key][scope] {
		set[n] = true
	}
	changed := false
	for _, n := range add {
		changed = changed || !set[n]
		set[n] = true
	}
	for _, n := range remove {
		changed = changed || set[n]
		delete(set, n)
	}
	if !changed {
		return nil
	}
	names := make([]string, 0, len(set))
	for n := range set {
		names = append(names, n)
	}
	sort.Strings(names)
	if l[key] == nil {
		l[key] = map[string][]string{}
	}
	if len(names) == 0 {
		delete(l[key], scope)
	} else {
		l[key][scope] = names
	}
	if len(l[key]) == 0 {
		delete(l, key)
	}
	b, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	return AtomicWrite(ResolvePath(ledgerPath()), append(b, '\n'))
}
