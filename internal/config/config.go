package config

import (
	"time"
)

type Config interface {
	TgBotEnabled() bool
	Token() string
	Secret() string
	IsWebhook() bool
	Host() string
	InterruptTimeout() time.Duration
	IsTildaParser() bool
	DefaultURL() string
}

type GlobalConfig struct {
	tgBotEnabled       bool
	isWebhook          bool
	token              string
	secret             string
	host               string
	defaultURL         string
	interruptTimeout   time.Duration
	tildaParserEnabled bool
}

func NewGlobalConfig(
	tgEnabled, isWebhook bool,
	token, secret, host, defUrl string,
	timeout time.Duration,
	tildaParser bool,
) Config {
	return &GlobalConfig{
		tgBotEnabled:       tgEnabled,
		isWebhook:          isWebhook,
		secret:             secret,
		token:              token,
		host:               host,
		defaultURL:         defUrl,
		interruptTimeout:   timeout,
		tildaParserEnabled: tildaParser,
	}
}

func (c *GlobalConfig) TgBotEnabled() bool {
	return c.tgBotEnabled
}

func (c *GlobalConfig) Token() string {
	return c.token
}

func (c *GlobalConfig) InterruptTimeout() time.Duration {
	return c.interruptTimeout
}

func (c *GlobalConfig) IsTildaParser() bool {
	return c.tildaParserEnabled
}

func (c *GlobalConfig) IsWebhook() bool {
	return c.isWebhook
}
func (c *GlobalConfig) Secret() string {
	return c.secret
}

func (c *GlobalConfig) Host() string {
	return c.host
}

func (c *GlobalConfig) DefaultURL() string {
	return c.defaultURL
}
