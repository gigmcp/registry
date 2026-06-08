package schema

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	sigsyaml "sigs.k8s.io/yaml"
)

// ToolSpec is the declarative tool→HTTP mapping consumed by the generic
// toolpack engine (github.com/gigmcp/toolpack). One spec pairs 1:1 with a
// manifest version: manifests/<name>/<version>.toolspec.yaml next to
// manifests/<name>/<version>.yaml. The spec is data, not code — consistent
// with the aggregator policy (manifests and build recipes only).
type ToolSpec struct {
	SchemaVersion int    `json:"schemaVersion"`
	Name          string `json:"name"`
	Version       string `json:"version"`
	// BaseURL is "https://<host>" only — no port, path, or userinfo. The host
	// must be allowed by the paired manifest's egress allowlist.
	BaseURL string `json:"baseUrl"`
	// Auth describes how the credential enters requests. Required iff the
	// paired manifest is entrusted-tier (its inject is env-only, so the spec
	// must say where the env secret goes); forbidden for sealed tier, where
	// the manifest credential's inject is authoritative and the egress proxy
	// performs the injection.
	Auth  *SpecAuth  `json:"auth,omitempty"`
	Tools []SpecTool `json:"tools"`
}

type SpecAuth struct {
	Header string `json:"header"`
	Format string `json:"format"` // must contain "{token}"
}

type SpecTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Method      string `json:"method"` // GET | POST | PUT | PATCH | DELETE
	// Path is absolute; "{placeholder}" segments bind to in:path params.
	Path string `json:"path"`
	// BaseURL optionally overrides the spec-level base URL for services with
	// multiple API hosts; same syntax and egress rules.
	BaseURL string `json:"baseUrl,omitempty"`
	// Encoding selects the body encoding for in:body params: "json" (default)
	// or "form" (application/x-www-form-urlencoded).
	Encoding string      `json:"encoding,omitempty"`
	Params   []SpecParam `json:"params,omitempty"`
}

type SpecParam struct {
	Name        string `json:"name"`
	In          string `json:"in"`   // path | query | body | header
	Type        string `json:"type"` // string | integer | number | boolean | object | array
	Required    bool   `json:"required,omitempty"`
	Description string `json:"description,omitempty"`
}

var pathPlaceholderRE = regexp.MustCompile(`\{([^{}/]+)\}`)

// ParseToolSpec decodes toolspec YAML. Strict: unknown fields are errors, so
// typos (e.g. "metod:") cannot silently produce a broken tool.
func ParseToolSpec(data []byte) (*ToolSpec, error) {
	var s ToolSpec
	if err := sigsyaml.UnmarshalStrict(data, &s); err != nil {
		return nil, fmt.Errorf("schema: parse toolspec: %w", err)
	}
	return &s, nil
}

// Validate checks structural rules that need no manifest. Pairing rules
// (tool-set equality, egress, tier/auth coherence) live in CheckAgainstManifest.
func (s *ToolSpec) Validate() error {
	var errs []error
	add := func(format string, a ...any) { errs = append(errs, fmt.Errorf(format, a...)) }

	if s.SchemaVersion != 1 {
		add("schemaVersion must be 1, got %d", s.SchemaVersion)
	}
	if !nameRE.MatchString(s.Name) {
		add("name %q invalid: lowercase [a-z0-9-], no underscores", s.Name)
	}
	if !versionRE.MatchString(s.Version) {
		add("version %q invalid: must be MAJOR.MINOR.PATCH", s.Version)
	}
	if err := checkBaseURL(s.BaseURL); err != nil {
		add("baseUrl: %v", err)
	}
	if s.Auth != nil {
		if s.Auth.Header == "" || !strings.Contains(s.Auth.Format, "{token}") {
			add("auth requires header and a format containing {token}")
		}
	}
	if len(s.Tools) == 0 {
		add("tools must not be empty")
	}
	seen := map[string]bool{}
	for _, t := range s.Tools {
		if t.Name == "" || seen[t.Name] {
			add("tool name %q empty or duplicate", t.Name)
		}
		seen[t.Name] = true
		if err := t.validate(s.Auth); err != nil {
			add("tool %q: %v", t.Name, err)
		}
	}
	return errors.Join(errs...)
}

