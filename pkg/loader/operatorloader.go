package loader

import (
	"fmt"

	"github.com/kudobuilder/kitt/pkg/internal/apis/operator"
	"github.com/kudobuilder/kitt/pkg/internal/apis/operator/encode"
)

// OperatorLoader allows to gather operators from different sources.
type OperatorLoader interface {
	Apply() ([]operator.Operator, error)
}

type operatorLoaderAdapter func() ([]operator.Operator, error)

func (f operatorLoaderAdapter) Apply() ([]operator.Operator, error) {
	return f()
}

// FromFiles reads operator definitions from multiple YAML files.
func FromFiles(paths []string) OperatorLoader {
	return operatorLoaderAdapter(func() ([]operator.Operator, error) {
		operators := make([]operator.Operator, 0, len(paths))

		for _, path := range paths {
			o, err := encode.FromFile(path)
			if err != nil {
				return operators, fmt.Errorf("failed to read %q: %v", path, err)
			}

			operators = append(operators, o)
		}

		return operators, nil
	})
}
