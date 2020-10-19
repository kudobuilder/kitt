package operator

import "fmt"

// Operator describes the location of a KUDO operator.
type Operator struct {
	// Name of the operator.
	Name string

	// GitSources are optional references to Git repositories.
	GitSources []GitSource

	// Versions of the operator.
	Versions []Version
}

// GitSource is the location of a git repository.
type GitSource struct {
	// Name of this source. This name is referenced by 'Version' entries setting
	// a 'Git' field.
	Name string

	URL string
}

// Version describes a version of a KUDO operator.
type Version struct {
	// OperatorVersion of the KUDO operator.
	OperatorVersion string

	// AppVersion of the KUDO operator, optional.
	AppVersion string

	// Git specifies a version as a directory in a Git repository with a
	// specific tag.
	Git *Git

	// URL specifies a version as a URL of a package tarball.
	URL *string

	// SkipVerify can be used to skip the KUDO package verification step to publish
	// old versions of an operator that would fail the package verification
	SkipVerify bool
}

// Version prints the version as a combination of appVersion and operatorVersion
// as described in KEP-19.
func (v Version) Version() string {
	if v.AppVersion != "" {
		return fmt.Sprintf("%s_%s", v.AppVersion, v.OperatorVersion)
	}

	return v.OperatorVersion
}

// Git references a specific tag of a Git repository of a KUDO operator.
type Git struct {
	// Source references a 'GitSource' name. The source's Git repository is
	// cloned and the specified tag is checked out.
	Source string

	// Directory where the KUDO operator is defined in the Git repository.
	Directory string

	// Tag of the KUDO operator version. Either this or 'SHA' has to be set.
	Tag string

	// SHA of the KUDO operator version if a branch is used instead of
	// a tag. Either this or 'Tag' has to be set.
	SHA string
}
