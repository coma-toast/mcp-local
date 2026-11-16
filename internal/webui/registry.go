package webui

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const defaultMCPRegistryURL = "https://registry.modelcontextprotocol.io/v0.1/servers"

type mcpRegistryResponse struct {
	Servers []struct {
		Server mcpRegistryServer `json:"server"`
	} `json:"servers"`
}

type mcpRegistryServer struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	WebsiteURL  string `json:"websiteUrl"`
	Repository  struct {
		URL string `json:"url"`
	} `json:"repository"`
	Packages []struct {
		RegistryType string `json:"registryType"`
		Identifier   string `json:"identifier"`
		RuntimeHint  string `json:"runtimeHint"`
		Transport    struct {
			Type string `json:"type"`
		} `json:"transport"`
		EnvironmentVariables []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			IsRequired  bool   `json:"isRequired"`
		} `json:"environmentVariables"`
	} `json:"packages"`
	Remotes []struct {
		Type    string `json:"type"`
		URL     string `json:"url"`
		Headers []struct {
			Name        string `json:"name"`
			Value       string `json:"value"`
			Description string `json:"description"`
		} `json:"headers"`
	} `json:"remotes"`
}

var nonSlugChars = regexp.MustCompile(`[^a-z0-9-]+`)

func slugifyServiceName(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	if idx := strings.LastIndex(s, "/"); idx >= 0 && idx+1 < len(s) {
		s = s[idx+1:]
	}
	s = strings.TrimPrefix(s, "mcp-server-")
	s = strings.TrimPrefix(s, "server-")
	s = strings.TrimSuffix(s, "-mcp-server")
	s = strings.TrimSuffix(s, "-mcp")
	s = nonSlugChars.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "mcp-service"
	}
	return s
}

func (s *Server) handleRegistrySearch(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		writeJSON(w, http.StatusOK, map[string]interface{}{"presets": []PresetDTO{}})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	u, err := url.Parse(defaultMCPRegistryURL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	params := u.Query()
	params.Set("search", q)
	params.Set("version", "latest")
	params.Set("limit", "20")
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "mcp-local-webui/1.0")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		// Return empty presets with error message so UI degrades gracefully offline.
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"presets": []PresetDTO{},
			"error":   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"presets": []PresetDTO{},
			"error":   err.Error(),
		})
		return
	}

	var parsed mcpRegistryResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"presets": []PresetDTO{},
			"error":   fmt.Sprintf("registry parse error: %v", err),
		})
		return
	}

	presets := make([]PresetDTO, 0, len(parsed.Servers))
	seen := make(map[string]bool)
	for _, entry := range parsed.Servers {
		srv := entry.Server
		if srv.Name == "" || seen[srv.Name] {
			continue
		}
		seen[srv.Name] = true
		if preset, ok := registryServerToPreset(srv); ok {
			presets = append(presets, preset)
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"presets": presets})
}

func registryServerToPreset(srv mcpRegistryServer) (PresetDTO, bool) {
	svcName := slugifyServiceName(srv.Name)
	title := strings.TrimSpace(srv.Title)
	if title == "" {
		title = srv.Name
	}
	srcURL := strings.TrimSpace(srv.Repository.URL)
	if srcURL == "" {
		srcURL = strings.TrimSpace(srv.WebsiteURL)
	}

	req := ServiceUpsertRequest{
		Name:        svcName,
		Description: strings.TrimSpace(srv.Description),
		Args:        []string{},
		Env:         map[string]string{},
		Headers:     map[string]string{},
		Deps:        []string{},
		ExtraPorts:  []int{},
	}

	var hints []string
	configured := false

	// Prefer stdio package when available
	for _, pkg := range srv.Packages {
		id := strings.TrimSpace(pkg.Identifier)
		if id == "" {
			continue
		}
		rt := strings.ToLower(strings.TrimSpace(pkg.RegistryType))
		switch rt {
		case "npm":
			req.Type = "stdio"
			req.Command = "npx"
			req.Args = []string{"-y", id}
			configured = true
		case "pypi":
			req.Type = "stdio"
			req.Command = "uvx"
			req.Args = []string{id}
			configured = true
		case "oci":
			req.Type = "stdio"
			req.Command = "docker"
			req.Args = []string{"run", "-i", "--rm", id}
			configured = true
		}
		if configured {
			for _, ev := range pkg.EnvironmentVariables {
				k := strings.TrimSpace(ev.Name)
				if k == "" {
					continue
				}
				req.Env[k] = ""
				if ev.Description != "" {
					hints = append(hints, fmt.Sprintf("%s: %s", k, ev.Description))
				} else {
					hints = append(hints, k)
				}
			}
			break
		}
	}

	// Fallback to remote Streamable HTTP / SSE endpoint
	if !configured && len(srv.Remotes) > 0 {
		rem := srv.Remotes[0]
		if strings.TrimSpace(rem.URL) != "" {
			req.Type = "http"
			req.MCPURL = strings.TrimSpace(rem.URL)
			for _, h := range rem.Headers {
				name := strings.TrimSpace(h.Name)
				if name == "" {
					continue
				}
				val := strings.TrimSpace(h.Value)
				if val == "" && strings.EqualFold(name, "Authorization") {
					val = "Bearer "
				}
				req.Headers[name] = val
				if h.Description != "" {
					hints = append(hints, fmt.Sprintf("Header %s: %s", name, h.Description))
				}
			}
			configured = true
		}
	}

	if !configured {
		return PresetDTO{}, false
	}

	hintText := ""
	if len(hints) > 0 {
		hintText = "Configure: " + strings.Join(hints, " • ")
	}

	return PresetDTO{
		ID:          "registry:" + srv.Name,
		Title:       title,
		Category:    "MCP Registry",
		Source:      "github.com/mcp",
		SourceURL:   srcURL,
		Description: req.Description,
		Tags:        []string{svcName, srv.Name, req.Type},
		SetupHint:   hintText,
		Service:     req,
	}, true
}
