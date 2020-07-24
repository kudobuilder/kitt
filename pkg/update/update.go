package update

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/kudobuilder/kitt/pkg/internal/apis/operator"
	"github.com/kudobuilder/kitt/pkg/internal/repo"
	"github.com/kudobuilder/kitt/pkg/internal/resolver"
	"github.com/kudobuilder/kitt/pkg/loader"
)

// Update resolves a list of operators and adds them to a repository.
func Update(
	ctx context.Context,
	operatorLoader loader.OperatorLoader,
	repoPath string,
	repoURL string,
	force bool,
) error {
	repoFs := afero.NewBasePathFs(afero.NewOsFs(), repoPath)

	isDir, err := afero.IsDir(repoFs, "")
	if err != nil {
		return fmt.Errorf("failed to open repository path %q: %v", repoPath, err)
	}

	if !isDir {
		return fmt.Errorf("repository path %q is not a directory", repoPath)
	}

	syncedRepo, err := repo.NewSyncedRepo(repoFs, repoURL)
	if err != nil {
		return fmt.Errorf("failed to open repository %q: %v", repoPath, err)
	}

	operators, err := operatorLoader.Apply()
	if err != nil {
		return fmt.Errorf("failed to load operator configurations: %v", err)
	}

	for _, operator := range operators {
		for _, version := range operator.Versions {
			log.WithField("operator", operator.Name).
				WithField("version", version.Version()).
				WithField("repository", repoURL).
				WithField("path", repoPath).
				Info("Updating operator")

			if err := updateOperator(ctx, operator, version, syncedRepo, force); err != nil {
				return err
			}
		}
	}

	return nil
}

func updateOperator(
	ctx context.Context,
	operator operator.Operator,
	version operator.Version,
	syncedRepo *repo.SyncedRepo,
	force bool,
) (err error) {
	operatorName := fmt.Sprintf("%s-%s", operator.Name, version.Version())

	resolver, err := resolver.New(operator, version)
	if err != nil {
		return fmt.Errorf("failed to resolve operator %q: %v", operatorName, err)
	}

	pkgFs, remover, err := resolver.Resolve(ctx)
	if err != nil {
		return fmt.Errorf("failed to resolve operator %q: %v", operatorName, err)
	}

	// The package resolver created a temporary directory for the package file system.
	// We remove it once we no longer need it.
	defer func() {
		if rerr := remover(); rerr != nil {
			err = fmt.Errorf("failed to remove temporary directory of operator %q: %v", operatorName, rerr)
		}
	}()

	pkg, err := repo.NewPackage(pkgFs)
	if err != nil {
		return fmt.Errorf("failed to extract package version of operator %q: %v", operatorName, err)
	}

	contains := syncedRepo.Contains(pkg)

	if !contains || force {
		pkgName, err := syncedRepo.Add(pkg)
		if err != nil {
			return fmt.Errorf("failed to add operator %q to the repository: %v", pkg.String(), err)
		}

		log.WithField("operator", operator.Name).
			WithField("version", version.Version()).
			WithField("repository", syncedRepo.URL).
			WithField("tarball", pkgName).
			Info("Added operator to the repository")
	} else {
		log.WithField("operator", operator.Name).
			WithField("version", version.Version()).
			WithField("repository", syncedRepo.URL).
			Info("Operator is already in the repository")
	}

	return nil
}
