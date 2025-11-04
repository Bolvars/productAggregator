package config

import "time"

type TgConfig interface {
	TgBotEnabled() bool
	Token() string
	InterruptTimeout() time.Duration
	IsTildaParser() bool
}

type tgConfig struct {
	tgBotEnabled       bool
	token              string
	interruptTimeout   time.Duration
	tildaParserEnabled bool
}

func NewGlobalConfig(
	tgEnabled bool,
	token string,
	timeout time.Duration,
	tildaParser bool,
) TgConfig {
	return &tgConfig{
		tgBotEnabled:       tgEnabled,
		token:              token,
		interruptTimeout:   timeout,
		tildaParserEnabled: tildaParser,
	}
}

func (c *tgConfig) TgBotEnabled() bool {
	return c.tgBotEnabled
}

func (c *tgConfig) Token() string {
	return c.token
}

func (c *tgConfig) InterruptTimeout() time.Duration {
	return c.interruptTimeout
}

func (c *tgConfig) IsTildaParser() bool {
	return c.tildaParserEnabled
}
