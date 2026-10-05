// Package jsonagent edits the MCP server block of an agent's JSON/JSONC config file.
//
// Files are parsed with hujson and only the managed key path is patched (RFC 6902),
// so comments, key order, trailing commas, and formatting elsewhere survive byte-for-byte.
// Every modification backs up the original under ~/.mcp-local/backups and is written atomically.
package jsonagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/tailscale/hujson"
)

type EntryConverter func(entry config.AgentEntry) map[string]interface{}

type Agent struct {
	configPath   string
	blockKey     string
	convertEntry EntryConverter
}

func New(configPath string, blockKey string, convert EntryConverter) Agent {
	return Agent{configPath: configPath, blockKey: blockKey, convertEntry: convert}
}

func (a Agent) ConfigPath() string { return a.configPath }

// ParseJSONC decodes JSON or JSONC (comments, trailing commas) into a generic map.
func ParseJSONC(raw []byte) (map[string]interface{}, error) {
	std, err := hujson.Standardize(raw)
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

// RegisterServices converts and upserts every service into the block.
func (a Agent) RegisterServices(services []config.ServiceConfig) error {
	entries := make(map[string]map[string]interface{}, len(services))
	for _, s := range services {
		entry, err := config.ServiceToEntry(s)
		if err != nil {
			return fmt.Errorf("%s: %w", s.Name, err)
		}
		entries[s.Name] = a.convertEntry(entry)
	}
	return a.SetEntries(entries)
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

func (a Agent) edit(fn func(d *doc) error) error {
	path := resolvePath(a.configPath)
	raw, err := os.ReadFile(path)
	exists := err == nil
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	src := raw
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
	d := &doc{v: &v, blockKey: a.blockKey, path: path}
	if err := fn(d); err != nil {
		return fmt.Errorf("%s: %w", a.configPath, err)
	}
	if !d.changed {
		return nil
	}
	if exists {
		if err := backup(a.configPath, raw); err != nil {
			return fmt.Errorf("backup %s: %w", a.configPath, err)
		}
	}
	if err := atomicWrite(path, v.Pack()); err != nil {
		return err
	}
	return d.saveOwnership()
}

// doc is a parsed config plus helpers that patch only /<blockKey>/<name> paths.
type doc struct {
	v        *hujson.Value
	blockKey string
	path     string
	changed  bool
	created  *bool // block ownership change to persist; nil = none
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
	want, err := normalize(entry)
	if err != nil {
		return err
	}
	if ev := d.v.Find(pointer(d.blockKey, name)); ev != nil {
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

func (d *doc) remove(name string) (bool, error) {
	obj := d.block()
	if obj == nil || d.v.Find(pointer(d.blockKey, name)) == nil {
		return false, nil
	}
	closeIndent := memberIndent(d.root(), d.blockKey, "")
	trailing := hasTrailingComma(obj)
	if err := d.patch("remove", pointer(d.blockKey, name), nil); err != nil {
		return false, err
	}
	if n := len(obj.Members); trailing && n > 0 && obj.Members[n-1].Value.AfterExtra == nil {
		obj.Members[n-1].Value.AfterExtra = hujson.Extra{}
	}
	if _, ok := indentOf(obj.AfterExtra); ok {
		obj.AfterExtra = append(bytes.TrimRight(obj.AfterExtra, " \t"), closeIndent...)
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
