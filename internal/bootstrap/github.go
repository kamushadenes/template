package bootstrap

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

func fetchGitHub(githubRepo string) (string, error) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "github_templates")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	if _, err := git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:      githubRepo,
		Progress: os.Stdout,
	}); err != nil {
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}

	return tempDir, nil
}
