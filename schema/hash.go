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
