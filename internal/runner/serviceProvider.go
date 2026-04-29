package runner

import (
	"productsParser/internal/config"
	i "productsParser/internal/domain/interface"
	"productsParser/internal/infrastructure/parsers/tilda"
)

type serviceProvider struct {
	config config.Config
	parser func() i.Parser
}

func newServiceProvider(config config.Config) *serviceProvider {
	return &serviceProvider{
		config: config,
	}
}

func (s *serviceProvider) Config() config.Config {
	return s.config
}

func (s *serviceProvider) ParserInit() func() i.Parser {
	if s.parser == nil {
		s.parser = func() i.Parser {
			return tilda.New()
		}
	}

	return s.parser
}
