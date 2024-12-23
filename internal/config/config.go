package config

type Config struct {
	Global  GlobalConfig  `toml:"global"`
	Log     LogConfig     `toml:"log"`
	Servers ServersConfig `toml:"servers"`
	Sentry  SentryConfig  `toml:"sentry"`
	Clients Clients       `toml:"clients"`
}

type GlobalConfig struct {
	Env string `toml:"env" validate:"required,oneof=dev stage prod"`
}

type LogConfig struct {
	Level string `toml:"level" validate:"required,oneof=debug info warn error"`
}

type ServersConfig struct {
	Debug  DebugServerConfig  `toml:"debug"`
	Client ClientServerConfig `toml:"client"`
}

type DebugServerConfig struct {
	Addr string `toml:"addr" validate:"required,hostname_port"`
}

type ClientServerConfig struct {
	Addr         string         `toml:"addr" validate:"required,hostname_port"`
	AllowOrigins []string       `toml:"allow_origins" validate:"required,dive,url"`
	Access       RequiredAccess `toml:"required_access" validate:"required"`
}

type RequiredAccess struct {
	Resource string `toml:"resource" validate:"required"`
	Role     string `toml:"role" validate:"required"`
}

type SentryConfig struct {
	DSN string `toml:"dsn" validate:"omitempty,url"`
}

type Clients struct {
	Keycloak `toml:"keycloak"`
}

type Keycloak struct {
	BasePath     string `toml:"base_path" validate:"required,url"`
	Realm        string `tom:"realm" validate:"required"`
	ClientID     string `toml:"client_id" validate:"required"`
	ClientSecret string `toml:"client_secret" validate:"required"`
	DebugMode    bool   `toml:"debug_mode" default:"false" validate:"omitempty"`
}
