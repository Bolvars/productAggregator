package telegrambot

import (
	"context"
	"productsParser/internal/application/common"
	service "productsParser/internal/service/users"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Bot struct {
	api *bot.Bot
	cmn *common.Common
}

func NewBot(token string, uService *service.UserService) (*bot.Bot, error) {
	bt := &Bot{
		cmn: common.New(uService),
	}
	b, err := bot.New(token,
		bot.WithMessageTextHandler(common.HandleNameHelp, bot.MatchTypeExact, bt.handleHelpCommand),
		bot.WithMessageTextHandler(common.HandleNameStart, bot.MatchTypeExact, bt.handleStartCommand),
		bot.WithMessageTextHandler(common.HandleNameCalc, bot.MatchTypeExact, bt.handleCalculate),
		bot.WithMessageTextHandler(common.HandleNameCalcAndSort, bot.MatchTypeExact, bt.handleCalculateAndSort),
		bot.WithDefaultHandler(bt.handleTextMessage),
	)
	if err != nil {
		return nil, err
	}
	bt.api = b
	return bt.api, nil
}

func (bt *Bot) handleHelpCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   common.Help,
	})
}

func (bt *Bot) handleStartCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   common.AboutBot,
	})
}

func (bt *Bot) handleTextMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.Chat.ID
	ok, txtResp, err := bt.cmn.AddOrder(ctx, strconv.FormatInt(userID, 10), update.Message.Text)

	message := &bot.SendMessageParams{
		ChatID: userID,
	}
	if err != nil {
		message.Text = err.Error()
		b.SendMessage(ctx, message)
		return
	}
	var menu models.ReplyMarkup
	if !ok {
		menu = &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{{Text: common.HandleNameCalc}},
				{{Text: common.HandleNameCalcAndSort}},
			},
			ResizeKeyboard: true, // чтобы не занимала полэкрана
		}
	}
	message.ReplyMarkup = menu
	message.Text = txtResp
	b.SendMessage(ctx, message)
}

func (bt *Bot) handleCalculate(ctx context.Context, b *bot.Bot, update *models.Update) {
	bt.calc(ctx, b, update, false)
}

func (bt *Bot) handleCalculateAndSort(ctx context.Context, b *bot.Bot, update *models.Update) {
	bt.calc(ctx, b, update, true)
}

func (bt *Bot) calc(ctx context.Context, b *bot.Bot, update *models.Update, sorted bool) {

	userID := update.Message.From.ID
	res, err := bt.cmn.Calculate(ctx, strconv.FormatInt(userID, 10), sorted)

	menu := &models.ReplyKeyboardRemove{
		RemoveKeyboard: true,
	}
	message := &bot.SendMessageParams{
		ChatID:      userID,
		ReplyMarkup: menu,
	}
	if err != nil {
		message.Text = err.Error()
		b.SendMessage(ctx, message)
		return
	}

	message.Text = common.Results + res.String()
	b.SendMessage(ctx, message)
}
