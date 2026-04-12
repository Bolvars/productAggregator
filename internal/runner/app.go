package runner

import (
	"context"
	"errors"
	"fmt"
	"log"
	"productsParser/internal/application/maxbot"
	"productsParser/internal/application/telegrambot"
	"productsParser/internal/config"
	service "productsParser/internal/service/users"

	"golang.org/x/sync/errgroup"
)

type Bot interface {
	Start(ctx context.Context)
}

// Структура, реализующая запуск бинаринка
type GatewayApp struct {
	serviceProvider *serviceProvider
	config          config.TgConfig

	b Bot

	g   *errgroup.Group
	ctx context.Context
}

func NewApp(ctx context.Context, config config.TgConfig) (*GatewayApp, error) {
	app := &GatewayApp{
		config: config,
	}

	app.g, app.ctx = errgroup.WithContext(ctx)

	var err error

	err = app.initDeps(app.ctx)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (a *GatewayApp) Ctx() context.Context {
	return a.ctx
}

func (a *GatewayApp) Run(ctx context.Context) error {
	a.runServers(ctx)
	if err := a.g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("server exited with error: %w", err)
	}
	log.Println("Shutdown done")
	return nil
}

func (a *GatewayApp) initDeps(ctx context.Context) error {
	inits := []func(ctx2 context.Context) error{
		a.initServiceProvider,
		a.initBot,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *GatewayApp) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider(a.config)
	return nil
}

func (a *GatewayApp) initBot(_ context.Context) error {
	userService := service.NewUserService(a.serviceProvider.ParserInit())
	var (
		b   Bot
		err error
	)
	if a.config.TgBotEnabled() {
		b, err = telegrambot.NewBot(a.serviceProvider.Config().Token(), userService)
		if err != nil {
			return err
		}
	} else {
		b, err = maxbot.NewBot(a.serviceProvider.Config().Token(), userService)
		if err != nil {
			return err
		}
	}

	a.b = b
	return nil
}

func (a *GatewayApp) runServers(ctx context.Context) {
	a.g.Go(func() error {
		if a.b != nil {
			log.Println("Bot was started")
			a.b.Start(ctx)
		}
		return nil
	})
}
