package git

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	testClone := func(ctx context.Context, tempDir, url, branch, sha string) error {
		if url == "example.org" && branch == "test" && sha == "" {
			return nil
		}

		return errors.New("wrong URL and branch")
	}

	resolver := &Resolver{
		URL:               "example.org",
		Branch:            "test",
		OperatorDirectory: "operator",
		gitClone:          testClone,
	}

	_, remover, err := resolver.Resolve(context.Background())
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, remover())
	}()
}
