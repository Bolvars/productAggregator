package telegrambot

import (
	"context"
	"fmt"
	"productsParser/internal/service"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const aboutBot = `👋 Привет! Я бот для обработки заказов. Пока что я умею обрабатывать только текст заказа из Tilda.

📦 Отправь сюда текст заказа — я распознаю товары, вес и количество.

После этого ты сможешь нажать кнопку «Рассчитать» (или ввести /calc),
и я покажу агрегированные результаты по заказам.

Попробуй — просто вставь сообщение из Тильды.`

const help = `ℹ️ Помощь

/start — описание бота
/calc — рассчитать заказы
/help — показать эту помощь

Пример сообщения tilda, который поддерживает бот:

Order #12345
1. Салат , pc: 3600 (50 pc x 72)
2. Апельсин: 585 (3 x 195) Вес: 1000 гр.
3. Лайм: 1350 (3 x 450) Вес: 1 гр.
4. Лимон: 460 (2 x 230) Вес: 1 шт.
`

type Bot struct {
	api         *bot.Bot
	userService *service.UserService
}

func NewBot(token string, uService *service.UserService) (*bot.Bot, error) {
	bt := &Bot{
		userService: uService,
	}
	b, err := bot.New(token,
		bot.WithMessageTextHandler("/help", bot.MatchTypeExact, bt.handleHelpCommand),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, bt.handleStartCommand),
		bot.WithMessageTextHandler("/calc", bot.MatchTypeExact, bt.handleCalculate),
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
		Text:   help,
	})
}

func (bt *Bot) handleStartCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   aboutBot,
	})
}

func (bt *Bot) handleTextMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.Chat.ID
	id := strconv.FormatInt(userID, 10)
	curUser, ok := bt.userService.GetOrAddUserService(id)
	order, err := curUser.AddOrder([]byte(update.Message.Text))

	message := &bot.SendMessageParams{
		ChatID: userID,
	}
	if err != nil {
		message.Text = "Не удалось распарсить текст: " + err.Error()
		b.SendMessage(ctx, message)
		return
	}
	var menu models.ReplyMarkup
	if !ok {
		menu = &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{{Text: "/calc"}},
			},
			ResizeKeyboard: true, // чтобы не занимала полэкрана
		}
	}
	message.ReplyMarkup = menu
	message.Text = fmt.Sprintf("Заказ #%s успешно сохранен. Чтобы выполнить расчёт, нажмите «Рассчитать».", order.Id())
	b.SendMessage(ctx, message)
}

func (bt *Bot) handleCalculate(ctx context.Context, b *bot.Bot, update *models.Update) {

	userID := update.Message.From.ID
	id := strconv.FormatInt(userID, 10)
	curUser, exist := bt.userService.GetUserService(id)

	menu := &models.ReplyKeyboardRemove{
		RemoveKeyboard: true,
	}
	message := &bot.SendMessageParams{
		ChatID:      userID,
		ReplyMarkup: menu,
	}
	if !exist {
		message.Text = "Пользователь не добавил заказы для расчета"
		b.SendMessage(ctx, message)
		return
	}

	defer bt.userService.DelUser(id)
	products, err := curUser.Compute()
	if err != nil {
		message.Text = "Не удалось произвести расчёт: " + err.Error()
		b.SendMessage(ctx, message)
		return
	}

	builder := strings.Builder{}
	for i, product := range products {
		builder.WriteString(fmt.Sprintf("%d. ", i+1) + product.Name() + ": ")
		builder.WriteString(product.ToString() + "\n")
	}
	message.Text = "Результаты расчётов:\n" + builder.String()
	b.SendMessage(ctx, message)

}
