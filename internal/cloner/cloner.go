package cloner

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

// clones git repository to temp dir
type Cloner interface {
	Clone(ctx context.Context, repoURL string) (dir string, cleanup func(), err error)
}

type cloner struct{}

func NewCloner() Cloner {
	return &cloner{}
}

func (c *cloner) Clone(ctx context.Context, repoURL string) (string, func(), error) {
	url, err := normalizeAndValidateURL(repoURL)
	if err != nil {
		return "", nil, fmt.Errorf("git: invalid repo url: %w", err)
	}

	if err := checkGitAvailable(); err != nil {
		return "", nil, err
	}

	tmpDir, err := os.MkdirTemp("", "depcheck-*")
	if err != nil {
		return "", nil, fmt.Errorf("git: create temp dir: %w", err)
	}

	clean := func() { _ = os.RemoveAll(tmpDir) }

	cmd := exec.CommandContext(ctx, "git", "clone", "--depth=1", "--quiet", url, tmpDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		clean()
		return "", nil, fmt.Errorf("git: clone %q: %w\n%s", url, err, strings.TrimSpace(string(out)))
	}

	return tmpDir, clean, nil
}

func normalizeAndValidateURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("empty repo url")
	}

	if strings.HasPrefix(raw, "git@") {
		if !strings.Contains(raw, ":") {
			return "", fmt.Errorf("invalid ssh repo url: %q", raw)
		}
		return raw, nil
	}

	if !strings.HasPrefix(raw, "http://") &&
		!strings.HasPrefix(raw, "https://") &&
		!strings.HasPrefix(raw, "ssh://") {
		raw = "https://" + raw
	}

	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return "", fmt.Errorf("invalid repo url %q: %w", raw, err)
	}
	if u.Host == "" {
		return "", fmt.Errorf("invalid repo url %q: empty host", raw)
	}

	switch u.Scheme {
	case "http", "https", "ssh":
		return raw, nil
	default:
		return "", fmt.Errorf("unsupported repo url scheme %q", u.Scheme)
	}
}

func checkGitAvailable() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git: binary not found in PATH: %w", err)
	}
	return nil
}
