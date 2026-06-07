package schema

import (
	"strings"
	"testing"
)

const goodYAML = `
schemaVersion: 1
name: slack-mcp
version: 1.4.2
source:
  repo: github.com/author/slack-mcp
  tag: v1.4.2
image:
  ref: ghcr.io/gigmcp/slack-mcp
  digest: sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
  entrypoint: /server
tier: sealed
entitlements:
  egress:
    - slack.com
    - "*.slack.com"
credentials:
  - id: slack_bot_token
    type: api_key
    provider: slack
    scopes: [chat:write]
    inject:
      header: Authorization
      format: "Bearer {token}"
tools:
  - name: send_message
    default: true
  - name: admin_set_workspace_settings
    default: false
`

func TestParseGoodManifest(t *testing.T) {
	m, err := Parse([]byte(goodYAML))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if m.Name != "slack-mcp" || m.Version != "1.4.2" || m.Tier != "sealed" {
		t.Fatalf("bad fields: %+v", m)
	}
	if m.Image.Entrypoint != "/server" || len(m.Entitlements.Egress) != 2 {
		t.Fatalf("bad image/egress: %+v", m)
	}
	if len(m.Credentials) != 1 || m.Credentials[0].Inject.Header != "Authorization" {
		t.Fatalf("bad credentials: %+v", m.Credentials)
	}
	if len(m.Tools) != 2 || !m.Tools[0].Default || m.Tools[1].Default {
		t.Fatalf("bad tools: %+v", m.Tools)
	}
}

func TestParseRejectsUnknownFields(t *testing.T) {
	bad := goodYAML + "\nsurprise: field\n"
	if _, err := Parse([]byte(bad)); err == nil {
		t.Fatal("expected error for unknown top-level field")
	}
	// Also verify transitive strictness: unknown field nested inside source.
	nestedBad := strings.ReplaceAll(goodYAML,
		"source:\n  repo:", "source:\n  surprise: nested\n  repo:")
	if _, err := Parse([]byte(nestedBad)); err == nil {
		t.Fatal("expected error for unknown nested field under source")
	}
}

func TestParseRejectsGarbage(t *testing.T) {
	if _, err := Parse([]byte("{[not yaml")); err == nil || !strings.Contains(err.Error(), "parse manifest") {
		t.Fatalf("expected parse error, got %v", err)
	}
}
