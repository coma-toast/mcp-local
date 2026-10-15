package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/coma-toast/mcp-local/internal/mgr/safefile"
	"gopkg.in/yaml.v3"
)

// SaveConfig writes cfg to config.yaml without disturbing what it does not change: values
// that still match the file's ${HOME}/~ form are written in that form, the existing YAML is
// edited in place (comments and unknown keys survive), and the file is backed up and replaced
// atomically through any symlink. A save that changes nothing is skipped.
func SaveConfig(cfg *ManagerConfig) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}
	doc, raw, agents := cfg.doc, cfg.raw, cfg.Agents
	var current []byte
	if data, err := os.ReadFile(path); err == nil {
		current = data
		if doc == nil {
			// cfg wasn't loaded from this file (e.g. the wizard): merge into what is there,
			// and keep its agents: block unless cfg sets one.
			onDisk, err := parseConfig(data)
			if err != nil {
				return fmt.Errorf("parse %s: %w", path, err)
			}
			doc, raw = onDisk.doc, onDisk.raw
			if agents == nil {
				agents = onDisk.Agents
			}
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	out := cfg.persistable(raw)
	out.Agents = agents
	if current != nil && sameConfig(current, out) {
		return nil
	}
	var fresh yaml.Node
	if err := fresh.Encode(out); err != nil {
		return err
	}
	if doc == nil || len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		doc = &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{&fresh}}
	} else {
		mergeNode(doc.Content[0], &fresh, reflect.TypeOf(ManagerConfig{}))
	}
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(detectIndent(current))
	if err := enc.Encode(doc); err != nil {
		return err
	}
	if err := enc.Close(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := safefile.Replace(path, buf.Bytes()); err != nil {
		return err
	}
	cfg.doc, cfg.raw = doc, byName(out.Services)
	return nil
}

// sameConfig reports whether the file already holds out (compared in normalized form).
func sameConfig(current []byte, out ManagerConfig) bool {
	var have ManagerConfig
	if yaml.Unmarshal(current, &have) != nil {
		return false
	}
	a, errA := yaml.Marshal(have)
	b, errB := yaml.Marshal(out)
	return errA == nil && errB == nil && bytes.Equal(a, b)
}

// persistable is cfg as it should be written: each value that still equals the expansion of
// what the file had keeps the file's form (e.g. ${HOME}/bin rather than /Users/me/bin).
func (m *ManagerConfig) persistable(raw map[string]ServiceConfig) ManagerConfig {
	out := ManagerConfig{Agents: m.Agents}
	for _, s := range m.Services {
		s.Deps = append([]string(nil), s.Deps...)
		if s.Env != nil {
			env := make(map[string]string, len(s.Env))
			for k, v := range s.Env {
				env[k] = v
			}
			s.Env = env
		}
		if r, ok := raw[s.Name]; ok {
			unexpand(&s, r)
		}
		out.Services = append(out.Services, s)
	}
	return out
}

func unexpand(s *ServiceConfig, raw ServiceConfig) {
	keep := func(cur *string, orig string) {
		if *cur != orig && ExpandPath(orig) == *cur {
			*cur = orig
		}
	}
	keep(&s.Command, raw.Command)
	keep(&s.HealthURL, raw.HealthURL)
	keep(&s.DashboardURL, raw.DashboardURL)
	keep(&s.MCPURL, raw.MCPURL)
	keep(&s.Log, raw.Log)
	keep(&s.Path, raw.Path)
	keep(&s.ToolsConfigPath, raw.ToolsConfigPath)
	for i := range s.Deps {
		if i < len(raw.Deps) {
			keep(&s.Deps[i], raw.Deps[i])
		}
	}
	for k, v := range s.Env {
		if orig, ok := raw.Env[k]; ok {
			keep(&v, orig)
			s.Env[k] = v
		}
	}
}

