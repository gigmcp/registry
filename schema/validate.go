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
)

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
	case "", "go-static", "node", "python":
		// valid
	default:
		add("image.builder %q invalid: must be one of \"go-static\", \"node\", \"python\" (or omitted for go-static default)", m.Image.Builder)
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
