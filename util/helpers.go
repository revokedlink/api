package util

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
)

// RestrictFields prevents specific fields from being included in the request body.
// If any of the provided fields are found in the request, a 400 Bad Request
// with detailed validation errors is returned. This is already possible by `isset`,
// but doesn't return a user understandable error code
//
// Usage:
// hooks.RestrictFields(e, util.Fields.User.Email, util.Fields.User.Verified)
func RestrictFields(e *core.RecordRequestEvent, fields ...string) error {
	info, err := e.RequestInfo()
	if err != nil {
		return nil // Fallback to next handler if request info is missing
	}

	if info.HasSuperuserAuth() {
		return nil
	}

	errs := validation.Errors{}
	for _, field := range fields {
		if _, ok := info.Body[field]; ok {
			errs[field] = validation.NewError(Errors.ValidationFieldRestricted.ErrorCode, Errors.ValidationFieldRestricted.ErrorText)
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// RequireFields ensures that specific fields are present in the request body.
// This is already possible with the `Required` annotation
func RequireFields(e *core.RecordRequestEvent, fields ...string) error {
	info, err := e.RequestInfo()
	if err != nil {
		return nil
	}

	errs := validation.Errors{}
	for _, field := range fields {
		if val, ok := info.Body[field]; !ok || val == nil || val == "" {
			errs[field] = validation.NewError(Errors.ValidationFieldRequired.ErrorCode, Errors.ValidationFieldRequired.ErrorText)
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
