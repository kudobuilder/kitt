package v1alpha1

// Operator describes the location of a KUDO operator.
type Operator struct {
	TypeMeta `yaml:",inline"`

	// Name of the operator.
	Name string `yaml:"name"`

	// GitSources are optional references to Git repositories.
	GitSources []GitSource `yaml:"gitSources,omitempty"`

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
	// OperatorVersion of the KUDO operator.
	OperatorVersion string `yaml:"operatorVersion"`

	// AppVersion of the KUDO operator, optional.
	AppVersion string `yaml:"appVersion,omitempty"`

	// Git specifies a version as a directory in a Git repository with a
	// specific tag.
	Git *Git `yaml:"git,omitempty"`

	// URL specifies a version as a URL of a package tarball.
	URL *string `yaml:"url,omitempty"`

	// SkipVerify can be used to skip the KUDO package verification step to publish
	// old versions of an operator that would fail the package verification
	SkipVerify *bool `yaml:"url,omitempty"`
}

// Git references a specific tag of a Git repository of a KUDO operator.
type Git struct {
	// Source references a 'GitSource' name. The source's Git repository is
	// cloned and the specified tag is checked out.
	Source string `yaml:"source"`

	// Directory where the KUDO operator is defined in the Git repository.
	Directory string `yaml:"directory"`

	// Tag of the KUDO operator version. Either this or 'SHA' has to be set.
	Tag string `yaml:"tag,omitempty"`

	// SHA of the KUDO operator version if a branch is used instead of
	// a tag. Either this or 'Tag' has to be set.
	SHA string `yaml:"sha,omitempty"`
}
