package repo

import (
	"fmt"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/kudobuilder/kudo/pkg/kudoctl/packages/reader"
	"github.com/spf13/afero"
)

// Package wraps a package file system with additional information contained
// in the package.
type Package struct {
	afero.Fs

	OperatorName    string
	OperatorVersion semver.Version
	AppVersion      *semver.Version
}

// NewPackage creates a new Package by extracting version information from a
// file system of an operator package.
// The base path of 'pkgFs' is assumed to contain an operator package.
func NewPackage(pkgFs afero.Fs) (Package, error) {
	p, err := reader.ReadDir(pkgFs, string(filepath.Separator))
	if err != nil {
		return Package{}, err
	}

	operatorName := p.Resources.Operator.Name

	operatorVersion, err := semver.NewVersion(p.Resources.OperatorVersion.Spec.Version)
	if err != nil {
		return Package{}, err
	}

	var appVersion *semver.Version

	// AppVersion is optional
	if p.Resources.OperatorVersion.Spec.AppVersion != "" {
		appVersion, err = semver.NewVersion(p.Resources.OperatorVersion.Spec.AppVersion)
		if err != nil {
			return Package{}, err
		}
	}

	return Package{
		Fs:              pkgFs,
		OperatorName:    operatorName,
		OperatorVersion: *operatorVersion,
		AppVersion:      appVersion,
	}, nil
}

// String returns a string representation of a package.
func (p Package) String() string {
	if p.AppVersion == nil {
		return fmt.Sprintf("%v-%v", p.OperatorName, p.OperatorVersion.String())
	}

	return fmt.Sprintf("%v-%v_%v", p.OperatorName, p.AppVersion.String(), p.OperatorVersion.String())
}

// Equal checks for equality of package versions.
func (p Package) Equal(other Package) bool {
	if p.OperatorName != other.OperatorName {
		return false
	}

	if !p.OperatorVersion.Equal(&other.OperatorVersion) {
		return false
	}

	if p.AppVersion == nil {
		return other.AppVersion == nil
	}

	if p.AppVersion != nil && other.AppVersion == nil {
		return false
	}

	if !p.AppVersion.Equal(other.AppVersion) {
		return false
	}

	return true
}
