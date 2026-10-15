// Package jsonagent edits the MCP server block of an agent's JSON/JSONC config file.
//
// Files are parsed with hujson and only the managed key path is patched (RFC 6902),
// so comments, key order, trailing commas, and formatting elsewhere survive byte-for-byte.
// Within an existing entry only the keys mcp-local manages are rewritten; the user's own keys
// (enabled, disabled, timeout, ...) are kept. Every modification backs up the original under
// ~/.mcp-local/backups and is written atomically; the names written are recorded in a ledger
// so removal can tell mcp-local's entries from the user's.
package jsonagent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/safefile"
	"github.com/tailscale/hujson"
)

type EntryConverter func(entry config.AgentEntry) map[string]interface{}

type Agent struct {
	configPath   string
	blockKey     string
	convertEntry EntryConverter
	managed      []string
}

// New returns an agent editing blockKey in configPath. managed lists the entry keys mcp-local
// owns: they are overwritten, and removed when the new entry no longer has them. Other keys the
// converter emits are defaults, written only when the entry lacks them. With no managed keys,
// every key the converter emits is managed and nothing is removed.
func New(configPath string, blockKey string, convert EntryConverter, managed ...string) Agent {
	return Agent{configPath: configPath, blockKey: blockKey, convertEntry: convert, managed: managed}
}

func (a Agent) ConfigPath() string { return a.configPath }

// Has reports whether the block has an entry named name.
func (a Agent) Has(name string) bool {
	m, err := ReadJSON(a.configPath)
	if err != nil {
		return false
	}
	b, _ := m[a.blockKey].(map[string]interface{})
	_, ok := b[name]
	return ok
}

// Owns reports whether mcp-local wrote the entry name (per the ledger).
func (a Agent) Owns(name string) bool { return safefile.Recorded(a.configPath, a.blockKey, name) }

