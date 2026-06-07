package schema

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

// CheckEgressEntry enforces EXACTLY the proxy's allowed() semantics
// (internal/proxy in gigmcp: case-insensitive exact host, or "*.suffix"
// suffix match): lowercase hostnames only; wildcards only as a leading "*."
// whose suffix keeps >= 2 labels; no ports, paths, or IP literals.
func CheckEgressEntry(e string) error {
	host := e
	if strings.HasPrefix(e, "*.") {
		host = strings.TrimPrefix(e, "*.")
		if strings.Count(host, ".") < 1 {
			return fmt.Errorf("egress %q: wildcard suffix must keep at least two labels (*.%s is too broad)", e, host)
		}
	}
	if strings.Contains(host, "*") {
		return fmt.Errorf("egress %q: bare or embedded wildcard not allowed", e)
	}
	if strings.ContainsAny(e, ":/") {
		return fmt.Errorf("egress %q: ports and paths not allowed (hostnames only)", e)
	}
	if net.ParseIP(host) != nil {
		return fmt.Errorf("egress %q: raw IP literals not allowed", e)
	}
	for _, label := range strings.Split(host, ".") {
		if !nameRE.MatchString(label) {
			return fmt.Errorf("egress %q: invalid label %q (lowercase a-z0-9, hyphens inside)", e, label)
		}
	}
	return nil
}

// Lint = Validate + registry policy: no egress entry may equal or be a
// subdomain of a denylisted exfil domain. Denylist lines starting with "#"
// or blank are ignored (the checked-in denylist file is passed line-split).
func (m *Manifest) Lint(denylist []string) error {
	if err := m.Validate(); err != nil {
		return err
	}
	var errs []error
	for _, e := range m.Entitlements.Egress {
		host := strings.TrimPrefix(e, "*.")
		for _, bad := range denylist {
			bad = strings.ToLower(strings.TrimSpace(bad))
			if bad == "" || strings.HasPrefix(bad, "#") {
				continue
			}
			if host == bad || strings.HasSuffix(host, "."+bad) {
				errs = append(errs, fmt.Errorf("egress %q matches denylist entry %q", e, bad))
			}
		}
	}
	return errors.Join(errs...)
}
