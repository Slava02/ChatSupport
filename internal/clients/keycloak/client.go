package keycloakclient

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

//go:generate options-gen -out-filename=client_options.gen.go -from-struct=Options
type Options struct {
	basePath             string `option:"mandatory" validate:"required,url"`
	keyCloakRealm        string `option:"mandatory" validate:"required"`
	keyCloakClientID     string `option:"mandatory" validate:"required"`
	keyCloakClientSecret string `option:"mandatory" validate:"required"`
	debugMode            bool   `validate:"omitempty"`
}

// Client is a tiny client to the KeyCloak realm operations. UMA configuration:
// http://localhost:3010/realms/Bank/.well-known/uma2-configuration
type Client struct {
	keyCloakRealm        string `validate:"required"`
	keyCloakClientID     string `validate:"required"`
	keyCloakClientSecret string `validate:"required"`
	cli                  *resty.Client
}

func New(opts Options) (*Client, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}

	cli := resty.New()
	cli.SetDebug(opts.debugMode)
	cli.SetBaseURL(opts.basePath)

	return &Client{
		keyCloakRealm:        opts.keyCloakRealm,
		keyCloakClientID:     opts.keyCloakClientID,
		keyCloakClientSecret: opts.keyCloakClientSecret,
		cli:                  cli,
	}, nil
}
