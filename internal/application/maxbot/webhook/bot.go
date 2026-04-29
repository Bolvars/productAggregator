package webhook

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"productsParser/internal/application/maxbot/common"
	"productsParser/internal/config"
	service "productsParser/internal/service/users"
	"time"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type Bot struct {
	*common.CommonBot
	server *fasthttp.Server
	secret string
	host   string
}

func New(config config.Config, uService *service.UserService) (*Bot, error) {
	cmnBot, err := common.NewBot(config.Token(), uService)
	if err != nil {
		return nil, err
	}

	return &Bot{
		CommonBot: cmnBot,
		host:      config.Host(),   // Например: https://mybot.com
		secret:    config.Secret(), // Для проверки подлинности запросов
	}, nil
}

func (bt *Bot) Start(ctx context.Context) {

	subs, _ := bt.Api().Subscriptions.GetSubscriptions(ctx)
	for _, s := range subs.Subscriptions {
		_, _ = bt.Api().Subscriptions.Unsubscribe(ctx, s.Url)
	}
	updateChan := make(chan schemes.UpdateInterface, 100)

	webhookURL := bt.host + common.WebhookHandle
	_, err := bt.Api().Subscriptions.Subscribe(ctx, webhookURL, []string{}, bt.secret)
	if err != nil {
		log.Fatalf("Failed to subscribe webhook: %v", err)
	}
	log.Printf("Webhook registered: %s", webhookURL)

	fastHandler := fasthttpadaptor.NewFastHTTPHandlerFunc(bt.Api().GetUpdateHandler(updateChan, bt.secret))

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/webhook":
			fastHandler(ctx)
		default:
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	}

	bt.server = &fasthttp.Server{
		Handler: requestHandler,
	}

	log.Println("Starting Webhook server on :10888")
	go func() {
		if err := bt.server.ListenAndServe(":10888"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			bt.stop()
			return
		case upd := <-updateChan:
			bt.Handle(ctx, upd)
		}
	}

}

func requestHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world! Requested path is %q", ctx.Path())
}

func (bt *Bot) stop() {
	if bt.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		bt.server.ShutdownWithContext(ctx)
		log.Println("Webhook server stopped")
	}
}
