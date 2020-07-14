package v1alpha1

// Operator describes the location of a KUDO operator.
type Operator struct {
	TypeMeta `yaml:",inline"`

	// Name of the operator.
	Name string `yaml:"name"`

	// GitSources are optional references to Git repositories.
	GitSources []GitSource `yaml:"git-sources,omitempty"`

	// Versions of the operator.
	Versions []Version `yaml:"versions"`
}

// TypeMeta partially copies apimachinery/pkg/apis/meta/v1.TypeMeta
// No need for a direct dependence; the fields are stable.
type TypeMeta struct {
	Kind       string `yaml:"kind,omitempty"`
	APIVersion string `yaml:"apiVersion,omitempty"`
}

// GitSource is the location of a git repository.
type GitSource struct {
	// Name of this source. This name is referenced by 'Version' entries setting
	// a 'Git' field.
	Name string `yaml:"name"`

	URL string `yaml:"url"`
}

// Version describes a version of a KUDO operator.
type Version struct {
	Version string `yaml:"version"`

	// Git specifies a version as a directory in a Git repository with a
	// specific tag.
	Git *Git `yaml:"git,omitempty"`

	// URL specifies a version as a URL of a package tarball.
	URL *string `yaml:"url,omitempty"`
}

// Git references a specific tag of a Git repository of a KUDO operator.
type Git struct {
	// Source references a 'GitSource' name. The source's Git repository is
	// cloned and the specified tag is checked out.
	Source string `yaml:"source"`

	// Directory where the KUDO operator is defined in the Git repository.
	Directory string `yaml:"directory"`

	// Tag (or branch) of the KUDO operator version.
	Tag string `yaml:"tag,omitempty"`

	// Optional SHA of the KUDO operator version if a branch is used instead of
	// a tag. If this isn't set, the latest commit of the referenced branch will
	// be used.
	SHA string `yaml:"sha,omitempty"`
}
