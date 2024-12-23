// Code generated by options-gen. DO NOT EDIT.
package serverclient

import (
	fmt461e464ebed9 "fmt"

	keycloakclient "github.com/Slava02/ChatSupport/internal/clients/keycloak"
	clientv1 "github.com/Slava02/ChatSupport/internal/server-client/v1"
	"github.com/getkin/kin-openapi/openapi3"
	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
	"go.uber.org/zap"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	logger *zap.Logger,
	addr string,
	allowOrigins []string,
	v1Swagger *openapi3.T,
	v1Handlers clientv1.ServerInterface,
	keycloak *keycloakclient.Client,
	resource string,
	role string,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.logger = logger

	o.addr = addr

	o.allowOrigins = allowOrigins

	o.v1Swagger = v1Swagger

	o.v1Handlers = v1Handlers

	o.keycloak = keycloak

	o.resource = resource

	o.role = role

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("logger", _validate_Options_logger(o)))
	errs.Add(errors461e464ebed9.NewValidationError("addr", _validate_Options_addr(o)))
	errs.Add(errors461e464ebed9.NewValidationError("allowOrigins", _validate_Options_allowOrigins(o)))
	errs.Add(errors461e464ebed9.NewValidationError("v1Swagger", _validate_Options_v1Swagger(o)))
	errs.Add(errors461e464ebed9.NewValidationError("v1Handlers", _validate_Options_v1Handlers(o)))
	errs.Add(errors461e464ebed9.NewValidationError("keycloak", _validate_Options_keycloak(o)))
	errs.Add(errors461e464ebed9.NewValidationError("resource", _validate_Options_resource(o)))
	errs.Add(errors461e464ebed9.NewValidationError("role", _validate_Options_role(o)))
	return errs.AsError()
}

func _validate_Options_logger(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.logger, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `logger` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_addr(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.addr, "required,hostname_port"); err != nil {
		return fmt461e464ebed9.Errorf("field `addr` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_allowOrigins(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.allowOrigins, "min=1"); err != nil {
		return fmt461e464ebed9.Errorf("field `allowOrigins` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_v1Swagger(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.v1Swagger, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `v1Swagger` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_v1Handlers(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.v1Handlers, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `v1Handlers` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_keycloak(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.keycloak, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `keycloak` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_resource(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.resource, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `resource` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_role(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.role, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `role` did not pass the test: %w", err)
	}
	return nil
}
