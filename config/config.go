package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	App        App        `toml:"app"`
	Statuses   []string   `toml:"statuses"`
	JWT        JWT        `toml:"JWT"`
	Cloudflare Cloudflare `toml:"cloudflare"`
	Redis      Redis      `toml:"redis"`
	MongoDB    MongoDB    `toml:"mongodb"`
}
type App struct {
	Port int `toml:"port"`
}
type JWT struct {
	secret string `toml:"secret"`
}
type Cloudflare struct {
	CFApiKey string `toml:"cf_apikey"`
}
type Redis struct {
	Hostname string `toml:"hostname"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
}
type MongoDB struct {
	String string `toml:"string"`
}

func LoadConfig(filepath string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(filepath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
