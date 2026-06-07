package schema

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"strings"
)

// Sign returns the hex-encoded ed25519 detached signature over data (the
// exact index.json bytes to publish). privHex: 64-byte private key, hex
// (lives in a GitHub Actions secret).
func Sign(privHex string, data []byte) (string, error) {
	key, err := hex.DecodeString(strings.TrimSpace(privHex))
	if err != nil || len(key) != ed25519.PrivateKeySize {
		return "", fmt.Errorf("schema: signing key must be %d hex-encoded bytes", ed25519.PrivateKeySize)
	}
	return hex.EncodeToString(ed25519.Sign(ed25519.PrivateKey(key), data)), nil
}

// Verify checks sigHex over data with a 32-byte hex public key (baked into
// the gateway / GIG_REGISTRY_PUBKEY). MUST pass before the index is parsed.
func Verify(pubHex string, data []byte, sigHex string) error {
	pub, err := hex.DecodeString(strings.TrimSpace(pubHex))
	if err != nil || len(pub) != ed25519.PublicKeySize {
		return fmt.Errorf("schema: public key must be %d hex-encoded bytes", ed25519.PublicKeySize)
	}
	sig, err := hex.DecodeString(strings.TrimSpace(sigHex))
	if err != nil || len(sig) != ed25519.SignatureSize {
		return fmt.Errorf("schema: signature must be %d hex-encoded bytes", ed25519.SignatureSize)
	}
	if !ed25519.Verify(ed25519.PublicKey(pub), data, sig) {
		return fmt.Errorf("schema: index signature verification FAILED")
	}
	return nil
}
