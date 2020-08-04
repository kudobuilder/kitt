package repo

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver/v3"
	kudo "github.com/kudobuilder/kudo/pkg/kudoctl/util/repo"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	repo := SyncedRepo{
		index: map[string]kudo.PackageVersions{
			"foo": {
				{
					Metadata: &kudo.Metadata{
						Name:            "foo",
						OperatorVersion: "1.0.0",
						AppVersion:      "1.0.0",
					},
				},
			},
		},
	}

	tests := []struct {
		name     string
		pkg      Package
		expected bool
	}{
		{
			name: "Unknown package",
			pkg: Package{
				OperatorName: "bar",
			},
			expected: false,
		},
		{
			name: "Package version not in repo",
			pkg: Package{
				OperatorName: "foo",
			},
			expected: false,
		},
		{
			name: "Package version in repo",
			pkg: Package{
				OperatorName:    "foo",
				OperatorVersion: *semver.MustParse("1.0.0"),
				AppVersion:      semver.MustParse("1.0.0"),
			},
			expected: true,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			actual := repo.Contains(test.pkg)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestAdd(t *testing.T) {
	// Using 'MemMapFs' with the default base path causes all kinds of trouble.
	// To avoid potential issues, all files are created in directories and
	// 'BasePathFs' is used to point to a different base path.
	repoFs := afero.NewMemMapFs()

	repoDir := filepath.Join(string(filepath.Separator), "repo")

	assert.NoError(t, repoFs.Mkdir(repoDir, 0755))
	repoFs = afero.NewBasePathFs(repoFs, repoDir)

	repo := SyncedRepo{
		fs: repoFs,
	}

	pkgFs := afero.NewMemMapFs()

	operatorDir := filepath.Join(string(filepath.Separator), "operator")

	assert.NoError(t, pkgFs.Mkdir(operatorDir, 0755))
	assert.NoError(t, afero.WriteFile(pkgFs, filepath.Join(operatorDir, "operator.yaml"), []byte(`name: foo
operatorVersion: "1.0.0"
appVersion: "1.0.0"
`), 0644))
	assert.NoError(t, afero.WriteFile(pkgFs, filepath.Join(operatorDir, "params.yaml"), []byte{}, 0644))

	pkgFs = afero.NewBasePathFs(pkgFs, operatorDir)

	pkg, err := NewPackage(pkgFs)
	assert.NoError(t, err)

	assert.False(t, repo.Contains(pkg))

	tarball, err := repo.Add(pkg)
	assert.NoError(t, err)

	assert.Equal(t, fmt.Sprintf("%s.tgz", pkg.String()), tarball)

	assert.True(t, repo.Contains(pkg))
}
