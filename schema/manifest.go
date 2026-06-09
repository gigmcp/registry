// Package schema is the authoritative parser/validator/linter for gigmcp
// registry manifests and the signed index. The SAME code runs in registry CI
// and in the gateway, so the two can never drift (handoff §4.1). Apache-2.0
// (DESIGN #16).
package schema

import (
	"fmt"

	sigsyaml "sigs.k8s.io/yaml"
)

// Manifest is an author-declared entitlements manifest (DESIGN #7, §4).
type Manifest struct {
	SchemaVersion int          `json:"schemaVersion"`
	Name          string       `json:"name"` // unique; tool-namespace prefix; no "_"
	Version       string       `json:"version"`
	Source        Source       `json:"source"`
	Image         Image        `json:"image"`
	Tier          string       `json:"tier"` // sealed | entrusted (DESIGN #6)
	Entitlements  Entitlements `json:"entitlements"`
	Credentials   []Credential `json:"credentials,omitempty"`
	Tools         []Tool       `json:"tools,omitempty"`

	// Presentation/grouping metadata (DESIGN: catalog branding). NOT part of
	// the runtime/security contract — excluded from RuntimeHash, so backfilling
	// these never diverges a baked image from the signed index and never forces
	// an image rebuild or re-consent. The gateway falls back to a generated
	// monogram + slug where these are absent.
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"` // closed enum (see Validate, once backfilled)
	Icon        string `json:"icon,omitempty"`     // repo-hosted icons/<slug>.svg (signed provenance)
}

type Source struct {
	Repo string `json:"repo"`
	Tag  string `json:"tag"`
	// Package is the directory within the repo containing the server's main
	// package/module (generic builder runs `go build .` there). Default ".".
	Package string `json:"package,omitempty"`
}

type Image struct {
	Ref string `json:"ref"`
	// Digest is the OCI image-index digest (multi-arch: linux/amd64 +
	// linux/arm64). The gateway pins it and resolves the matching per-arch
	// image manifest at pull time, so one digest serves every host arch.
	Digest     string `json:"digest"`
	Entrypoint string `json:"entrypoint"` // absolute path of the static binary inside the image
	// Builder selects the registry build recipe (images/<builder>/Dockerfile).
	// Empty = "go-static" (static Go binary, FROM scratch). "toolpack" builds
	// the generic engine and bakes in the paired manifests/<name>/<version>.toolspec.yaml.
	// "node" and "python" produce runtime-rootfs images — installable only
	// once the gateway's rootfs sandbox extension ships.
	Builder string `json:"builder,omitempty"`
}

type Entitlements struct {
	Egress []string `json:"egress,omitempty"` // exact host or "*.suffix" — proxy.allowed() semantics
}

type Credential struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`     // oauth2 | api_key | basic | custom_env
	Provider string   `json:"provider"` // connector slug — drives least-privilege scopes + incremental consent
	Scopes   []string `json:"scopes,omitempty"`
	Inject   Inject   `json:"inject"`
	// Vendor is the canonical OAuth-app grouping key (e.g. "google" for the 12
	// gmail/googlecalendar/... connectors). Lets one operator OAuth app + one
	// user Connected Account cover a whole vendor family. Grouping/presentation
	// only — excluded from RuntimeHash. Provider stays the per-connector slug.
	Vendor string `json:"vendor,omitempty"`
}

// Inject carries the secret-delivery mode. Validate() enforces that exactly one mode is active per tier: sealed→Header+Format, entrusted→Env.
type Inject struct {
	Header string `json:"header,omitempty"`
	Format string `json:"format,omitempty"` // must contain "{token}"
	Env    string `json:"env,omitempty"`
}

type Tool struct {
	Name    string `json:"name"`
	Default bool   `json:"default"` // exposed by default (DESIGN #11)
}

// Parse decodes manifest YAML. Strict: unknown fields are errors, so typos in
// security-relevant fields (e.g. "egres:") cannot silently grant nothing.
// Structural invariants (tier, schemaVersion, digests, egress) require Validate().
func Parse(data []byte) (*Manifest, error) {
	var m Manifest
	if err := sigsyaml.UnmarshalStrict(data, &m); err != nil {
		return nil, fmt.Errorf("schema: parse manifest: %w", err)
	}
	return &m, nil
}
