// Command registryctl is the registry CI tool: lint manifests, compile the
// signed index, sign and verify it. Thin shell over the schema package so CI
// and the gateway share one validator.
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gigmcp/registry/schema"
)

// stdout is the writer used for normal command output. Tests may swap it for
// a buffer to capture output without touching os.Stdout.
var stdout io.Writer = os.Stdout

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "registryctl:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: registryctl lint|lint-toolspecs|build-index|sign|verify|keygen|build-args ...")
	}
	switch args[0] {
	case "lint":
		if len(args) < 2 {
			return fmt.Errorf("usage: registryctl lint <manifests-dir> [-denylist file]")
		}
		fs := flag.NewFlagSet("lint", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		denyPath := fs.String("denylist", "", "path to exfil denylist file")
		if err := fs.Parse(args[2:]); err != nil {
			return err
		}
		var deny []string
		if *denyPath != "" {
			raw, err := os.ReadFile(*denyPath)
			if err != nil {
				return err
			}
			deny = strings.Split(string(raw), "\n")
		}
		manifests, err := loadAll(args[1], deny)
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "lint OK: %d manifest(s)\n", len(manifests))
		return nil

	case "lint-toolspecs":
		if len(args) < 2 {
			return fmt.Errorf("usage: registryctl lint-toolspecs <manifests-dir | manifests/<name>/<version>.toolspec.yaml>")
		}
		return lintToolspecs(args[1])

	case "build-index":
		if len(args) < 2 {
			return fmt.Errorf("usage: registryctl build-index <manifests-dir> [-out file]")
		}
		fs := flag.NewFlagSet("build-index", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		out := fs.String("out", "index.json", "output path")
		if err := fs.Parse(args[2:]); err != nil {
			return err
		}
		manifests, err := loadAll(args[1], nil)
		if err != nil {
			return err
		}
		ix, err := schema.BuildIndex(manifests, time.Now().UTC().Format(time.RFC3339))
		if err != nil {
			return err
		}
		raw, err := json.Marshal(ix)
		if err != nil {
			return err
		}
		return os.WriteFile(*out, raw, 0o644)

	case "sign":
		fs := flag.NewFlagSet("sign", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		in := fs.String("in", "index.json", "index file")
		out := fs.String("out", "index.json.sig", "signature output")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		key := os.Getenv("GIG_SIGNING_KEY")
		if key == "" {
			return fmt.Errorf("GIG_SIGNING_KEY must be set (64-byte ed25519 private key, hex)")
		}
		raw, err := os.ReadFile(*in)
		if err != nil {
			return err
		}
		sig, err := schema.Sign(key, raw)
		if err != nil {
			return err
		}
		return os.WriteFile(*out, []byte(sig+"\n"), 0o644)

	case "verify":
		fs := flag.NewFlagSet("verify", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		in := fs.String("in", "index.json", "index file")
		sigPath := fs.String("sig", "index.json.sig", "signature file")
		pub := fs.String("pub", "", "32-byte ed25519 public key, hex")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		raw, err := os.ReadFile(*in)
		if err != nil {
			return err
		}
		sig, err := os.ReadFile(*sigPath)
		if err != nil {
			return err
		}
		if err := schema.Verify(*pub, raw, string(sig)); err != nil {
			return err
		}
		fmt.Fprintln(stdout, "signature OK")
		return nil

	case "keygen":
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}
		// private key → GitHub Actions secret GIG_INDEX_SIGNING_KEY;
		// public key → gateway GIG_REGISTRY_PUBKEY / baked default.
		fmt.Fprintf(stdout, "public:  %s\nprivate: %s\n", hex.EncodeToString(pub), hex.EncodeToString(priv))
		return nil

	case "build-args":
		if len(args) < 2 {
			return fmt.Errorf("usage: registryctl build-args <manifests/<name> | manifests/<name>/<version>.yaml>")
		}
		target := args[1]
		info, err := os.Stat(target)
		if err != nil {
			return fmt.Errorf("build-args: %w", err)
		}
		var m *schema.Manifest
		if info.IsDir() {
			// Directory: load all manifests and pick the latest version.
			all, err := loadAll(target, nil)
			if err != nil {
				return fmt.Errorf("build-args: %w", err)
			}
			ix, err := schema.BuildIndex(all, "")
			if err != nil {
				return fmt.Errorf("build-args: %w", err)
			}
			// The server name is the last path component of the directory.
			name := filepath.Base(target)
			s, ok := ix.Servers[name]
			if !ok {
				// Fall back to the only server present (handles single-server dirs).
				for _, sv := range ix.Servers {
					s = sv
					break
				}
			}
			m = s.Versions[s.Latest]
		} else {
			// Specific .yaml file.
			raw, err := os.ReadFile(target)
			if err != nil {
				return fmt.Errorf("build-args: %w", err)
			}
			m, err = schema.Parse(raw)
			if err != nil {
				return fmt.Errorf("build-args: %w", err)
			}
			if err := m.Validate(); err != nil {
				return fmt.Errorf("build-args: %w", err)
			}
		}
		pkg := m.Source.Package
		if pkg == "" {
			pkg = "."
		}
		builder := m.Image.Builder
		if builder == "" {
			builder = "go-static"
		}
		fmt.Fprintf(stdout, "SOURCE_REPO=%s\n", m.Source.Repo)
		fmt.Fprintf(stdout, "SOURCE_TAG=%s\n", m.Source.Tag)
		fmt.Fprintf(stdout, "PACKAGE=%s\n", pkg)
		fmt.Fprintf(stdout, "BUILDER=%s\n", builder)
		fmt.Fprintf(stdout, "NAME=%s\n", m.Name)
		fmt.Fprintf(stdout, "VERSION=%s\n", m.Version)
		return nil

	default:
		return fmt.Errorf("unknown subcommand %q", args[0])
	}
}

// lintToolspecs validates toolspecs against their paired manifests. Toolspecs
// live NEXT TO their manifest: manifests/<name>/<version>.toolspec.yaml beside
// manifests/<name>/<version>.yaml. target is either the manifests dir
// (validates every spec, and requires a spec for every manifest selecting
// builder "toolpack") or a single *.toolspec.yaml (validates just that
// pairing — used by authoring agents to lint their own output without seeing
// siblings' work-in-progress).
func lintToolspecs(target string) error {
	lintOne := func(path string) (*schema.ToolSpec, error) {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		s, err := schema.ParseToolSpec(raw)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		if err := s.Validate(); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		wantDir, wantFile := s.Name, s.Version+".toolspec.yaml"
		if filepath.Base(filepath.Dir(path)) != wantDir || filepath.Base(path) != wantFile {
			return nil, fmt.Errorf("%s: path must be manifests/%s/%s", path, wantDir, wantFile)
		}
		// The paired manifest sits beside the spec; loading just it keeps
		// single-file mode free of sibling-directory parse interference.
		mPath := filepath.Join(filepath.Dir(path), s.Version+".yaml")
		rawM, err := os.ReadFile(mPath)
		if err != nil {
			return nil, fmt.Errorf("%s: no paired manifest: %w", path, err)
		}
		m, err := schema.Parse(rawM)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", mPath, err)
		}
		if err := m.Validate(); err != nil {
			return nil, fmt.Errorf("%s: %w", mPath, err)
		}
		if err := s.CheckAgainstManifest(m); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		return s, nil
	}

	info, err := os.Stat(target)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		if _, err := lintOne(target); err != nil {
			return err
		}
		fmt.Fprintln(stdout, "lint-toolspecs OK: 1 toolspec")
		return nil
	}

	specs := map[string]bool{}
	count := 0
	err = filepath.WalkDir(target, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Type()&os.ModeSymlink != 0 || !strings.HasSuffix(path, ".toolspec.yaml") {
			return nil
		}
		s, err := lintOne(path)
		if err != nil {
			return err
		}
		specs[s.Name+"@"+s.Version] = true
		count++
		return nil
	})
	if err != nil {
		return err
	}
	// Every manifest that selects the toolpack builder must have a spec —
	// without one the built image would have no /app/toolspec.yaml. (Dir mode
	// only; loadAll here is safe because dir mode runs in CI, not from
	// concurrent authoring agents.)
	manifests, err := loadAll(target, nil)
	if err != nil {
		return err
	}
	var missing []string
	for _, m := range manifests {
		if m.Image.Builder == "toolpack" && !specs[m.Name+"@"+m.Version] {
			missing = append(missing, m.Name+"@"+m.Version)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("manifests with builder \"toolpack\" but no toolspec: %s", strings.Join(missing, ", "))
	}
	fmt.Fprintf(stdout, "lint-toolspecs OK: %d toolspec(s)\n", count)
	return nil
}

// loadAll parses every manifests/<name>/<version>.yaml, enforces that the
// path matches the manifest's name/version (prevents PR sleight-of-hand),
// and lints when a denylist is given (Validate-only otherwise).
func loadAll(dir string, deny []string) ([]*schema.Manifest, error) {
	var manifests []*schema.Manifest
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// *.toolspec.yaml files sit beside manifests but are linted by
		// lint-toolspecs, not parsed as manifests.
		if d.IsDir() || d.Type()&os.ModeSymlink != 0 || !strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".toolspec.yaml") {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		m, err := schema.Parse(raw)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if deny != nil {
			if err := m.Lint(deny); err != nil {
				return fmt.Errorf("%s: %w", path, err)
			}
		} else if err := m.Validate(); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		wantDir, wantFile := m.Name, m.Version+".yaml"
		if filepath.Base(filepath.Dir(path)) != wantDir || filepath.Base(path) != wantFile {
			return fmt.Errorf("%s: path must be manifests/%s/%s", path, wantDir, wantFile)
		}
		manifests = append(manifests, m)
		return nil
	})
	if err == nil && len(manifests) == 0 {
		return nil, fmt.Errorf("no manifests found in %s", dir)
	}
	return manifests, err
}