// ParseJSONC decodes JSON or JSONC (comments, trailing commas) into a generic map.
func ParseJSONC(raw []byte) (map[string]interface{}, error) {
	std, err := hujson.Standardize(bytes.TrimPrefix(raw, bom))
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(std, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// ReadJSON reads a JSON or JSONC file into a generic map.
func ReadJSON(path string) (map[string]interface{}, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	m, err := ParseJSONC(raw)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return m, nil
}

// RegisterServices converts and upserts every service into the block. A service that can't be
// converted is reported but doesn't stop the others.
func (a Agent) RegisterServices(services []config.ServiceConfig) error {
	entries := make(map[string]map[string]interface{}, len(services))
	var errs []error
	for _, s := range services {
		entry, err := config.ServiceToEntry(s)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", s.Name, err))
			continue
		}
		entries[s.Name] = a.convertEntry(entry)
	}
	return errors.Join(append(errs, a.SetEntries(entries))...)
}

// RegisterRemote upserts a URL entry; an existing entry with the same url is left untouched.
func (a Agent) RegisterRemote(name, url string, extra map[string]interface{}) error {
	entry := map[string]interface{}{"url": url}
	for k, v := range extra {
		entry[k] = v
	}
	return a.edit(func(d *doc) error {
		if existing, ok := d.entry(name)["url"]; ok && existing == url {
			return nil
		}
		return d.set(name, entry)
	})
}

func (a Agent) RegisterLocal(name string, command []string, env map[string]string, envKey, argsKey string) error {
	if len(command) == 0 {
		return fmt.Errorf("empty command for %q", name)
	}
	entry := map[string]interface{}{"command": command[0]}
	if len(command) > 1 {
		entry[argsKey] = command[1:]
	}
	if len(env) > 0 {
		entry[envKey] = env
	}
	return a.SetEntries(map[string]map[string]interface{}{name: entry})
}

// SetEntries upserts the given entries (by exact name) into the block.
func (a Agent) SetEntries(entries map[string]map[string]interface{}) error {
	if len(entries) == 0 {
		return nil
	}
	names := make([]string, 0, len(entries))
	for n := range entries {
		names = append(names, n)
	}
	sort.Strings(names)
	return a.edit(func(d *doc) error {
		for _, n := range names {
			if err := d.set(n, entries[n]); err != nil {
				return err
			}
		}
		return nil
	})
}

// Deregister removes the entry with exactly this name. A missing file or entry is not an error.
func (a Agent) Deregister(name string) (bool, error) {
	if _, err := os.Stat(a.configPath); os.IsNotExist(err) {
		return false, nil
	}
	found := false
	err := a.edit(func(d *doc) error {
		var err error
		found, err = d.remove(name)
		return err
	})
	return found, err
}

var bom = []byte("\xef\xbb\xbf")

func (a Agent) edit(fn func(d *doc) error) error {
	path := resolvePath(a.configPath)
	raw, err := os.ReadFile(path)
	exists := err == nil
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	// Edit without the BOM and with LF endings; both are restored on write.
	src := bytes.TrimPrefix(raw, bom)
	hasBOM := len(src) != len(raw)
	crlf := bytes.Contains(src, []byte("\r\n"))
	if crlf {
		src = bytes.ReplaceAll(src, []byte("\r\n"), []byte("\n"))
	}
	if len(bytes.TrimSpace(src)) == 0 {
		src = []byte("{\n}\n")
	}
	v, err := hujson.Parse(src)
	if err != nil {
		return fmt.Errorf("parse %s: %w", a.configPath, err)
	}
	if _, ok := v.Value.(*hujson.Object); !ok {
		return fmt.Errorf("parse %s: top-level value is not an object", a.configPath)
	}
	d := &doc{v: &v, blockKey: a.blockKey, path: path, managed: a.managed}
	if err := duplicateKey(d.root()); err != nil {
		return fmt.Errorf("%s: %w", a.configPath, err)
	}
	if b := d.block(); b != nil {
		if err := duplicateKey(b); err != nil {
			return fmt.Errorf("%s: %s: %w", a.configPath, a.blockKey, err)
		}
	}
	if err := fn(d); err != nil {
		return fmt.Errorf("%s: %w", a.configPath, err)
	}
	if d.changed {
		out := v.Pack()
		if crlf {
			out = bytes.ReplaceAll(out, []byte("\n"), []byte("\r\n"))
		}
		if hasBOM {
			out = append(append([]byte{}, bom...), out...)
		}
		if exists {
			if err := safefile.Backup(a.configPath, raw); err != nil {
				return fmt.Errorf("backup %s: %w", a.configPath, err)
			}
		}
		if err := safefile.AtomicWrite(path, out); err != nil {
			return err
		}
		if err := d.saveOwnership(); err != nil {
			return err
		}
	}
	return safefile.Record(a.configPath, a.blockKey, d.written, d.removed)
}

// doc is a parsed config plus helpers that patch only /<blockKey>/<name> paths.
type doc struct {
	v        *hujson.Value
	blockKey string
	path     string
	managed  []string
	changed  bool
	created  *bool    // block ownership change to persist; nil = none
	written  []string // entry names set (for the ledger)
	removed  []string // entry names removed
}

// duplicateKey rejects an object with a repeated key: patches would edit the first while
// editors and JSON parsers read the last.
func duplicateKey(obj *hujson.Object) error {
	seen := map[string]bool{}
	for _, m := range obj.Members {
		lit, ok := m.Name.Value.(hujson.Literal)
		if !ok {
			continue
		}
		if seen[lit.String()] {
			return fmt.Errorf("duplicate key %q; remove one copy and retry", lit.String())
		}
		seen[lit.String()] = true
	}
	return nil
}

var pointerEscaper = strings.NewReplacer("~", "~0", "/", "~1")

func pointer(parts ...string) string {
	var b strings.Builder
	for _, p := range parts {
		b.WriteByte('/')
		b.WriteString(pointerEscaper.Replace(p))
	}
	return b.String()
}

func (d *doc) root() *hujson.Object { return d.v.Value.(*hujson.Object) }

func (d *doc) block() *hujson.Object {
	if bv := d.v.Find(pointer(d.blockKey)); bv != nil {
		obj, _ := bv.Value.(*hujson.Object)
		return obj
	}
	return nil
}

func (d *doc) entry(name string) map[string]interface{} {
	ev := d.v.Find(pointer(d.blockKey, name))
	if ev == nil {
		return nil
	}
	var m map[string]interface{}
	_ = decodeValue(*ev, &m)
	return m
}

func decodeValue(v hujson.Value, out interface{}) error {
	c := v.Clone()
	c.Standardize()
	return json.Unmarshal(c.Pack(), out)
}

func (d *doc) patch(op, path string, value []byte) error {
	var p bytes.Buffer
	pathJSON, _ := json.Marshal(path)
	fmt.Fprintf(&p, `[{"op":%q,"path":%s`, op, pathJSON)
	if value != nil {
		p.WriteString(`,"value":`)
		p.Write(value)
	}
	p.WriteString(`}]`)
	if err := d.v.Patch(p.Bytes()); err != nil {
		return err
	}
	d.changed = true
	return nil
}

func (d *doc) ensureBlock() error {
	if bv := d.v.Find(pointer(d.blockKey)); bv != nil {
		if _, ok := bv.Value.(*hujson.Object); !ok {
			return fmt.Errorf("%q is not an object", d.blockKey)
		}
		return nil
	}
	unit := indentUnit(d.root())
	trailing := hasTrailingComma(d.root())
	if err := d.patch("add", pointer(d.blockKey), []byte("{}")); err != nil {
		return err
	}
	fixInsert(d.root(), "", unit, trailing)
	created := true
	d.created = &created
	return nil
}

func (d *doc) set(name string, entry map[string]interface{}) error {
	d.written = append(d.written, name)
	want, err := normalize(entry)
	if err != nil {
		return err
	}
	if ev := d.v.Find(pointer(d.blockKey, name)); ev != nil {
		if obj, ok := ev.Value.(*hujson.Object); ok {
			return d.merge(name, ev, obj, entry)
		}
		var have interface{}
		if decodeValue(*ev, &have) == nil && reflect.DeepEqual(have, want) {
			return nil
		}
	}
	if err := d.ensureBlock(); err != nil {
		return err
	}
	unit := indentUnit(d.root())
	blockIndent := memberIndent(d.root(), d.blockKey, unit)
	obj := d.block()
	entryIndent := firstIndent(obj, blockIndent+unit)
	raw, err := json.MarshalIndent(entry, entryIndent, unit)
	if err != nil {
		return err
	}
	isNew := d.v.Find(pointer(d.blockKey, name)) == nil
	trailing := hasTrailingComma(obj)
	if err := d.patch("add", pointer(d.blockKey, name), raw); err != nil {
		return err
	}
	if isNew {
		fixInsert(obj, blockIndent, strings.TrimPrefix(entryIndent, blockIndent), trailing)
	}
	return nil
}

func (d *doc) isManaged(key string) bool {
	if len(d.managed) == 0 {
		return true
	}
	for _, k := range d.managed {
		if k == key {
			return true
		}
	}
	return false
}

// merge updates an existing entry object in place: managed keys are set (or removed when
// entry lacks them), other keys of entry are added only when missing, and the rest is kept.
func (d *doc) merge(name string, ev *hujson.Value, obj *hujson.Object, entry map[string]interface{}) error {
	if err := duplicateKey(obj); err != nil {
		return fmt.Errorf("%s/%s: %w", d.blockKey, name, err)
	}
	var have map[string]interface{}
	if err := decodeValue(*ev, &have); err != nil {
		return err
	}
	keys := make([]string, 0, len(entry))
	for k := range entry {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		want, err := normalize(entry[k])
		if err != nil {
			return err
		}
		cur, present := have[k]
		if present && (!d.isManaged(k) || reflect.DeepEqual(cur, want)) {
			continue
		}
		if err := d.setMember(name, obj, k, entry[k], present); err != nil {
			return err
		}
	}
	for _, k := range d.managed {
		if _, keep := entry[k]; keep {
			continue
		}
		if _, present := have[k]; present {
			if err := d.removeMember(obj, pointer(d.blockKey, name, k), memberIndent(d.block(), name, "")); err != nil {
				return err
			}
		}
	}
	return nil
}

// setMember writes key in the entry object, matching its layout (multi-line or one line).
func (d *doc) setMember(name string, obj *hujson.Object, key string, value interface{}, present bool) error {
	unit := indentUnit(d.root())
	blockIndent := memberIndent(d.root(), d.blockKey, unit)
	entryIndent := memberIndent(d.block(), name, blockIndent+unit)
	multiline := len(obj.Members) == 0 || bytes.Contains(obj.Members[0].Name.BeforeExtra, []byte("\n"))
	keyIndent := firstIndent(obj, entryIndent+unit)
	var raw []byte
	var err error
	if multiline {
		raw, err = json.MarshalIndent(value, keyIndent, unit)
	} else {
		raw, err = json.Marshal(value)
	}
	if err != nil {
		return err
	}
	var before, after hujson.Extra
	if present {
		for _, m := range obj.Members {
			if lit, ok := m.Name.Value.(hujson.Literal); ok && lit.String() == key {
				before, after = m.Value.BeforeExtra, m.Value.AfterExtra
			}
		}
	}
	trailing := hasTrailingComma(obj)
	if err := d.patch("add", pointer(d.blockKey, name, key), raw); err != nil {
		return err
	}
	for i := range obj.Members {
		m := &obj.Members[i]
		if lit, ok := m.Name.Value.(hujson.Literal); !ok || lit.String() != key {
			continue
		}
		switch {
		case present:
			m.Value.BeforeExtra, m.Value.AfterExtra = before, after
		case multiline:
			fixInsert(obj, entryIndent, strings.TrimPrefix(keyIndent, entryIndent), trailing)
		default:
			m.Name.BeforeExtra, m.Value.BeforeExtra = hujson.Extra(" "), hujson.Extra(" ")
		}
	}
	return nil
}

// removeMember removes the member at path from obj, keeping obj's trailing-comma style and
// closing-brace indentation.
func (d *doc) removeMember(obj *hujson.Object, path, closeIndent string) error {
	trailing := hasTrailingComma(obj)
	if err := d.patch("remove", path, nil); err != nil {
		return err
	}
	if n := len(obj.Members); trailing && n > 0 && obj.Members[n-1].Value.AfterExtra == nil {
		obj.Members[n-1].Value.AfterExtra = hujson.Extra{}
	}
	if _, ok := indentOf(obj.AfterExtra); ok {
		obj.AfterExtra = append(bytes.TrimRight(obj.AfterExtra, " \t"), closeIndent...)
	}
	return nil
}

func (d *doc) remove(name string) (bool, error) {
	obj := d.block()
	if obj == nil || d.v.Find(pointer(d.blockKey, name)) == nil {
		return false, nil
	}
	d.removed = append(d.removed, name)
	if err := d.removeMember(obj, pointer(d.blockKey, name), memberIndent(d.root(), d.blockKey, "")); err != nil {
		return false, err
	}
	if len(obj.Members) > 0 {
		return true, nil
	}
	if createdBlock(d.path, d.blockKey) {
		f := false
		d.created = &f
		return true, d.patch("remove", pointer(d.blockKey), nil)
	}
	if len(bytes.TrimSpace(obj.AfterExtra)) == 0 {
		obj.AfterExtra = nil
	}
	return true, nil
}

func (d *doc) saveOwnership() error {
	if d.created == nil {
		return nil
	}
	return setCreatedBlock(d.path, d.blockKey, *d.created)
}

func normalize(v interface{}) (interface{}, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var out interface{}
	return out, json.Unmarshal(b, &out)
}

func hasTrailingComma(obj *hujson.Object) bool {
	n := len(obj.Members)
	return n > 0 && obj.Members[n-1].Value.AfterExtra != nil
}

// fixInsert lays out a member just appended to obj: newline + indent before its name,
// a newline before the closing brace, and the object's existing trailing-comma style.
func fixInsert(obj *hujson.Object, closeIndent, unit string, trailingComma bool) {
	last := &obj.Members[len(obj.Members)-1]
	last.Name.BeforeExtra = withIndent(last.Name.BeforeExtra, closeIndent+unit)
	last.Value.BeforeExtra = hujson.Extra(" ")
	if trailingComma {
		last.Value.AfterExtra = hujson.Extra{}
	}
	if !bytes.Contains(obj.AfterExtra, []byte("\n")) && len(bytes.TrimSpace(obj.AfterExtra)) == 0 {
		obj.AfterExtra = hujson.Extra("\n" + closeIndent)
	}
}

func withIndent(b hujson.Extra, indent string) hujson.Extra {
	out := bytes.TrimRight(b, " \t")
	if !bytes.HasSuffix(out, []byte("\n")) {
		out = append(out, '\n')
	}
	return append(out, indent...)
}

func indentOf(b hujson.Extra) (string, bool) {
	i := bytes.LastIndexByte(b, '\n')
	if i < 0 {
		return "", false
	}
	tail := b[i+1:]
	if len(bytes.Trim(tail, " \t")) != 0 {
		return "", false
	}
	return string(tail), true
}

// indentUnit is the indentation of the first top-level member (default two spaces).
func indentUnit(root *hujson.Object) string {
	return firstIndent(root, "  ")
}

func memberIndent(obj *hujson.Object, name, fallback string) string {
	for _, m := range obj.Members {
		if lit, ok := m.Name.Value.(hujson.Literal); ok && lit.String() == name {
			if s, ok := indentOf(m.Name.BeforeExtra); ok {
				return s
			}
		}
	}
	return fallback
}

func firstIndent(obj *hujson.Object, fallback string) string {
	if len(obj.Members) > 0 {
		if s, ok := indentOf(obj.Members[0].Name.BeforeExtra); ok && s != "" {
			return s
		}
	}
	return fallback
}
