package schema

import (
	"strings"
	"testing"
)

const goodSpecYAML = `
schemaVersion: 1
name: slack-mcp
version: 1.4.2
baseUrl: https://slack.com
tools:
  - name: send_message
    description: Post a message to a channel
    method: POST
    path: /api/chat.postMessage
    params:
      - {name: channel, in: body, type: string, required: true}
      - {name: text, in: body, type: string, required: true}
  - name: admin_set_workspace_settings
    description: Update workspace settings
    method: POST
    path: /api/admin.workspaces.settings.set/{workspace}
    baseUrl: https://api.slack.com
    params:
      - {name: workspace, in: path, type: string, required: true}
      - {name: settings, in: body, type: object, required: true}
`

func mustSpec(t *testing.T, yaml string) *ToolSpec {
	t.Helper()
	s, err := ParseToolSpec([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseToolSpec: %v", err)
	}
	return s
}

func TestToolSpecGood(t *testing.T) {
	s := mustSpec(t, goodSpecYAML)
	if err := s.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	m, err := Parse([]byte(goodYAML))
	if err != nil {
		t.Fatalf("Parse manifest: %v", err)
	}
	if err := s.CheckAgainstManifest(m); err != nil {
		t.Fatalf("CheckAgainstManifest: %v", err)
	}
}

func TestToolSpecRejectsUnknownFields(t *testing.T) {
	if _, err := ParseToolSpec([]byte(goodSpecYAML + "\nsurprise: field\n")); err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestToolSpecValidateRejects(t *testing.T) {
	cases := []struct {
		name string
		mut  func(string) string
		want string
	}{
		{"http baseUrl", func(y string) string {
			return strings.Replace(y, "https://slack.com", "http://slack.com", 1)
		}, "scheme must be https"},
		{"baseUrl with path", func(y string) string {
			return strings.Replace(y, "https://slack.com", "https://slack.com/api", 1)
		}, "no path"},
		{"baseUrl with port", func(y string) string {
			return strings.Replace(y, "https://slack.com", "https://slack.com:8443", 1)
		}, "no port"},
		{"bad method", func(y string) string {
			return strings.Replace(y, "method: POST", "method: YEET", 1)
		}, "method"},
		{"body on GET", func(y string) string {
			return strings.Replace(y, "method: POST", "method: GET", 1)
		}, "not allowed on GET"},
		{"relative path", func(y string) string {
			return strings.Replace(y, "path: /api/chat.postMessage", "path: api/chat.postMessage", 1)
		}, "absolute"},
		{"missing description", func(y string) string {
			return strings.Replace(y, "    description: Post a message to a channel\n", "", 1)
		}, "description"},
		{"orphan placeholder", func(y string) string {
			return strings.Replace(y, "/api/chat.postMessage", "/api/{team}/chat.postMessage", 1)
		}, "no in:path param"},
		{"unused path param", func(y string) string {
			return strings.Replace(y, "path: /api/admin.workspaces.settings.set/{workspace}",
				"path: /api/admin.workspaces.settings.set", 1)
		}, "does not appear in path"},
		{"optional path param", func(y string) string {
			return strings.Replace(y, "{name: workspace, in: path, type: string, required: true}",
				"{name: workspace, in: path, type: string, required: false}", 1)
		}, "must be required"},
		{"bad param in", func(y string) string {
			return strings.Replace(y, "in: query", "in: cookie", 1)
		}, "invalid"},
		{"bad type", func(y string) string {
			return strings.Replace(y, "type: object", "type: blob", 1)
		}, "invalid"},
		{"duplicate tool", func(y string) string {
			return strings.Replace(y, "name: admin_set_workspace_settings", "name: send_message", 1)
		}, "duplicate"},
		{"auth without token", func(y string) string {
			return strings.Replace(y, "baseUrl: https://slack.com\n",
				"baseUrl: https://slack.com\nauth: {header: Authorization, format: Bearer xyz}\n", 1)
		}, "{token}"},
	}
	// goodSpecYAML has no in:query param; give "bad param in" one to mutate.
	withQuery := strings.Replace(goodSpecYAML,
		"{name: text, in: body, type: string, required: true}",
		"{name: text, in: body, type: string, required: true}\n      - {name: thread, in: query, type: string}", 1)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := goodSpecYAML
			if tc.name == "bad param in" {
				src = withQuery
			}
			s := mustSpec(t, tc.mut(src))
			err := s.Validate()
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("want error containing %q, got %v", tc.want, err)
			}
		})
	}
}

