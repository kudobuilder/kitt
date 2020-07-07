package update

import (
	"fmt"

	"github.com/kudobuilder/kitt/pkg/internal/apis/operator"
	"github.com/kudobuilder/kitt/pkg/internal/apis/operator/encode"
)

// OperatorOption allows to gather operators from different sources.
type OperatorOption interface {
	apply() ([]operator.Operator, error)
}

type operatorOptionAdapter func() ([]operator.Operator, error)

func (f operatorOptionAdapter) apply() ([]operator.Operator, error) {
	return f()
}

// FromFiles reads operator definitions from multiple YAML files.
func FromFiles(paths []string) OperatorOption {
	return operatorOptionAdapter(func() ([]operator.Operator, error) {
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