func (t *SpecTool) validate(auth *SpecAuth) error {
	var errs []error
	add := func(format string, a ...any) { errs = append(errs, fmt.Errorf(format, a...)) }

	switch t.Method {
	case "GET", "POST", "PUT", "PATCH", "DELETE":
	default:
		add("method %q invalid: GET|POST|PUT|PATCH|DELETE", t.Method)
	}
	if !strings.HasPrefix(t.Path, "/") {
		add("path %q must be absolute", t.Path)
	}
	if t.BaseURL != "" {
		if err := checkBaseURL(t.BaseURL); err != nil {
			add("baseUrl: %v", err)
		}
	}
	switch t.Encoding {
	case "", "json", "form":
	default:
		add("encoding %q invalid: json|form (or omitted for json)", t.Encoding)
	}
	if t.Description == "" {
		add("description is required (it becomes the MCP tool description)")
	}

	bodyAllowed := t.Method == "POST" || t.Method == "PUT" || t.Method == "PATCH"
	pathParams := map[string]bool{}
	seen := map[string]bool{}
	for _, p := range t.Params {
		// Params share one flat namespace: they become the properties of the
		// MCP tool's input schema, regardless of where each one is sent.
		if p.Name == "" || seen[p.Name] {
			add("param %q empty or duplicate", p.Name)
		}
		seen[p.Name] = true
		switch p.In {
		case "path":
			pathParams[p.Name] = true
			if !p.Required {
				add("path param %q must be required", p.Name)
			}
		case "query", "header":
		case "body":
			if !bodyAllowed {
				add("body param %q not allowed on %s", p.Name, t.Method)
			}
		default:
			add("param %q: in %q invalid: path|query|body|header", p.Name, p.In)
		}
		switch p.Type {
		case "string", "integer", "number", "boolean", "object", "array":
		default:
			add("param %q: type %q invalid: string|integer|number|boolean|object|array", p.Name, p.Type)
		}
		if p.In == "header" && auth != nil && strings.EqualFold(p.Name, auth.Header) {
			add("header param %q collides with the auth header", p.Name)
		}
	}
	// Every {placeholder} in the path needs an in:path param, and vice versa.
	holes := map[string]bool{}
	for _, m := range pathPlaceholderRE.FindAllStringSubmatch(t.Path, -1) {
		holes[m[1]] = true
		if !pathParams[m[1]] {
			add("path placeholder {%s} has no in:path param", m[1])
		}
	}
	for name := range pathParams {
		if !holes[name] {
			add("in:path param %q does not appear in path %q", name, t.Path)
		}
	}
	return errors.Join(errs...)
}

// checkBaseURL enforces "https://<host>" with nothing else: the engine joins
// tool paths onto it, and the host must be matchable against egress entries.
func checkBaseURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if u.Scheme != "https" {
		return fmt.Errorf("%q: scheme must be https", raw)
	}
	if u.Host == "" || u.Host != u.Hostname() {
		return fmt.Errorf("%q: host only — no port", raw)
	}
	if u.Path != "" || u.RawQuery != "" || u.Fragment != "" || u.User != nil {
		return fmt.Errorf("%q: must be https://<host> with no path, query, or userinfo", raw)
	}
	return CheckEgressEntry(u.Hostname())
}

// egressAllows mirrors the proxy's allowed() semantics EXACTLY
// (internal/proxy in gigmcp): case-insensitive exact host, or "*.suffix"
// matching only true subdomains — the bare suffix itself does NOT match.
func egressAllows(entries []string, host string) bool {
	host = strings.ToLower(host)
	for _, e := range entries {
		e = strings.ToLower(e)
		if host == e {
			return true
		}
		if strings.HasPrefix(e, "*.") && strings.HasSuffix(host, e[1:]) {
			return true
		}
	}
	return false
}

// CheckAgainstManifest enforces the pairing rules between a toolspec and its
// manifest: identity, tool-set equality, egress coverage of every base URL,
// and tier/auth coherence. The manifest is assumed already Validate()d.
func (s *ToolSpec) CheckAgainstManifest(m *Manifest) error {
	var errs []error
	add := func(format string, a ...any) { errs = append(errs, fmt.Errorf(format, a...)) }

	if s.Name != m.Name || s.Version != m.Version {
		add("toolspec %s@%s does not match manifest %s@%s", s.Name, s.Version, m.Name, m.Version)
	}
	specTools := map[string]bool{}
	for _, t := range s.Tools {
		specTools[t.Name] = true
	}
	for _, t := range m.Tools {
		if !specTools[t.Name] {
			add("manifest tool %q missing from toolspec", t.Name)
		}
		delete(specTools, t.Name)
	}
	for name := range specTools {
		add("toolspec tool %q not declared in manifest", name)
	}
	checkHost := func(raw string) {
		u, err := url.Parse(raw)
		if err != nil {
			return // Validate already rejected it
		}
		if !egressAllows(m.Entitlements.Egress, u.Hostname()) {
			add("base URL host %q not allowed by manifest egress %v", u.Hostname(), m.Entitlements.Egress)
		}
	}
	checkHost(s.BaseURL)
	for _, t := range s.Tools {
		if t.BaseURL != "" {
			checkHost(t.BaseURL)
		}
	}
	switch m.Tier {
	case TierEntrusted:
		if len(m.Credentials) > 0 && s.Auth == nil {
			add("entrusted tier requires a toolspec auth block (manifest inject is env-only)")
		}
	case TierSealed:
		if s.Auth != nil {
			add("sealed tier forbids a toolspec auth block (manifest credential inject is authoritative)")
		}
	}
	return errors.Join(errs...)
}
