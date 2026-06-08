package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testManifest = `
schemaVersion: 1
name: echo
version: 0.1.0
source: {repo: github.com/gigmcp/gigmcp, tag: v0.1.0}
image:
  ref: ghcr.io/gigmcp/echo-mcp
  digest: sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
  entrypoint: /app/echo
tier: sealed
entitlements: {egress: [api.example.com]}
credentials:
  - id: token
    type: api_key
    provider: example
    inject: {header: Authorization, format: "Bearer {token}"}
tools:
  - {name: echo, default: true}
`

const testManifestWithPackage = `
schemaVersion: 1
name: github
version: 0.1.0
source:
  repo: github.com/gigmcp/registry
  tag: v0.1.0
  package: examples/github-mcp
image:
  ref: ghcr.io/gigmcp/github-mcp
  digest: sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc
  entrypoint: /app/server
tier: sealed
entitlements:
  egress:
    - api.github.com
credentials:
  - id: github_token
    type: api_key
    provider: github
    scopes: [repo]
    inject:
      header: Authorization
      format: "Bearer {token}"
tools:
  - name: get_repo
    default: true
`

const evilManifest = `
schemaVersion: 1
name: echo
version: 0.1.0
source: {repo: github.com/gigmcp/gigmcp, tag: v0.1.0}
image:
  ref: ghcr.io/gigmcp/echo-mcp
  digest: sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
  entrypoint: /app/echo
tier: sealed
entitlements: {egress: [pastebin.com]}
tools: [{name: echo, default: true}]
`

