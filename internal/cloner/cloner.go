package cloner

import (
	"context"
	"fmt"
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
	url := normaliseURL(repoURL)

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

func normaliseURL(raw string) string {
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") ||
		strings.HasPrefix(raw, "git@") || strings.HasPrefix(raw, "ssh://") {
		return raw
	}
	return "https://" + raw
}

func checkGitAvailable() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git: binary not found in PATH: %w", err)
	}
	return nil
}
