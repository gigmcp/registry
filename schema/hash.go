package schema

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// Hash is the canonical manifest hash: SHA-256 over the canonical JSON
// encoding of the PARSED manifest (struct field order is deterministic;
// no maps), hex-encoded. Basis for re-consent (DESIGN #7): a changed
// manifest changes the hash; YAML formatting changes do not.
func (m *Manifest) Hash() (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

// RuntimeHash is the canonical hash over ONLY the runtime/security-relevant
// fields of a manifest. It is the comparator the gateway uses for its bundle
// defense-in-depth check: it asserts RuntimeHash(indexManifest) ==
// RuntimeHash(bundledImageManifest) so a tampered image can't diverge its
// egress/credentials/tier/entrypoint/tools from what was signed + consented.
//
// It EXCLUDES, by clearing before hashing:
//   - image.digest — self-referential; image integrity already comes from the
//     pinned, signature-verified image-index digest.
//   - presentation/grouping metadata — credential.vendor and manifest
//     displayName/description/category/icon.
//
// Because those last fields are NOT part of this hash, backfilling vendor or
// branding into the manifests never diverges an already-built image from the
// signed index (the baked manifest still RuntimeHash-matches), so it needs no
// image rebuild and is safe to do incrementally. Everything else is included by
// clearing only the known-excluded fields (secure by default: a newly added
// runtime field is hashed unless explicitly excluded here).
//
// Both registry CI and the gateway MUST call this same function from the shared
// schema module so the two comparisons can never drift.
func (m *Manifest) RuntimeHash() (string, error) {
	r := *m // shallow copy; Image is a value field, so clearing its Digest can't touch m
	r.Image.Digest = ""
	r.DisplayName = ""
	r.Description = ""
	r.Category = ""
	r.Icon = ""
	// Credentials is a slice (shared backing array under the shallow copy):
	// reallocate before clearing per-credential Vendor so m is left untouched.
	if len(m.Credentials) > 0 {
		r.Credentials = make([]Credential, len(m.Credentials))
		copy(r.Credentials, m.Credentials)
		for i := range r.Credentials {
			r.Credentials[i].Vendor = ""
		}
	}
	b, err := json.Marshal(&r)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}
