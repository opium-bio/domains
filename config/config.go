package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	App        App        `toml:"app"`
	Statuses   []string   `toml:"statuses"`
	JWT        JWT        `toml:"JWT"`
	Discord    Discord    `toml:"discord"`
	Cloudflare Cloudflare `toml:"cloudflare"`
	Redis      Redis      `toml:"redis"`
	MongoDB    MongoDB    `toml:"mongodb"`
}
type App struct {
	Port int `toml:"port"`
}
type JWT struct {
	Secret string `toml:"secret"`
}
type Discord struct {
	Username string `toml:"username"`
	Webhook  string `toml:"webhook"`
}
type Cloudflare struct {
	CFApiKey  string `toml:"cf_apikey"`
	CFAccount string `toml:"account_id"`
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
