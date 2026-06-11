package schema

import (
	"strings"
	"testing"
)

func good(t *testing.T) *Manifest {
	t.Helper()
	m, err := Parse([]byte(goodYAML))
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestValidateGood(t *testing.T) {
	if err := good(t).Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestValidateRejects(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*Manifest)
		want string
	}{
		{"bad schemaVersion", func(m *Manifest) { m.SchemaVersion = 2 }, "schemaVersion"},
		{"underscore name", func(m *Manifest) { m.Name = "slack_mcp" }, "name"},
		{"uppercase name", func(m *Manifest) { m.Name = "Slack" }, "name"},
		{"bad version", func(m *Manifest) { m.Version = "1.4" }, "version"},
		{"leading-zero version", func(m *Manifest) { m.Version = "01.4.2" }, "version"},
		{"missing source", func(m *Manifest) { m.Source.Repo = "" }, "source"},
		{"bad digest", func(m *Manifest) { m.Image.Digest = "sha256:short" }, "digest"},
		{"missing ref", func(m *Manifest) { m.Image.Ref = "" }, "image.ref"},
		{"relative entrypoint", func(m *Manifest) { m.Image.Entrypoint = "server" }, "entrypoint"},
		{"bad tier", func(m *Manifest) { m.Tier = "yolo" }, "tier"},
		{"sealed without header", func(m *Manifest) { m.Credentials[0].Inject = Inject{Env: "TOKEN"} }, "sealed"},
		{"format without token", func(m *Manifest) { m.Credentials[0].Inject.Format = "Bearer xyz" }, "{token}"},
		{"dup credential id", func(m *Manifest) { m.Credentials = append(m.Credentials, m.Credentials[0]) }, "credential"},
		{"dup tool", func(m *Manifest) { m.Tools = append(m.Tools, m.Tools[0]) }, "tool"},
		{"empty tool name", func(m *Manifest) { m.Tools[0].Name = "" }, "tool"},
		{"dotdot package", func(m *Manifest) { m.Source.Package = "../evil" }, "package"},
		{"absolute package", func(m *Manifest) { m.Source.Package = "/etc" }, "package"},
		{"bad builder", func(m *Manifest) { m.Image.Builder = "ruby" }, "builder"},
		{"oauth2 missing vendor", func(m *Manifest) { m.Credentials[0].Type = "oauth2" }, "vendor is required"},
		{"bad vendor slug", func(m *Manifest) { m.Credentials[0].Vendor = "Google_Inc" }, "vendor"},
		{"uppercase vendor", func(m *Manifest) { m.Credentials[0].Vendor = "Google" }, "vendor"},
		{"bad category", func(m *Manifest) { m.Category = "news" }, "category"},
		{"foreign icon path", func(m *Manifest) { m.Icon = "icons/evil.svg" }, "icon"},
		{"icon not under icons/", func(m *Manifest) { m.Icon = "evil.svg" }, "icon"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := good(t)
			tc.mut(m)
			err := m.Validate()
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("want error containing %q, got %v", tc.want, err)
			}
		})
	}
}

func TestValidateOAuth2Vendor(t *testing.T) {
	m := good(t)
	m.Credentials[0].Type = "oauth2"
	m.Credentials[0].Vendor = "google"
	if err := m.Validate(); err != nil {
		t.Fatalf("oauth2 credential with a valid vendor should validate: %v", err)
	}
	// non-oauth2 credentials may omit vendor (the api_key fixture already does).
	m2 := good(t)
	if err := m2.Validate(); err != nil {
		t.Fatalf("api_key credential without vendor should validate: %v", err)
	}
}

func TestValidateGoodIcon(t *testing.T) {
	m := good(t) // name: slack-mcp
	m.Icon = "icons/slack-mcp.svg"
	if err := m.Validate(); err != nil {
		t.Fatalf("icon icons/<name>.svg should validate: %v", err)
	}
}

func TestValidateGoodCategories(t *testing.T) {
	for _, c := range []string{"", "ai", "analytics", "comms", "CRM", "data", "design", "dev-tools", "documents", "e-commerce", "finance", "HR/ATS", "marketing", "productivity", "search", "social", "storage", "support"} {
		m := good(t)
		m.Category = c
		if err := m.Validate(); err != nil {
			t.Fatalf("category %q should validate, got: %v", c, err)
		}
	}
}

func TestValidateGoodPackage(t *testing.T) {
	m := good(t)
	m.Source.Package = "examples/github-mcp"
	if err := m.Validate(); err != nil {
		t.Fatalf("manifest with valid source.package should validate: %v", err)
	}
}

func TestValidateGoodBuilder(t *testing.T) {
	for _, b := range []string{"", "go-static", "node", "python"} {
		m := good(t)
		m.Image.Builder = b
		if err := m.Validate(); err != nil {
			t.Fatalf("Builder %q should validate, got: %v", b, err)
		}
	}
}

func TestValidateEntrustedInjection(t *testing.T) {
	m := good(t)
	m.Tier = TierEntrusted
	m.Credentials[0].Inject = Inject{Env: "SLACK_TOKEN"}
	if err := m.Validate(); err != nil {
		t.Fatalf("entrusted env injection should validate: %v", err)
	}
	m.Credentials[0].Inject.Header = "Authorization"
	if err := m.Validate(); err == nil {
		t.Fatal("entrusted with header should fail (exactly one mode)")
	}
}

func TestCheckEgressEntry(t *testing.T) {
	ok := []string{"slack.com", "*.slack.com", "api.github.com", "*.s3.us-east-1.amazonaws.com"}
	for _, e := range ok {
		if err := CheckEgressEntry(e); err != nil {
			t.Errorf("CheckEgressEntry(%q) = %v, want nil", e, err)
		}
	}
	bad := map[string]string{
		"*":             "wildcard",
		"*.com":         "two labels",
		"foo.*.com":     "wildcard",
		"slack.com:443": "ports",
		"slack.com/path": "ports",
		"10.0.0.1":      "IP",
		"*.10.0.0.1":    "IP",
		"":              "label",
		"UPPER.com":     "label",
		"-bad.com":      "label",
	}
	for e, want := range bad {
		err := CheckEgressEntry(e)
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("CheckEgressEntry(%q) = %v, want error containing %q", e, err, want)
		}
	}
}

func TestLintDenylist(t *testing.T) {
	deny := []string{"# exfil domains", "pastebin.com", "webhook.site", "ngrok.io"}
	m := good(t)
	if err := m.Lint(deny); err != nil {
		t.Fatalf("clean manifest should lint: %v", err)
	}
	for _, evil := range []string{"pastebin.com", "api.pastebin.com", "*.ngrok.io"} {
		m := good(t)
		m.Entitlements.Egress = []string{evil}
		if err := m.Lint(deny); err == nil || !strings.Contains(err.Error(), "denylist") {
			t.Errorf("Lint with egress %q = %v, want denylist error", evil, err)
		}
	}
}
