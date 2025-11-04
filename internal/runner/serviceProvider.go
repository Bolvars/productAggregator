package runner

import (
	"productsParser/internal/config"
	"productsParser/internal/infrastructure/parsers/tilda"
	"productsParser/internal/service"
)

type serviceProvider struct {
	config config.TgConfig
	parser func() service.Parser
}

func newServiceProvider(config config.TgConfig) *serviceProvider {
	return &serviceProvider{
		config: config,
	}
}

func (s *serviceProvider) Config() config.TgConfig {
	return s.config
}

func (s *serviceProvider) ParserInit() func() service.Parser {
	if s.parser == nil {
		s.parser = func() service.Parser {
			return tilda.New()
		}
	}

	return s.parser
}
