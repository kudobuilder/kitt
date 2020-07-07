package git

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	testClone := func(ctx context.Context, url, branch, tempDir string) error {
		if url == "example.org" && branch == "test" {
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
