package schema

import (
	"fmt"
	"strconv"
	"strings"
)

// Index is the compiled, signed registry artifact. Runners consume ONLY this
// (never the raw repo, DESIGN #8). The signature covers the exact published
// index.json bytes — verify raw bytes, never re-serialized JSON.
type Index struct {
	SchemaVersion int                    `json:"schemaVersion"`
	Generated     string                 `json:"generated"` // RFC3339, set by CI
	Servers       map[string]IndexServer `json:"servers"`
}

type IndexServer struct {
	Latest   string               `json:"latest"`
	Versions map[string]*Manifest `json:"versions"`
}

// BuildIndex validates every manifest and compiles the index. Latest is the
// highest semver per server.
func BuildIndex(manifests []*Manifest, generated string) (*Index, error) {
	ix := &Index{SchemaVersion: 1, Generated: generated, Servers: map[string]IndexServer{}}
	for _, m := range manifests {
		if err := m.Validate(); err != nil {
			return nil, fmt.Errorf("manifest %s@%s: %w", m.Name, m.Version, err)
		}
		s := ix.Servers[m.Name]
		if s.Versions == nil {
			s.Versions = map[string]*Manifest{}
		}
		if _, dup := s.Versions[m.Version]; dup {
			return nil, fmt.Errorf("duplicate manifest %s@%s", m.Name, m.Version)
		}
		s.Versions[m.Version] = m
		if s.Latest == "" || semverLess(s.Latest, m.Version) {
			s.Latest = m.Version
		}
		ix.Servers[m.Name] = s
	}
	return ix, nil
}

// Resolve resolves "name" (latest), "name@version", or "sha256:<digest>"
// (exact image-digest lookup) — the installer ref grammar (handoff §6.B).
// If multiple manifests pin the same image digest, digest resolution is non-deterministic.
func (ix *Index) Resolve(ref string) (*Manifest, error) {
	if strings.HasPrefix(ref, "sha256:") {
		for _, s := range ix.Servers {
			for _, m := range s.Versions {
				if m.Image.Digest == ref {
					return m, nil
				}
			}
		}
		return nil, fmt.Errorf("schema: no manifest pins image digest %s", ref)
	}
	name, version, _ := strings.Cut(ref, "@")
	if strings.HasSuffix(ref, "@") {
		return nil, fmt.Errorf("schema: malformed ref %q (trailing @)", ref)
	}
	s, ok := ix.Servers[name]
	if !ok {
		return nil, fmt.Errorf("schema: unknown server %q", name)
	}
	if version == "" {
		version = s.Latest
	}
	m, ok := s.Versions[version]
	if !ok {
		return nil, fmt.Errorf("schema: unknown version %s@%s", name, version)
	}
	return m, nil
}

// semverLess reports a < b for versionRE-validated MAJOR.MINOR.PATCH strings.
func semverLess(a, b string) bool {
	pa, pb := strings.Split(a, "."), strings.Split(b, ".")
	if len(pa) < 3 || len(pb) < 3 {
		return false
	}
	for i := 0; i < 3; i++ {
		na, _ := strconv.Atoi(pa[i])
		nb, _ := strconv.Atoi(pb[i])
		if na != nb {
			return na < nb
		}
	}
	return false
}
