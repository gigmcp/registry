package schema

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"
)

func TestHashStableAndSensitive(t *testing.T) {
	m1, m2 := good(t), good(t)
	h1, err := m1.Hash()
	if err != nil || len(h1) != 64 {
		t.Fatalf("Hash: %v %q", err, h1)
	}
	h2, _ := m2.Hash()
	if h1 != h2 {
		t.Fatal("same manifest must hash identically")
	}
	m2.Entitlements.Egress = append(m2.Entitlements.Egress, "evil.example")
	if h3, _ := m2.Hash(); h3 == h1 {
		t.Fatal("changed manifest must change hash")
	}
}

func TestBuildIndexAndResolve(t *testing.T) {
	a, b := good(t), good(t)
	b.Version = "1.5.0"
	b.Image.Digest = "sha256:" + strings.Repeat("b", 64)
	ix, err := BuildIndex([]*Manifest{a, b}, "2026-06-06T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if ix.Servers["slack-mcp"].Latest != "1.5.0" {
		t.Fatalf("latest = %q, want 1.5.0", ix.Servers["slack-mcp"].Latest)
	}
	for ref, wantVer := range map[string]string{
		"slack-mcp":       "1.5.0",
		"slack-mcp@1.4.2": "1.4.2",
		b.Image.Digest:    "1.5.0",
	} {
		m, err := ix.Resolve(ref)
		if err != nil || m.Version != wantVer {
			t.Errorf("Resolve(%q) = %v, %v; want version %s", ref, m, err, wantVer)
		}
	}
	for _, ref := range []string{"nope", "slack-mcp@9.9.9", "slack-mcp@", "sha256:" + strings.Repeat("c", 64)} {
		if _, err := ix.Resolve(ref); err == nil {
			t.Errorf("Resolve(%q) should fail", ref)
		}
	}
	if _, err := BuildIndex([]*Manifest{a, a}, "x"); err == nil {
		t.Fatal("duplicate name@version must fail")
	}
}

func TestSignVerifyRoundTrip(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	ix, _ := BuildIndex([]*Manifest{good(t)}, "2026-06-06T00:00:00Z")
	raw, err := json.Marshal(ix)
	if err != nil {
		t.Fatal(err)
	}
	sig, err := Sign(hex.EncodeToString(priv), raw)
	if err != nil {
		t.Fatal(err)
	}
	if err := Verify(hex.EncodeToString(pub), raw, sig); err != nil {
		t.Fatalf("Verify: %v", err)
	}
	tampered := append([]byte(nil), raw...)
	tampered[20] ^= 1
	if err := Verify(hex.EncodeToString(pub), tampered, sig); err == nil {
		t.Fatal("tampered bytes must fail verification")
	}
	pub2, _, _ := ed25519.GenerateKey(rand.Reader)
	if err := Verify(hex.EncodeToString(pub2), raw, sig); err == nil {
		t.Fatal("verify accepted wrong key")
	}
	if _, err := Sign("zz", raw); err == nil {
		t.Fatal("bad key hex must fail")
	}
}
