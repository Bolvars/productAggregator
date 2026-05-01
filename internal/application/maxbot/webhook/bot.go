package webhook

import (
	"context"
	"errors"
	"log"
	"net/http"
	"productsParser/internal/application/maxbot/common"
	"productsParser/internal/config"
	service "productsParser/internal/service/users"
	"time"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type Bot struct {
	*common.CommonBot
	server *http.Server
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
		host:      config.Host(),
		secret:    config.Secret(),
	}, nil
}

func (bt *Bot) Start(ctx context.Context) {

	subs, _ := bt.Api().Subscriptions.GetSubscriptions(ctx)
	for _, s := range subs.Subscriptions {
		_, _ = bt.Api().Subscriptions.Unsubscribe(ctx, s.Url)
	}
	webhookURL := bt.host + common.WebhookHandle
	_, err := bt.Api().Subscriptions.Subscribe(ctx, webhookURL, []string{}, bt.secret)
	if err != nil {
		log.Fatalf("Failed to subscribe webhook: %v", err)
	}
	log.Printf("Webhook registered: %s", webhookURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", bt.Api().GetUpdateHandlerFunc(func(ui schemes.UpdateInterface) {
		bt.Handle(ctx, ui)
	}, bt.secret))

	bt.server = &http.Server{
		Addr:    ":10888",
		Handler: mux,
	}

	log.Println("Starting Webhook server on :10888")
	go func() {
		if err := bt.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	bt.stop()
}

func (bt *Bot) stop() {
	if bt.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		bt.server.Shutdown(ctx)
		log.Println("Webhook server stopped")
	}
}
