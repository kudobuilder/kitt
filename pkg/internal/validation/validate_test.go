package validation

import (
	"path/filepath"
	"testing"

	"github.com/kudobuilder/kitt/pkg/internal/apis/operator"
	"github.com/kudobuilder/kitt/pkg/internal/repo"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name    string
		pkg     repo.Package
		version operator.Version
		result  Result
	}{
		{
			name: "operatorVersion isn't semver",
			pkg: createPkg(t, `name: foo
operatorVersion: "1.0.0"`),
			version: operator.Version{
				OperatorVersion: "next",
			},
			result: Result{
				Errors: []string{"operatorVersion isn't semver"},
			},
		},
		{
			name: "operatorVersion doesn't match",
			pkg: createPkg(t, `name: foo
operatorVersion: "1.0.0"`),
			version: operator.Version{
				OperatorVersion: "1.1.0",
			},
			result: Result{
				Warnings: []string{"operatorVersion \"1.1.0\" doesn't match operatorVersion \"1.0.0\" in operator package"},
			},
		},
		{
			name: "appVersion isn't semver",
			pkg: createPkg(t, `name: foo
operatorVersion: "1.0.0"
appVersion: "1.0.0"`),
			version: operator.Version{
				OperatorVersion: "1.0.0",
				AppVersion:      "next",
			},
			result: Result{
				Errors: []string{"appVersion isn't semver"},
			},
		},
		{
			name: "appVersion doesn't match",
			pkg: createPkg(t, `name: foo
operatorVersion: "1.0.0"
appVersion: "1.0.0"`),
			version: operator.Version{
				OperatorVersion: "1.0.0",
				AppVersion:      "1.1.0",
			},
			result: Result{
				Warnings: []string{"appVersion \"1.1.0\" doesn't match appVersion \"1.0.0\" in operator package"},
			},
		},
		{
			name: "appVersion in reference but not in package",
			pkg: createPkg(t, `name: foo
operatorVersion: "1.0.0"`),
			version: operator.Version{
				OperatorVersion: "1.0.0",
				AppVersion:      "1.0.0",
			},
			result: Result{
				Warnings: []string{"appVersion provided but not set in operator package"},
			},
		},
		{
			name: "appVersion in package but not in reference",
			pkg: createPkg(t, `name: foo
operatorVersion: "1.0.0"
appVersion: "1.0.0"`),
			version: operator.Version{
				OperatorVersion: "1.0.0",
			},
			result: Result{
				Warnings: []string{"appVersion not provided but set in operator package"},
			},
		},
		{
			name: "everything matches",
			pkg: createPkg(t, `name: foo
operatorVersion: "1.0.0"
appVersion: "1.0.0"`),
			version: operator.Version{
				OperatorVersion: "1.0.0",
				AppVersion:      "1.0.0",
			},
			result: Result{},
		},
	}

	for _, test := range tests {
		var result Result

		validateVersion(test.version, test.pkg, &result)
		assert.Equal(t, test.result, result, test.name)
	}
}

func createPkg(t *testing.T, operator string) repo.Package {
	pkgFs := afero.NewMemMapFs()

	assert.NoError(
		t,
		afero.WriteFile(pkgFs, filepath.Join(string(filepath.Separator), "operator.yaml"), []byte(operator), 0644))
	assert.NoError(t, afero.WriteFile(pkgFs, filepath.Join(string(filepath.Separator), "params.yaml"), []byte{}, 0644))

	pkg, err := repo.NewPackage(pkgFs)
	assert.NoError(t, err)

	return pkg
}