func TestCheckAgainstManifestRejects(t *testing.T) {
	m, err := Parse([]byte(goodYAML))
	if err != nil {
		t.Fatalf("Parse manifest: %v", err)
	}

	// Tool-set mismatch: spec missing a manifest tool.
	missing := mustSpec(t, goodSpecYAML)
	missing.Tools = missing.Tools[:1]
	if err := missing.CheckAgainstManifest(m); err == nil || !strings.Contains(err.Error(), "missing from toolspec") {
		t.Errorf("missing tool: got %v", err)
	}

	// Tool-set mismatch: spec declares an extra tool.
	extra := mustSpec(t, goodSpecYAML)
	extra.Tools = append(extra.Tools, SpecTool{Name: "rogue", Description: "x", Method: "GET", Path: "/x"})
	if err := extra.CheckAgainstManifest(m); err == nil || !strings.Contains(err.Error(), "not declared in manifest") {
		t.Errorf("extra tool: got %v", err)
	}

	// Egress: host not allowlisted.
	badHost := mustSpec(t, strings.Replace(goodSpecYAML, "https://api.slack.com", "https://api.evil.com", 1))
	if err := badHost.CheckAgainstManifest(m); err == nil || !strings.Contains(err.Error(), "not allowed by manifest egress") {
		t.Errorf("bad host: got %v", err)
	}

	// Wildcard egress allows subdomains: api.slack.com matches *.slack.com.
	wild := mustSpec(t, goodSpecYAML)
	if err := wild.CheckAgainstManifest(m); err != nil {
		t.Errorf("wildcard host should pass: %v", err)
	}

	// Sealed tier forbids a spec auth block.
	sealedAuth := mustSpec(t, strings.Replace(goodSpecYAML, "baseUrl: https://slack.com\n",
		"baseUrl: https://slack.com\nauth: {header: X-Auth, format: \"Bearer {token}\"}\n", 1))
	if err := sealedAuth.CheckAgainstManifest(m); err == nil || !strings.Contains(err.Error(), "sealed tier forbids") {
		t.Errorf("sealed auth: got %v", err)
	}

	// Entrusted tier requires a spec auth block.
	entrustedYAML := strings.ReplaceAll(goodYAML, "tier: sealed", "tier: entrusted")
	entrustedYAML = strings.Replace(entrustedYAML,
		"inject:\n      header: Authorization\n      format: \"Bearer {token}\"",
		"inject:\n      env: SLACK_TOKEN", 1)
	em, err := Parse([]byte(entrustedYAML))
	if err != nil {
		t.Fatalf("Parse entrusted manifest: %v", err)
	}
	noAuth := mustSpec(t, goodSpecYAML)
	if err := noAuth.CheckAgainstManifest(em); err == nil || !strings.Contains(err.Error(), "entrusted tier requires") {
		t.Errorf("entrusted no-auth: got %v", err)
	}

	// Identity mismatch.
	other := mustSpec(t, strings.Replace(goodSpecYAML, "name: slack-mcp", "name: other", 1))
	if err := other.CheckAgainstManifest(m); err == nil || !strings.Contains(err.Error(), "does not match manifest") {
		t.Errorf("identity: got %v", err)
	}
}

func TestEgressAllows(t *testing.T) {
	entries := []string{"api.example.com", "*.slack.com"}
	for host, want := range map[string]bool{
		"api.example.com":  true,
		"API.Example.com":  true,
		"sub.example.com":  false,
		"slack.com":        false, // *.suffix matches only true subdomains (proxy semantics)
		"a.b.slack.com":    true,
		"evilslack.com":    false,
		"slack.com.evil.x": false,
	} {
		if got := egressAllows(entries, host); got != want {
			t.Errorf("egressAllows(%q) = %v, want %v", host, got, want)
		}
	}
}
