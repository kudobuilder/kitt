package encode

import (
	"fmt"

	"github.com/kudobuilder/kitt/pkg/apis/operator/v1alpha1"
	"github.com/kudobuilder/kitt/pkg/internal/apis/operator"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// FromFile reads a YAML file containing an 'Operator' in any of the supported
// external APIs.
func FromFile(path string) (operator.Operator, error) {
	fs := afero.NewReadOnlyFs(afero.NewOsFs())

	content, err := afero.ReadFile(fs, path)
	if err != nil {
		return operator.Operator{}, err
	}

	return fromYAML(content)
}

func fromYAML(input []byte) (operator.Operator, error) {
	tm := typeMeta{}

	if err := yaml.Unmarshal(input, &tm); err != nil {
		return operator.Operator{}, fmt.Errorf("could not determine the API version: %v", err)
	}

	switch tm.APIVersion {
	case "index.kudo.dev/v1alpha1":
		if tm.Kind != "Operator" {
			return operator.Operator{}, fmt.Errorf("unknown kind %q for API version %q", tm.Kind, tm.APIVersion)
		}

		o := v1alpha1.Operator{}

		if err := yaml.Unmarshal(input, &o); err != nil {
			return operator.Operator{}, fmt.Errorf("could not decode operator config: %v", err)
		}

		return ConvertV1Alpha1(o), nil
	default:
		return operator.Operator{}, fmt.Errorf("unknown API version %q", tm.APIVersion)
	}
}

type typeMeta struct {
	Kind       string `yaml:"kind,omitempty"`
	APIVersion string `yaml:"apiVersion,omitempty"`
}
