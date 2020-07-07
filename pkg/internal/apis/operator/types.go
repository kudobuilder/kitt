package operator

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
	Version string

	// Git specifies a version as a directory in a Git repository with a
	// specific tag.
	Git *Git

	// URL specifies a version as a URL of a package tarball.
	URL *string
}

// Git references a specific tag of a Git repository of a KUDO operator.
type Git struct {
	// Source references a 'GitSource' name. The source's Git repository is
	// cloned and the specified tag is checked out.
	Source string

	// Tag of the KUDO operator version.
	Tag string

	// Directory where the KUDO operator is defined in the Git repository.
	Directory string
}
