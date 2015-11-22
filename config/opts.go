package config

type ConfigOpts struct {
	File string `flag:"config-file" env:"CONFIG_FILE" default:"config.toml"`
}