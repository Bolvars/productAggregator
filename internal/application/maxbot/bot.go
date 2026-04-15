package maxbot

import (
	"context"
	"log"
	"productsParser/internal/application/common"
	service "productsParser/internal/service/users"
	"strconv"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type Bot struct {
	api      *maxbot.Api
	commands map[string]func(context.Context, *schemes.MessageCreatedUpdate) error
	cmn      *common.Common
}

func NewBot(token string, uService *service.UserService) (*Bot, error) {
	bt := &Bot{
		cmn: common.New(uService),
	}
	api, err := maxbot.New(token)
	if err != nil {
		return nil, err
	}

	commands := map[string]func(context.Context, *schemes.MessageCreatedUpdate) error{
		common.HandleNameHelp:        bt.handleHelpCommand,
		common.HandleNameStart:       bt.handleStartCommand,
		common.HandleNameCalc:        bt.handleCalculate,
		common.HandleNameCalcAndSort: bt.handleCalculateAndSort,
	}

	bt.api = api
	bt.commands = commands
	return bt, nil
}

func (bt *Bot) Start(ctx context.Context) {
	for upd := range bt.api.GetUpdates(ctx) {
		switch upd := upd.(type) {
		case *schemes.BotStartedUpdate: // Определение типа пришедшего обновления
			bt.api.Messages.Send(ctx, maxbot.NewMessage().SetChat(upd.ChatId).SetText(common.AboutBot))
		case *schemes.MessageCreatedUpdate:
			handleFunc, ok := bt.commands[upd.GetCommand()]
			if !ok {
				if err := bt.handleTextMessage(ctx, upd); err != nil {
					log.Println("error", err.Error())
				}
			} else {
				if err := handleFunc(ctx, upd); err != nil {
					log.Println("error", err.Error())
				}
			}
		}
	}

}

func (bt *Bot) handleHelpCommand(ctx context.Context, msg *schemes.MessageCreatedUpdate) error {
	return bt.api.Messages.Send(ctx, maxbot.NewMessage().SetChat(msg.GetChatID()).SetText(common.Help))
}

func (bt *Bot) handleStartCommand(ctx context.Context, msg *schemes.MessageCreatedUpdate) error {
	return bt.api.Messages.Send(ctx, maxbot.NewMessage().SetChat(msg.GetChatID()).SetText(common.AboutBot))
}

func (bt *Bot) handleTextMessage(ctx context.Context, msg *schemes.MessageCreatedUpdate) error {
	chatId := msg.GetChatID()

	var txt string
	if msg.Message.Link != nil {
		txt = msg.Message.Link.Message.Text
	} else {
		txt = msg.GetText()
	}
	_, txtResp, err := bt.cmn.AddOrder(ctx, strconv.FormatInt(msg.GetUserID(), 10), txt)
	if err != nil {
		return bt.api.Messages.Send(ctx, maxbot.NewMessage().SetChat(chatId).SetText(err.Error()))
	}
	var message *maxbot.Message

	keyboard := bt.api.Messages.NewKeyboardBuilder()
	keyboard.AddRow().
		AddMessage(common.HandleNameCalc).
		AddMessage(common.HandleNameCalcAndSort)
	message = maxbot.NewMessage().
		SetChat(chatId).
		AddKeyboard(keyboard).
		SetText(txtResp)

	return bt.api.Messages.Send(ctx, message)
}

func (bt *Bot) handleCalculate(ctx context.Context, msg *schemes.MessageCreatedUpdate) error {
	return bt.calc(ctx, msg, false)
}

func (bt *Bot) handleCalculateAndSort(ctx context.Context, msg *schemes.MessageCreatedUpdate) error {
	return bt.calc(ctx, msg, true)
}

func (bt *Bot) calc(ctx context.Context, msg *schemes.MessageCreatedUpdate, sorted bool) error {
	chatId := msg.GetChatID()
	res, err := bt.cmn.Calculate(ctx, strconv.FormatInt(msg.GetUserID(), 10), sorted)

	var message *maxbot.Message
	if err != nil {
		message = maxbot.NewMessage().SetChat(chatId).SetText(err.Error())
		return bt.api.Messages.Send(ctx, message)
	}

	message = maxbot.NewMessage().SetChat(chatId).SetText("Результаты расчётов:\n" + res.String())
	return bt.api.Messages.Send(ctx, message)
}
