package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Masterminds/semver"
	"github.com/kudobuilder/kudo/pkg/kudoctl/packages/writer"
	kudo "github.com/kudobuilder/kudo/pkg/kudoctl/util/repo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// SyncedRepo manages operator packages of a file system.
// If packages are added, an updated index file is created.
type SyncedRepo struct {
	fs    afero.Fs
	index map[string]kudo.PackageVersions

	URL string
}

// NewSyncedRepo create a new repository in a file system.
// The repoURL parameter determines the URL to use for package file URLs
// in the repository index.
func NewSyncedRepo(fs afero.Fs, repoURL string) (*SyncedRepo, error) {
	index := map[string]kudo.PackageVersions{}

	indexExists, err := afero.Exists(fs, "index.yaml")
	if err != nil {
		return nil, err
	}

	if indexExists {
		log.WithField("repository", repoURL).
			Debug("Using existing index file")

		indexFile, err := afero.ReadFile(fs, "index.yaml")
		if err != nil {
			return nil, err
		}

		i, err := kudo.ParseIndexFile(indexFile)
		if err != nil {
			return nil, err
		}

		index = i.Entries
	}

	return &SyncedRepo{
		fs:    fs,
		index: index,
		URL:   repoURL,
	}, nil
}

// Contains checks if a specific operator package is in the repository.
func (s SyncedRepo) Contains(pkg Package) bool {
	entries, ok := s.index[pkg.OperatorName]
	if ok {
		for _, entry := range entries {
			var appVersion *semver.Version

			// AppVersion is optional
			if entry.AppVersion != "" {
				appVersion = semver.MustParse(entry.AppVersion)
			}

			if pkg.Equal(Package{
				OperatorName:    entry.Name,
				OperatorVersion: *semver.MustParse(entry.OperatorVersion),
				AppVersion:      appVersion,
			}) {
				return true
			}
		}
	}

	return false
}

// Add adds an operator package to the repository.
// The package contents are provided as a file system.
func (s *SyncedRepo) Add(pkg Package) (tarballName string, err error) {
	tarballName = fmt.Sprintf("%s.tgz", pkg.String())

	log.WithField("repository", s.URL).
		WithField("tarball", tarballName).
		Debug("Creating operator package")

	tarball, err := s.fs.OpenFile(tarballName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create operator package %q : %v", tarballName, err)
	}

	defer func() {
		if cerr := tarball.Close(); cerr != nil {
			err = cerr
		}
	}()

	// Path needs to be an empty string, otherwise wrong filenames will be created
	if err := writer.TgzDir(pkg, "", tarball); err != nil {
		return "", fmt.Errorf("failed to tar operator package %q: %v", tarballName, err)
	}

	now := time.Now()

	newIndex, err := kudo.IndexDirectory(s.fs, string(filepath.Separator), s.URL, &now)
	if err != nil {
		return "", fmt.Errorf("failed to create new index: %v", err)
	}

	s.index = newIndex.Entries

	log.WithField("repository", s.URL).
		WithField("tarball", tarballName).
		Debug("Writing new index file")

	// The repository file system is the source of truth.
	// We don't merge an existing index but overwrite instead.
	if err := newIndex.WriteFile(s.fs, "index.yaml"); err != nil {
		return "", fmt.Errorf("failed to write new index: %v", err)
	}

	return tarballName, nil
}
