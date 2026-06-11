package schema

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	TierSealed    = "sealed"
	TierEntrusted = "entrusted"
)

var (
	nameRE    = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)
	versionRE = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`)
	digestRE  = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
	packageRE = regexp.MustCompile(`^[a-zA-Z0-9._/-]+$`)
	// vendorRE is the canonical OAuth-app grouping key shape: a lowercase slug,
	// same alphabet as a server name (lets "google" group the 12 google-* connectors).
	vendorRE = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)
)

// categories is the closed enum for the presentation-only Category field. It is
// enforced ONLY when Category is set (the field is optional/omitempty); absent
// categories are fine and the gateway falls back to no badge. The set is sized
// to the catalog so every server gets an honest bucket (catalog filtering),
// not force-fit into a handful of generic ones.
var categories = map[string]bool{
	"ai": true, "analytics": true, "comms": true, "CRM": true,
	"data": true, "design": true, "dev-tools": true, "documents": true,
	"e-commerce": true, "finance": true, "HR/ATS": true, "marketing": true,
	"productivity": true, "search": true, "social": true, "storage": true,
	"support": true,
}

// Validate checks structural rules. Registry-policy rules that need external
// data (the exfil denylist) live in Lint.
func (m *Manifest) Validate() error {
	var errs []error
	add := func(format string, a ...any) { errs = append(errs, fmt.Errorf(format, a...)) }

	if m.SchemaVersion != 1 {
		add("schemaVersion must be 1, got %d", m.SchemaVersion)
	}
	if !nameRE.MatchString(m.Name) {
		add("name %q invalid: lowercase [a-z0-9-], no underscores (name is the tool-namespace prefix)", m.Name)
	}
	if !versionRE.MatchString(m.Version) {
		add("version %q invalid: must be MAJOR.MINOR.PATCH", m.Version)
	}
	if m.Source.Repo == "" || m.Source.Tag == "" {
		add("source.repo and source.tag are required")
	}
	if m.Source.Package != "" {
		if !packageRE.MatchString(m.Source.Package) {
			add("source.package %q invalid: must match ^[a-zA-Z0-9._/-]+$", m.Source.Package)
		} else if strings.Contains(m.Source.Package, "..") {
			add("source.package %q invalid: must not contain \"..\" (path injection guard)", m.Source.Package)
		} else if strings.HasPrefix(m.Source.Package, "/") {
			add("source.package %q invalid: must not be an absolute path (path injection guard)", m.Source.Package)
		}
	}
	if m.Image.Ref == "" {
		add("image.ref is required")
	}
	if !digestRE.MatchString(m.Image.Digest) {
		add("image.digest %q invalid: must be sha256:<64 hex> (the approved digest is what runs)", m.Image.Digest)
	}
	if !strings.HasPrefix(m.Image.Entrypoint, "/") {
		add("image.entrypoint %q must be an absolute path inside the image", m.Image.Entrypoint)
	}
	switch m.Image.Builder {
	case "", "go-static", "node", "python", "toolpack":
		// valid
	default:
		add("image.builder %q invalid: must be one of \"go-static\", \"node\", \"python\", \"toolpack\" (or omitted for go-static default)", m.Image.Builder)
	}
	if m.Tier != TierSealed && m.Tier != TierEntrusted {
		add("tier %q invalid: must be %q or %q", m.Tier, TierSealed, TierEntrusted)
	}
	for _, e := range m.Entitlements.Egress {
		if err := CheckEgressEntry(e); err != nil {
			errs = append(errs, err)
		}
	}
	seenCred := map[string]bool{}
	for _, c := range m.Credentials {
		if c.ID == "" || seenCred[c.ID] {
			add("credential id %q empty or duplicate", c.ID)
		}
		seenCred[c.ID] = true
		switch m.Tier {
		case TierSealed:
			if c.Inject.Header == "" || !strings.Contains(c.Inject.Format, "{token}") || c.Inject.Env != "" {
				add("credential %q: sealed tier requires inject.header + inject.format containing {token}, and no inject.env", c.ID)
			}
		case TierEntrusted:
			if c.Inject.Env == "" || c.Inject.Header != "" || c.Inject.Format != "" {
				add("credential %q: entrusted tier requires inject.env only", c.ID)
			}
		}
		// Vendor is the canonical OAuth-app grouping key (one operator app + one
		// user Connected Account per vendor family). Required for oauth2 so the
		// gateway can group connectors; optional for other credential types.
		// Validated for shape whenever present.
		if c.Type == "oauth2" && c.Vendor == "" {
			add("credential %q: vendor is required for oauth2 (canonical OAuth-app grouping key, e.g. \"google\")", c.ID)
		}
		if c.Vendor != "" && !vendorRE.MatchString(c.Vendor) {
			add("credential %q: vendor %q invalid: lowercase slug [a-z0-9-]", c.ID, c.Vendor)
		}
	}
	if m.Category != "" && !categories[m.Category] {
		add("category %q invalid: must be one of ai, analytics, comms, CRM, data, design, dev-tools, documents, e-commerce, finance, HR/ATS, marketing, productivity, search, social, storage, support", m.Category)
	}
	// Icon, when present, must be the repo-hosted path icons/<name>.svg — tied to
	// the server name so a ref can never point at an unrelated/foreign asset. File
	// existence (the asset is actually committed) is checked by the linter, which
	// has the repo root; Validate only enforces the shape.
	if m.Icon != "" && m.Icon != "icons/"+m.Name+".svg" {
		add("icon %q invalid: must be %q (repo-hosted, signed-provenance asset tied to the server name)", m.Icon, "icons/"+m.Name+".svg")
	}
	seenTool := map[string]bool{}
	for _, tl := range m.Tools {
		if tl.Name == "" || seenTool[tl.Name] {
			add("tool name %q empty or duplicate", tl.Name)
		}
		seenTool[tl.Name] = true
	}
	return errors.Join(errs...)
}
