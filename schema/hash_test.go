package schema

import "testing"

// baseManifest returns a fully-populated, valid manifest for hash tests.
func baseManifest() *Manifest {
	return &Manifest{
		SchemaVersion: 1,
		Name:          "gmail",
		Version:       "0.1.0",
		Source:        Source{Repo: "github.com/gigmcp/toolpack", Tag: "v0.1.0"},
		Image: Image{
			Ref:        "ghcr.io/gigmcp/gmail-mcp",
			Digest:     "sha256:" + "ab12cd34" + "00000000000000000000000000000000000000000000000000000000",
			Entrypoint: "/app/server",
			Builder:    "toolpack",
		},
		Tier:         TierSealed,
		Entitlements: Entitlements{Egress: []string{"gmail.googleapis.com"}},
		Credentials: []Credential{{
			ID:       "google-oauth",
			Type:     "oauth2",
			Provider: "gmail",
			Scopes:   []string{"https://www.googleapis.com/auth/gmail.readonly"},
			Inject:   Inject{Header: "Authorization", Format: "Bearer {token}"},
		}},
		Tools: []Tool{{Name: "gmail_list_messages", Default: true}},
	}
}

func mustRuntimeHash(t *testing.T, m *Manifest) string {
	t.Helper()
	h, err := m.RuntimeHash()
	if err != nil {
		t.Fatalf("RuntimeHash: %v", err)
	}
	return h
}

// The gateway's cross-check scenario: an image baked BEFORE vendor/branding
// existed (no presentation fields) must RuntimeHash-equal the same manifest
// after vendor + branding are backfilled into the signed index. This is the
// whole point of Option B — backfilling presentation never forces a rebuild.
func TestRuntimeHash_IgnoresPresentationAndDigest(t *testing.T) {
	baked := baseManifest() // as the old image baked it: no vendor, no branding
	index := baseManifest()
	// what the re-signed index will carry after backfill:
	index.Credentials[0].Vendor = "google"
	index.DisplayName = "Gmail"
	index.Description = "Send and read Gmail messages"
	index.Category = "comms"
	index.Icon = "icons/gmail.svg"
	// and a different (e.g. multi-arch re-pinned) image digest:
	index.Image.Digest = "sha256:" + "ffffffff" + "00000000000000000000000000000000000000000000000000000000"

	if got, want := mustRuntimeHash(t, index), mustRuntimeHash(t, baked); got != want {
		t.Fatalf("RuntimeHash diverged on presentation/digest-only change:\n baked=%s\n index=%s", want, got)
	}
}

// Every runtime/security field must move the hash — otherwise a tampered image
// could change it and slip past the cross-check.
func TestRuntimeHash_SensitiveToRuntimeFields(t *testing.T) {
	base := mustRuntimeHash(t, baseManifest())
	mutate := map[string]func(*Manifest){
		"name":             func(m *Manifest) { m.Name = "other" },
		"version":          func(m *Manifest) { m.Version = "0.2.0" },
		"schemaVersion":    func(m *Manifest) { m.SchemaVersion = 2 },
		"source.repo":      func(m *Manifest) { m.Source.Repo = "github.com/evil/x" },
		"source.tag":       func(m *Manifest) { m.Source.Tag = "v9.9.9" },
		"image.ref":        func(m *Manifest) { m.Image.Ref = "ghcr.io/evil/x" },
		"image.entrypoint": func(m *Manifest) { m.Image.Entrypoint = "/app/evil" },
		"image.builder":    func(m *Manifest) { m.Image.Builder = "go-static" },
		"tier":             func(m *Manifest) { m.Tier = TierEntrusted },
		"egress":           func(m *Manifest) { m.Entitlements.Egress = []string{"evil.example.com"} },
		"cred.id":          func(m *Manifest) { m.Credentials[0].ID = "x" },
		"cred.type":        func(m *Manifest) { m.Credentials[0].Type = "api_key" },
		"cred.provider":    func(m *Manifest) { m.Credentials[0].Provider = "evil" },
		"cred.scopes":      func(m *Manifest) { m.Credentials[0].Scopes = []string{"admin"} },
		"cred.inject":      func(m *Manifest) { m.Credentials[0].Inject.Header = "X-Evil" },
		"tool.name":        func(m *Manifest) { m.Tools[0].Name = "evil_tool" },
		"tool.default":     func(m *Manifest) { m.Tools[0].Default = false },
	}
	for field, mut := range mutate {
		t.Run(field, func(t *testing.T) {
			m := baseManifest()
			mut(m)
			if got := mustRuntimeHash(t, m); got == base {
				t.Fatalf("RuntimeHash unchanged after mutating runtime field %q — security gap", field)
			}
		})
	}
}

// Presentation/digest fields must NOT move the hash, individually.
func TestRuntimeHash_InsensitiveToPresentationFields(t *testing.T) {
	base := mustRuntimeHash(t, baseManifest())
	mutate := map[string]func(*Manifest){
		"image.digest": func(m *Manifest) {
			m.Image.Digest = "sha256:" + "deadbeef" + "00000000000000000000000000000000000000000000000000000000"
		},
		"cred.vendor": func(m *Manifest) { m.Credentials[0].Vendor = "google" },
		"displayName": func(m *Manifest) { m.DisplayName = "Gmail" },
		"description": func(m *Manifest) { m.Description = "x" },
		"category":    func(m *Manifest) { m.Category = "comms" },
		"icon":        func(m *Manifest) { m.Icon = "icons/gmail.svg" },
	}
	for field, mut := range mutate {
		t.Run(field, func(t *testing.T) {
			m := baseManifest()
			mut(m)
			if got := mustRuntimeHash(t, m); got != base {
				t.Fatalf("RuntimeHash changed after mutating presentation field %q — should be excluded", field)
			}
		})
	}
}

// RuntimeHash must not mutate its receiver (it clears fields on a copy).
func TestRuntimeHash_DoesNotMutateReceiver(t *testing.T) {
	m := baseManifest()
	m.Credentials[0].Vendor = "google"
	m.DisplayName = "Gmail"
	digestBefore := m.Image.Digest
	if _, err := m.RuntimeHash(); err != nil {
		t.Fatalf("RuntimeHash: %v", err)
	}
	if m.Image.Digest != digestBefore {
		t.Fatalf("RuntimeHash cleared receiver image.digest: %q", m.Image.Digest)
	}
	if m.Credentials[0].Vendor != "google" {
		t.Fatalf("RuntimeHash cleared receiver credential vendor: %q", m.Credentials[0].Vendor)
	}
	if m.DisplayName != "Gmail" {
		t.Fatalf("RuntimeHash cleared receiver displayName: %q", m.DisplayName)
	}
}

// Hash (full) MUST still move when presentation changes — it's the broad
// equality/identity hash; RuntimeHash is the narrowed security comparator.
// This documents the deliberate difference between the two.
func TestHash_StillSensitiveToPresentation(t *testing.T) {
	a, err := baseManifest().Hash()
	if err != nil {
		t.Fatal(err)
	}
	m := baseManifest()
	m.Credentials[0].Vendor = "google"
	b, err := m.Hash()
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Fatal("full Hash() should differ when vendor changes (it hashes everything)")
	}
}
