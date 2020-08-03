package git

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	tests := []struct {
		name      string
		cloneFake func(context.Context, string, string, string, string) error
		branch    string
		sha       string
		expectErr bool
	}{
		{
			name: "resolve branch",
			cloneFake: func(ctx context.Context, tempDir, url, branch, sha string) error {
				if url == "example.org" && branch == "test" && sha == "" {
					return nil
				}

				return errors.New("wrong URL and branch")
			},
			branch:    "test",
			sha:       "",
			expectErr: false,
		},
		{
			name: "resolve SHA",
			cloneFake: func(ctx context.Context, tempDir, url, branch, sha string) error {
				if url == "example.org" && branch == "" && sha == "abcdefg" {
					return nil
				}

				return errors.New("wrong URL and SHA")
			},
			branch:    "",
			sha:       "abcdefg",
			expectErr: false,
		},
		{
			name:      "neither branch nor SHA set",
			cloneFake: nil,
			branch:    "",
			sha:       "",
			expectErr: true,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			resolver := &Resolver{
				URL:               "example.org",
				Branch:            test.branch,
				SHA:               test.sha,
				OperatorDirectory: "operator",
				gitClone:          test.cloneFake,
			}

			_, remover, err := resolver.Resolve(context.Background())

			if remover != nil {
				defer func() {
					assert.NoError(t, remover())
				}()
			}

			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
