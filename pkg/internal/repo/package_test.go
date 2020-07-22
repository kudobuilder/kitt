package repo

import (
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNewPackage(t *testing.T) {
	pkgFs := afero.NewMemMapFs()

	assert.NoError(t, afero.WriteFile(pkgFs, filepath.Join(string(filepath.Separator), "operator.yaml"), []byte(`name: foo
operatorVersion: "1.0.0"
appVersion: "1.0.0"`), 0644))
	assert.NoError(t, afero.WriteFile(pkgFs, filepath.Join(string(filepath.Separator), "params.yaml"), []byte{}, 0644))

	pkg, err := NewPackage(pkgFs)
	assert.NoError(t, err)

	assert.Equal(t, "foo", pkg.OperatorName)
	assert.Equal(t, *semver.MustParse("1.0.0"), pkg.OperatorVersion)
	assert.Equal(t, semver.MustParse("1.0.0"), pkg.AppVersion)
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name     string
		a        Package
		b        Package
		expected bool
	}{
		{
			name: "operator name doesn't match",
			a: Package{
				OperatorName: "foo",
			},
			b: Package{
				OperatorName: "bar",
			},
			expected: false,
		},
		{
			name: "operator version doesn't match",
			a: Package{
				OperatorVersion: *semver.MustParse("1.0.0"),
			},
			b: Package{
				OperatorVersion: *semver.MustParse("1.1.0"),
			},
			expected: false,
		},
		{
			name: "a has app version, b doesn't",
			a: Package{
				AppVersion: semver.MustParse("1.0.0"),
			},
			b:        Package{},
			expected: false,
		},
		{
			name: "b has app version, a doesn't",
			a:    Package{},
			b: Package{
				AppVersion: semver.MustParse("1.0.0"),
			},
			expected: false,
		},
		{
			name:     "a and b are equal, no app version",
			a:        Package{},
			b:        Package{},
			expected: true,
		},
		{
			name: "a and b are equal, with app version",
			a: Package{
				AppVersion: semver.MustParse("1.0.0"),
			},
			b: Package{
				AppVersion: semver.MustParse("1.0.0"),
			},
			expected: true,
		},
	}

	for _, test := range tests {
		actual := test.a.Equal(test.b)
		assert.Equal(t, test.expected, actual, test.name)
	}
}