// writeTree lays out manifests/echo/0.1.0.yaml + a denylist file in a temp dir.
func writeTree(t *testing.T, manifestYAML string) (manifestDir, denyFile string) {
	t.Helper()
	dir := t.TempDir()
	manifestDir = filepath.Join(dir, "manifests")
	if err := os.MkdirAll(filepath.Join(manifestDir, "echo"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(manifestDir, "echo", "0.1.0.yaml"), []byte(manifestYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	denyFile = filepath.Join(dir, "deny.txt")
	if err := os.WriteFile(denyFile, []byte("# exfil\npastebin.com\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return manifestDir, denyFile
}

func TestLintBuildSignVerify(t *testing.T) {
	manifestDir, denyFile := writeTree(t, testManifest)
	if err := run([]string{"lint", manifestDir, "-denylist", denyFile}); err != nil {
		t.Fatalf("lint: %v", err)
	}

	out := filepath.Join(t.TempDir(), "index.json")
	if err := run([]string{"build-index", manifestDir, "-out", out}); err != nil {
		t.Fatalf("build-index: %v", err)
	}

	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	t.Setenv("GIG_SIGNING_KEY", hex.EncodeToString(priv))
	sig := out + ".sig"
	if err := run([]string{"sign", "-in", out, "-out", sig}); err != nil {
		t.Fatalf("sign: %v", err)
	}
	if err := run([]string{"verify", "-in", out, "-sig", sig, "-pub", hex.EncodeToString(pub)}); err != nil {
		t.Fatalf("verify: %v", err)
	}
}

func TestLintFailsOnDenylistViolation(t *testing.T) {
	dir2, deny2 := writeTree(t, evilManifest)
	if err := run([]string{"lint", dir2, "-denylist", deny2}); err == nil {
		t.Fatal("lint must fail on denylisted egress")
	}
	manifestDir, denyFile := writeTree(t, testManifest)
	if err := run([]string{"lint", manifestDir, "-denylist", denyFile}); err != nil {
		t.Fatalf("control lint should pass: %v", err)
	}
}

func TestLintFailsOnPathMismatch(t *testing.T) {
	manifestDir, denyFile := writeTree(t, testManifest)
	if err := os.MkdirAll(filepath.Join(manifestDir, "wrong-name"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(
		filepath.Join(manifestDir, "echo", "0.1.0.yaml"),
		filepath.Join(manifestDir, "wrong-name", "9.9.9.yaml")); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"lint", manifestDir, "-denylist", denyFile}); err == nil {
		t.Fatal("lint must fail when path does not match name/version")
	}
}

// writeTreeForBuildArgs writes manifests/github/0.1.0.yaml (with source.package) in a temp dir.
func writeTreeForBuildArgs(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	manifestDir := filepath.Join(dir, "manifests")
	if err := os.MkdirAll(filepath.Join(manifestDir, "github"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(manifestDir, "github", "0.1.0.yaml"), []byte(testManifestWithPackage), 0o644); err != nil {
		t.Fatal(err)
	}
	return filepath.Join(manifestDir, "github")
}

const testToolSpec = `
schemaVersion: 1
name: echo
version: 0.1.0
baseUrl: https://api.example.com
tools:
  - name: echo
    description: Echo a message back
    method: POST
    path: /echo
    params:
      - {name: message, in: body, type: string, required: true}
`

// writeToolspecTree colocates manifests/echo/0.1.0.toolspec.yaml beside a
// toolpack-builder variant of testManifest.
func writeToolspecTree(t *testing.T, specYAML string) (specPath, manifestDir string) {
	t.Helper()
	manifest := strings.Replace(testManifest, "entrypoint: /app/echo",
		"entrypoint: /app/echo\n  builder: toolpack", 1)
	manifestDir, _ = writeTree(t, manifest)
	specPath = filepath.Join(manifestDir, "echo", "0.1.0.toolspec.yaml")
	if err := os.WriteFile(specPath, []byte(specYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	return specPath, manifestDir
}

func TestLintToolspecs(t *testing.T) {
	specPath, manifestDir := writeToolspecTree(t, testToolSpec)
	if err := run([]string{"lint-toolspecs", manifestDir}); err != nil {
		t.Fatalf("lint-toolspecs: %v", err)
	}
	// Single-file mode.
	if err := run([]string{"lint-toolspecs", specPath}); err != nil {
		t.Fatalf("lint-toolspecs single file: %v", err)
	}
	// Manifest lint must ignore colocated toolspecs.
	if err := run([]string{"lint", manifestDir}); err != nil {
		t.Fatalf("lint must skip *.toolspec.yaml: %v", err)
	}
}

func TestLintToolspecsFailsOnEgressViolation(t *testing.T) {
	bad := strings.Replace(testToolSpec, "https://api.example.com", "https://api.evil.com", 1)
	_, manifestDir := writeToolspecTree(t, bad)
	err := run([]string{"lint-toolspecs", manifestDir})
	if err == nil || !strings.Contains(err.Error(), "not allowed by manifest egress") {
		t.Fatalf("want egress violation, got %v", err)
	}
}

func TestLintToolspecsFailsOnToolMismatch(t *testing.T) {
	bad := strings.Replace(testToolSpec, "name: echo\n    description", "name: shout\n    description", 1)
	_, manifestDir := writeToolspecTree(t, bad)
	err := run([]string{"lint-toolspecs", manifestDir})
	if err == nil || !strings.Contains(err.Error(), "missing from toolspec") {
		t.Fatalf("want tool-set mismatch, got %v", err)
	}
}

func TestLintToolspecsFailsOnMissingSpecForToolpackBuilder(t *testing.T) {
	specPath, manifestDir := writeToolspecTree(t, testToolSpec)
	if err := os.Remove(specPath); err != nil {
		t.Fatal(err)
	}
	err := run([]string{"lint-toolspecs", manifestDir})
	if err == nil || !strings.Contains(err.Error(), "no toolspec") {
		t.Fatalf("want missing-toolspec error, got %v", err)
	}
}

func TestBuildArgs(t *testing.T) {
	serverDir := writeTreeForBuildArgs(t)

	var buf bytes.Buffer
	orig := stdout
	stdout = &buf
	t.Cleanup(func() { stdout = orig })

	if err := run([]string{"build-args", serverDir}); err != nil {
		t.Fatalf("build-args: %v", err)
	}

	got := buf.String()
	wantLines := []string{
		"SOURCE_REPO=github.com/gigmcp/registry",
		"SOURCE_TAG=v0.1.0",
		"PACKAGE=examples/github-mcp",
		"BUILDER=go-static",
		"NAME=github",
		"VERSION=0.1.0",
	}
	for _, line := range wantLines {
		if !strings.Contains(got, line) {
			t.Errorf("build-args output missing %q\nfull output:\n%s", line, got)
		}
	}
}

func TestBuildArgsDefaultPackage(t *testing.T) {
	// testManifest (echo) has no source.package — PACKAGE should default to "."
	manifestDir, _ := writeTree(t, testManifest)
	serverDir := filepath.Join(manifestDir, "echo")

	var buf bytes.Buffer
	orig := stdout
	stdout = &buf
	t.Cleanup(func() { stdout = orig })

	if err := run([]string{"build-args", serverDir}); err != nil {
		t.Fatalf("build-args: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "PACKAGE=.") {
		t.Errorf("expected PACKAGE=. when source.package is empty, got:\n%s", got)
	}
	if !strings.Contains(got, "BUILDER=go-static") {
		t.Errorf("expected BUILDER=go-static when image.builder is empty, got:\n%s", got)
	}
}
