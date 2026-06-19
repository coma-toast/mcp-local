package embedfs

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type toolDef struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

type toolHandlerFn func(root string, args json.RawMessage) (string, error)

func toolDefinitions() []toolDef {
	return []toolDef{
		{
			Name:        "read_file",
			Description: "Read the contents of a file",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Path relative to server root",
					},
				},
				"required": []string{"path"},
			},
		},
		{
			Name:        "write_file",
			Description: "Write content to a file (creates parent directories)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Path relative to server root",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "File content to write",
					},
				},
				"required": []string{"path", "content"},
			},
		},
		{
			Name:        "list_directory",
			Description: "List files and directories in a directory",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Directory path relative to server root",
					},
				},
				"required": []string{"path"},
			},
		},
		{
			Name:        "search_files",
			Description: "Search for files matching a glob pattern",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"pattern": map[string]interface{}{
						"type":        "string",
						"description": "Glob pattern (e.g. **/*.go, *.md)",
					},
				},
				"required": []string{"pattern"},
			},
		},
		{
			Name:        "get_file_info",
			Description: "Get metadata for a file or directory",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Path relative to server root",
					},
				},
				"required": []string{"path"},
			},
		},
		{
			Name:        "create_directory",
			Description: "Create a directory (including parents)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Directory path relative to server root",
					},
				},
				"required": []string{"path"},
			},
		},
	}
}

func toolHandler(name string) toolHandlerFn {
	switch name {
	case "read_file":
		return handleReadFile
	case "write_file":
		return handleWriteFile
	case "list_directory":
		return handleListDir
	case "search_files":
		return handleSearchFiles
	case "get_file_info":
		return handleGetFileInfo
	case "create_directory":
		return handleCreateDir
	}
	return nil
}

type toolArgs struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Pattern string `json:"pattern"`
}

func handleReadFile(root string, raw json.RawMessage) (string, error) {
	var a toolArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return "", fmt.Errorf("invalid args: %w", err)
	}
	full, err := safeJoin(root, a.Path)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(full)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", a.Path, err)
	}
	return string(data), nil
}

func handleWriteFile(root string, raw json.RawMessage) (string, error) {
	var a toolArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return "", fmt.Errorf("invalid args: %w", err)
	}
	full, err := safeJoin(root, a.Path)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return "", fmt.Errorf("mkdir: %w", err)
	}
	if err := os.WriteFile(full, []byte(a.Content), 0644); err != nil {
		return "", fmt.Errorf("write %s: %w", a.Path, err)
	}
	return fmt.Sprintf("Written %d bytes to %s", len(a.Content), a.Path), nil
}

func handleListDir(root string, raw json.RawMessage) (string, error) {
	var a toolArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return "", fmt.Errorf("invalid args: %w", err)
	}
	full, err := safeJoin(root, a.Path)
	if err != nil {
		return "", err
	}
	entries, err := os.ReadDir(full)
	if err != nil {
		return "", fmt.Errorf("list %s: %w", a.Path, err)
	}
	var b strings.Builder
	for _, e := range entries {
		info, _ := e.Info()
		mode := info.Mode().String()
		size := info.Size()
		name := e.Name()
		if e.IsDir() {
			name += "/"
		}
		fmt.Fprintf(&b, "%-6s %8d  %s\n", mode, size, name)
	}
	return b.String(), nil
}

func handleSearchFiles(root string, raw json.RawMessage) (string, error) {
	var a toolArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return "", fmt.Errorf("invalid args: %w", err)
	}
	var matches []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		matched, err := filepath.Match(a.Pattern, d.Name())
		if err != nil {
			return nil
		}
		if matched {
			rel, _ := filepath.Rel(root, path)
			matches = append(matches, rel)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("search: %w", err)
	}
	if len(matches) == 0 {
		return "No matches found", nil
	}
	return strings.Join(matches, "\n"), nil
}

func handleGetFileInfo(root string, raw json.RawMessage) (string, error) {
	var a toolArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return "", fmt.Errorf("invalid args: %w", err)
	}
	full, err := safeJoin(root, a.Path)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(full)
	if err != nil {
		return "", fmt.Errorf("stat %s: %w", a.Path, err)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Name:  %s\n", info.Name())
	fmt.Fprintf(&b, "Size:  %d bytes\n", info.Size())
	fmt.Fprintf(&b, "Mode:  %s\n", info.Mode().String())
	fmt.Fprintf(&b, "Mod:   %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
	if info.IsDir() {
		fmt.Fprintf(&b, "Type:  directory\n")
	} else {
		fmt.Fprintf(&b, "Type:  file\n")
	}
	return b.String(), nil
}

func handleCreateDir(root string, raw json.RawMessage) (string, error) {
	var a toolArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return "", fmt.Errorf("invalid args: %w", err)
	}
	full, err := safeJoin(root, a.Path)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(full, 0755); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", a.Path, err)
	}
	return fmt.Sprintf("Created directory %s", a.Path), nil
}
