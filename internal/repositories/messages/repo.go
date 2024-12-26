package messagesrepo

import (
	"github.com/Slava02/ChatSupport/internal/store"
)

//go:generate options-gen -out-filename=repo_options.gen.go -from-struct=Options
type Options struct {
	db *store.Database `option:"mandatory" validate:"required"`
}

type Repo struct {
	Options
}

func New(opts Options) (*Repo, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return &Repo{Options: opts}, nil
}
