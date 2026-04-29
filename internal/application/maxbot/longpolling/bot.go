package longpolling

import (
	"context"
	"productsParser/internal/application/maxbot/common"
	service "productsParser/internal/service/users"
)

type Bot struct {
	*common.CommonBot
}

func New(token string, uService *service.UserService) (*Bot, error) {
	cmnBot, err := common.NewBot(token, uService)
	if err != nil {
		return nil, err
	}
	bt := &Bot{
		CommonBot: cmnBot,
	}
	return bt, nil
}

func (bt *Bot) Start(ctx context.Context) {
	for upd := range bt.Api().GetUpdates(ctx) {
		bt.Handle(ctx, upd)
	}
}
