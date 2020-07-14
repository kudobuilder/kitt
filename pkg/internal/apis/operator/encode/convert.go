package encode

import (
	"github.com/kudobuilder/kitt/pkg/apis/operator/v1alpha1"
	"github.com/kudobuilder/kitt/pkg/internal/apis/operator"
)

// ConvertV1Alpha1 creates an internal 'Operator' instance from the external
// v1alpha1 API.
func ConvertV1Alpha1(in v1alpha1.Operator) operator.Operator {
	out := operator.Operator{
		Name:       in.Name,
		GitSources: make([]operator.GitSource, len(in.GitSources)),
		Versions:   make([]operator.Version, len(in.Versions)),
	}

	for i := range in.GitSources {
		out.GitSources[i] = convertV1Alpha1GitSource(in.GitSources[i])
	}

	for i := range in.Versions {
		out.Versions[i] = convertV1Alpha1Version(in.Versions[i])
	}

	return out
}

func convertV1Alpha1GitSource(in v1alpha1.GitSource) operator.GitSource {
	out := operator.GitSource{
		Name: in.Name,
		URL:  in.URL,
	}

	return out
}

func convertV1Alpha1Version(in v1alpha1.Version) operator.Version {
	out := operator.Version{
		Version: in.Version,
		URL:     in.URL,
	}

	if in.Git != nil {
		git := convertV1Alpha1Git(*in.Git)
		out.Git = &git
	}

	return out
}

func convertV1Alpha1Git(in v1alpha1.Git) operator.Git {
	out := operator.Git{
		Source:    in.Source,
		Directory: in.Directory,
		Tag:       in.Tag,
		SHA:       in.SHA,
	}

	return out
}
