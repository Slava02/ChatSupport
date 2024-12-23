// Code generated by options-gen. DO NOT EDIT.
package keycloakclient

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	basePath string,
	keyCloakRealm string,
	keyCloakClientID string,
	keyCloakClientSecret string,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.basePath = basePath

	o.keyCloakRealm = keyCloakRealm

	o.keyCloakClientID = keyCloakClientID

	o.keyCloakClientSecret = keyCloakClientSecret

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func WithDebugMode(opt bool) OptOptionsSetter {
	return func(o *Options) {
		o.debugMode = opt

	}
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("basePath", _validate_Options_basePath(o)))
	errs.Add(errors461e464ebed9.NewValidationError("keyCloakRealm", _validate_Options_keyCloakRealm(o)))
	errs.Add(errors461e464ebed9.NewValidationError("keyCloakClientID", _validate_Options_keyCloakClientID(o)))
	errs.Add(errors461e464ebed9.NewValidationError("keyCloakClientSecret", _validate_Options_keyCloakClientSecret(o)))
	errs.Add(errors461e464ebed9.NewValidationError("debugMode", _validate_Options_debugMode(o)))
	return errs.AsError()
}

func _validate_Options_basePath(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.basePath, "required,url"); err != nil {
		return fmt461e464ebed9.Errorf("field `basePath` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_keyCloakRealm(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.keyCloakRealm, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `keyCloakRealm` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_keyCloakClientID(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.keyCloakClientID, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `keyCloakClientID` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_keyCloakClientSecret(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.keyCloakClientSecret, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `keyCloakClientSecret` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_debugMode(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.debugMode, "omitempty"); err != nil {
		return fmt461e464ebed9.Errorf("field `debugMode` did not pass the test: %w", err)
	}
	return nil
}