// mergeNode edits old in place so it decodes like nw, keeping old's comments, styles and
// any keys whose presence doesn't change the decoded value (unknown keys, explicit defaults).
// t is the Go type both nodes decode into (nil when unknown).
func mergeNode(old, nw *yaml.Node, t reflect.Type) {
	if nodesEqual(old, nw) {
		return
	}
	for t != nil && t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	switch {
	case t == nil || old.Kind != nw.Kind:
		replaceNode(old, nw)
	case old.Kind == yaml.MappingNode && (t.Kind() == reflect.Struct || t.Kind() == reflect.Map):
		mergeMapping(old, nw, t)
	case old.Kind == yaml.SequenceNode && t.Kind() == reflect.Slice && hasNameField(t.Elem()):
		mergeNamedSeq(old, nw, t.Elem())
	default:
		replaceNode(old, nw)
	}
}

func mergeMapping(old, nw *yaml.Node, t reflect.Type) {
	idx := map[string]int{}
	for i := 0; i+1 < len(old.Content); i += 2 {
		idx[old.Content[i].Value] = i
	}
	inNew := map[string]bool{}
	for i := 0; i+1 < len(nw.Content); i += 2 {
		k := nw.Content[i].Value
		inNew[k] = true
		if j, ok := idx[k]; ok {
			mergeNode(old.Content[j+1], nw.Content[i+1], fieldType(t, k))
		} else {
			old.Content = append(old.Content, nw.Content[i], nw.Content[i+1])
		}
	}
	for i := 0; i+1 < len(old.Content); {
		if inNew[old.Content[i].Value] || !changesValue(old, i, t) {
			i += 2
			continue
		}
		old.Content = append(old.Content[:i], old.Content[i+2:]...)
	}
}

// changesValue reports whether dropping the key at index i of mapping m changes what m
// decodes to as t. Unknown keys and keys set to their default do not.
func changesValue(m *yaml.Node, i int, t reflect.Type) bool {
	without := *m
	without.Content = append(append([]*yaml.Node{}, m.Content[:i]...), m.Content[i+2:]...)
	a, b := reflect.New(t), reflect.New(t)
	if m.Decode(a.Interface()) != nil || without.Decode(b.Interface()) != nil {
		return true
	}
	return !reflect.DeepEqual(a.Elem().Interface(), b.Elem().Interface())
}

func mergeNamedSeq(old, nw *yaml.Node, elem reflect.Type) {
	byName := map[string]*yaml.Node{}
	for _, item := range old.Content {
		if n := nameOf(item); n != "" {
			byName[n] = item
		}
	}
	content := make([]*yaml.Node, 0, len(nw.Content))
	for _, item := range nw.Content {
		if o, ok := byName[nameOf(item)]; ok {
			delete(byName, nameOf(item))
			mergeNode(o, item, elem)
			item = o
		}
		content = append(content, item)
	}
	old.Content = content
}

func nameOf(n *yaml.Node) string {
	if n.Kind != yaml.MappingNode {
		return ""
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		if n.Content[i].Value == "name" {
			return n.Content[i+1].Value
		}
	}
	return ""
}

func hasNameField(t reflect.Type) bool {
	return t.Kind() == reflect.Struct && fieldType(t, "name") != nil
}

// fieldType is the Go type of key in t: the field with that yaml name for a struct (nil if
// none), or the element type for a map.
func fieldType(t reflect.Type, key string) reflect.Type {
	if t.Kind() == reflect.Map {
		return t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	if t == reflect.TypeOf(ToolConfig{}) {
		t = reflect.TypeOf(toolYAML{})
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name, _, _ := strings.Cut(f.Tag.Get("yaml"), ",")
		if name == key && f.IsExported() {
			return f.Type
		}
	}
	return nil
}

func nodesEqual(a, b *yaml.Node) bool {
	var va, vb interface{}
	if a.Decode(&va) != nil || b.Decode(&vb) != nil {
		return false
	}
	return reflect.DeepEqual(va, vb)
}

func replaceNode(old, nw *yaml.Node) {
	head, line, foot := old.HeadComment, old.LineComment, old.FootComment
	*old = *nw
	old.HeadComment, old.LineComment, old.FootComment = head, line, foot
}

// detectIndent is the file's indentation step: its smallest non-zero leading indent (2 if none).
func detectIndent(data []byte) int {
	best := 0
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimLeft(line, " ")
		n := len(line) - len(trimmed)
		if n == 0 || trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if best == 0 || n < best {
			best = n
		}
	}
	if best < 2 || best > 8 {
		return 2
	}
	return best
}
