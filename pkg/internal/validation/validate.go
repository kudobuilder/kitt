package validation

import (
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/kudobuilder/kudo/pkg/kudoctl/cmd/verify"
	"github.com/kudobuilder/kudo/pkg/kudoctl/packages/reader"

	"github.com/kudobuilder/kitt/pkg/internal/apis/operator"
	"github.com/kudobuilder/kitt/pkg/internal/repo"
)

// Validate runs several checks on an operator and it's referenced package.
func Validate(operator operator.Operator, version operator.Version, pkg repo.Package) Result {
	result := Result{}

	validateVersion(version, pkg, &result)
	validateVerify(pkg, &result)

	return result
}

func validateVersion(version operator.Version, pkg repo.Package, result *Result) {
	operatorVersion, err := semver.NewVersion(version.OperatorVersion)
	if err != nil {
		result.AddError("operatorVersion isn't semver")
	} else if !operatorVersion.Equal(&pkg.OperatorVersion) {
		result.AddWarningf(
			"operatorVersion %q doesn't match operatorVersion %q in operator package",
			operatorVersion,
			pkg.OperatorVersion)
	}

	if version.AppVersion != "" {
		appVersion, err := semver.NewVersion(version.AppVersion)
		if err != nil {
			result.AddError("appVersion isn't semver")
		} else {
			if pkg.AppVersion == nil {
				result.AddWarning("appVersion provided but not set in operator package")
			} else if !appVersion.Equal(pkg.AppVersion) {
				result.AddWarningf(
					"appVersion %q doesn't match appVersion %q in operator package",
					appVersion,
					pkg.AppVersion)
			}
		}
	} else if pkg.AppVersion != nil {
		result.AddWarning("appVersion not provided but set in operator package")
	}
}

func validateVerify(pkg repo.Package, result *Result) {
	p, err := reader.ReadDir(pkg, string(filepath.Separator))
	if err != nil {
		// 'repo.Package' has been created by 'reader.ReadDir'.
		// We don't expect it to fail when running 'reader.ReadDir' again.
		panic(err)
	}

	verifyResult := verify.PackageFiles(p.Files)

	for _, warning := range verifyResult.Warnings {
		result.AddWarning(warning)
	}

	for _, error := range verifyResult.Errors {
		result.AddError(error)
	}
}
