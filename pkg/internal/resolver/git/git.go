package git

import (
	"bufio"
	"context"
	"errors"
	"os/exec"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Resolver resolves operator packages from a Git repository.
type Resolver struct {
	URL               string
	Branch            string
	SHA               string
	OperatorDirectory string

	// Extracted function to simplify testing.
	gitClone func(ctx context.Context, tempDir, url, branch, sha string) error
}

// NewResolver creates a new Resolver for a Git repository at the specified URL.
func NewResolver(url, branch, sha string, operatorDirectory string) Resolver {
	return Resolver{
		URL:               url,
		Branch:            branch,
		SHA:               sha,
		OperatorDirectory: operatorDirectory,

		gitClone: gitClone,
	}
}

// Resolve clones a specific branch of a git repository and returns a file system
// pointing at the operator directory.
// The repository is cloned into a temporary directory. Callers are responsible
// for removing this directory by running the returned remover function.
func (r Resolver) Resolve(ctx context.Context) (afero.Fs, func() error, error) {
	fs := afero.NewOsFs()

	if r.Branch == "" && r.SHA == "" {
		return nil, nil, errors.New("neither branch nor SHA provided")
	}

	tempDir, err := afero.TempDir(fs, "", "")
	if err != nil {
		return nil, nil, err
	}

	log.WithField("directory", tempDir).
		Debug("Created temporary directory")

	log.WithField("url", r.URL).
		WithField("branch", r.Branch).
		WithField("sha", r.SHA).
		Info("Cloning Git repository")

	if err := r.gitClone(ctx, tempDir, r.URL, r.Branch, r.SHA); err != nil {
		return nil, nil, err
	}

	remover := func() error {
		log.WithField("directory", tempDir).
			Debug("Removing temporary directory")
		return fs.RemoveAll(tempDir)
	}

	return afero.NewBasePathFs(fs, path.Join(tempDir, r.OperatorDirectory)), remover, nil
}

func gitClone(ctx context.Context, tempDir, url, branch, sha string) error {
	if branch != "" {
		logger := log.WithField("url", url).WithField("branch", branch)

		if err := runAndLog(ctx, logger, "git", "clone", "--branch", branch, "--single-branch", url, tempDir); err != nil {
			return err
		}
	} else {
		logger := log.WithField("url", url).WithField("sha", sha)

		if err := runAndLog(ctx, logger, "git", "clone", url, tempDir); err != nil {
			return err
		}

		if err := runAndLog(ctx, logger, "git", "-C", tempDir, "checkout", sha); err != nil {
			return err
		}
	}

	return nil
}

func runAndLog(ctx context.Context, logger *log.Entry, name string, args ...string) error {
	//nolint:gosec
	cmd := exec.CommandContext(ctx, name, args...)

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
		logger.Debug(t)
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
