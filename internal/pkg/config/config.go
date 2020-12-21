package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"gopkg.in/square/go-jose.v2/json"
)

type Config struct {
	Port        string   `env:"PORT" envDefault:"9090"`
	Certificate string   `env:"CERTIFICATE,file" envDefault:"../../certificates/auth0.pem" json:"-"`
	Issuer      string   `env:"ISSUER" envDefault:"https://dev-dgoly5h6.eu.auth0.com/"`
	Audience    []string `env:"AUDIENCE" envDefault:"casseur_flutter"`
	ApiKey      string   `env:"API_KEY,required"`
	BaseID      string   `env:"BASE_ID,required"`
}

func (c *Config) String() string {
	payload, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(payload)
}

func NewConfig() (*Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, err
	}
	// Heroku override the variable and provide an invalid variable
	c.Port = fmt.Sprintf(":%s", c.Port)
	return &c, nil
}
