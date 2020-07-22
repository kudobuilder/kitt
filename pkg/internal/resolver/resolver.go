package resolver

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/afero"

	o "github.com/kudobuilder/kitt/pkg/internal/apis/operator"
	"github.com/kudobuilder/kitt/pkg/internal/resolver/git"
	"github.com/kudobuilder/kitt/pkg/internal/resolver/url"
)

// Resolver retrieves an operator package into a file system.
type Resolver interface {
	Resolve(context.Context) (fs afero.Fs, remover func() error, err error)
}

// New returns a new resolver for the kind of reference provided by 'version'.
func New(operator o.Operator, version o.Version) (Resolver, error) {
	if version.Git != nil {
		source := findSource(operator.GitSources, version.Git.Source)
		if source == nil {
			return nil, fmt.Errorf("unknown git source %q", version.Git.Source)
		}

		// TODO: cache git sources to ensure that repositories are only cloned once per source
		resolver := git.NewResolver(source.URL, version.Git.Tag, version.Git.SHA, version.Git.Directory)

		return resolver, nil
	}

	if version.URL != nil {
		resolver := url.NewResolver(*version.URL)

		return resolver, nil
	}

	return nil, errors.New("unknown version resolver")
}

func findSource(sources []o.GitSource, name string) *o.GitSource {
	for _, source := range sources {
		if source.Name == name {
			return &source
		}
	}

	return nil
}
