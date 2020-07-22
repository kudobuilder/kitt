package validation

import (
	"fmt"
)

// Result is a validation result consisting of a list of warnings and errors.
type Result struct {
	Warnings []string
	Errors   []string
}

// AddWarning adds a warning to a validation result.
func (r *Result) AddWarning(warning string) {
	r.Warnings = append(r.Warnings, warning)
}

// AddWarningf adds a formatted warning to a validation result.
func (r *Result) AddWarningf(warning string, a ...interface{}) {
	r.Warnings = append(r.Warnings, fmt.Sprintf(warning, a...))
}

// AddError adds an error to a validation result.
func (r *Result) AddError(error string) {
	r.Errors = append(r.Errors, error)
}

// AddErrorf adds a formatted error to a validation result.
func (r *Result) AddErrorf(error string, a ...interface{}) {
	r.Errors = append(r.Errors, fmt.Sprintf(error, a...))
}
