package git

import (
	"bufio"
	"context"
	"os/exec"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/kudobuilder/kitt/pkg/internal/resolvers"
)

// Resolver resolves operator packages from a Git repository.
type Resolver struct {
	URL               string
	Branch            string
	OperatorDirectory string

	// Extracted function to simplify testing.
	gitClone func(ctx context.Context, url, branch, tempDir string) error
}

// NewResolver creates a new Resolver for a Git repository at the specified URL.
func NewResolver(url, branch, operatorDirectory string) Resolver {
	return Resolver{
		URL:               url,
		Branch:            branch,
		OperatorDirectory: operatorDirectory,

		gitClone: gitClone,
	}
}

// Resolve clones a specific branch of a git repository and returns a file system
// pointing at the operator directory.
// The repository is cloned into a temporary directory. Callers are responsible
// for removing this directory by running the returned remover function.
func (r Resolver) Resolve(ctx context.Context) (afero.Fs, resolvers.Remover, error) {
	fs := afero.NewOsFs()

	tempDir, err := afero.TempDir(fs, "", "")
	if err != nil {
		return nil, nil, err
	}

	log.WithField("directory", tempDir).
		Debug("Created temporary directory")

	log.WithField("url", r.URL).
		WithField("branch", r.Branch).
		Info("Cloning Git repository")

	if err := r.gitClone(ctx, r.URL, r.Branch, tempDir); err != nil {
		return nil, nil, err
	}

	remover := func() error {
		log.WithField("directory", tempDir).
			Debug("Removing temporary directory")
		return fs.RemoveAll(tempDir)
	}

	return afero.NewBasePathFs(fs, path.Join(tempDir, r.OperatorDirectory)), remover, nil
}

func gitClone(ctx context.Context, url, branch, tempDir string) error {
	//nolint:gosec
	cmd := exec.CommandContext(ctx, "git", "clone", "--branch", branch, "--single-branch", url, tempDir)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Output git's command output while running the command.
	scanner := bufio.NewScanner(stderr)

	for scanner.Scan() {
		t := scanner.Text()
		log.WithField("url", url).
			WithField("branch", branch).
			Debug(t)
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
