package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Slava02/ChatSupport/internal/config"
)

func TestGlobalConfig_IsProduction(t *testing.T) {
	assert.True(t, config.GlobalConfig{Env: "prod"}.IsProduction())
	assert.False(t, config.GlobalConfig{Env: "neprod"}.IsProduction())
}
