package validate

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/kudobuilder/kitt/pkg/internal/apis/operator"
	"github.com/kudobuilder/kitt/pkg/internal/repo"
	"github.com/kudobuilder/kitt/pkg/internal/resolver"
	"github.com/kudobuilder/kitt/pkg/internal/validation"
	"github.com/kudobuilder/kitt/pkg/loader"
)

// Validate runs several checks on the operator reference as well as the
// referenced package. It checks that metadata provided in the reference is
// consistent with the metadata provided in the referenced package and also
// verifies all referenced packages.
func Validate(
	ctx context.Context,
	operatorLoader loader.OperatorLoader,
	strict bool,
) error {
	operators, err := operatorLoader.Apply()
	if err != nil {
		return fmt.Errorf("failed to load operator configurations: %v", err)
	}

	for _, operator := range operators {
		for _, version := range operator.Versions {
			log.WithField("operator", operator.Name).
				WithField("version", version.Version()).
				Info("Validating operator")

			if err := validateOperator(ctx, operator, version, strict); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateOperator(
	ctx context.Context,
	operator operator.Operator,
	version operator.Version,
	strict bool,
) (err error) {
	operatorName := fmt.Sprintf("%s-%s", operator.Name, version.Version())

	r, err := resolver.New(operator, version)
	if err != nil {
		return fmt.Errorf("failed to resolve operator %q: %v", operatorName, err)
	}

	pkgFs, remover, err := r.Resolve(ctx)
	if err != nil {
		return fmt.Errorf("failed to resolve operator %q: %v", operatorName, err)
	}

	// The package resolver created a temporary directory for the package file system.
	// We remove it once we no longer need it.
	defer func() {
		if rerr := remover(); rerr != nil {
			err = fmt.Errorf("failed to remove temporary directory of operator %q: %v", operatorName, rerr)
		}
	}()

	pkg, err := repo.NewPackage(pkgFs)
	if err != nil {
		return fmt.Errorf("failed to extract package version of operator %q: %v", operatorName, err)
	}

	validationResult := validation.Validate(operator, version, pkg)

	var warnings, errors string

	if strict {
		errors = strings.Join(append(validationResult.Warnings, validationResult.Errors...), "\n")
	} else {
		warnings = strings.Join(validationResult.Warnings, "\n")
		errors = strings.Join(validationResult.Errors, "\n")
	}

	if warnings != "" {
		fmt.Printf("validation warnings for operator %q:\n%s", operatorName, warnings)
	}

	if errors != "" {
		return fmt.Errorf("validation failed for operator %q:\n%s", operatorName, errors)
	}

	return nil
}
