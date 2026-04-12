package common

import (
	"context"
	"errors"
	"fmt"
	service "productsParser/internal/service/users"
	"strings"
)

type Common struct {
	userService *service.UserService
}

func New(us *service.UserService) *Common {
	return &Common{
		userService: us,
	}
}

func (c *Common) Calculate(ctx context.Context, userID string, sorted bool) (*strings.Builder, error) {

	curUser, exist := c.userService.GetUserService(userID)
	if !exist {
		return nil, errors.New("Пользователь не добавил заказы для расчета")
	}

	defer c.userService.DelUser(userID)
	productService, err := curUser.Compute()
	if err != nil {
		return nil, errors.New("Не удалось произвести расчёт: " + err.Error())
	}

	builder := &strings.Builder{}
	for i, product := range productService.Products(sorted) {
		builder.WriteString(fmt.Sprintf("%d. ", i+1) + product.Name() + ": ")
		builder.WriteString(product.ToString() + "\n")
	}
	return builder, nil
}

func (c *Common) AddOrder(ctx context.Context, userId, text string) (bool, string, error) {
	curUser, ok := c.userService.GetOrAddUserService(userId)
	order, err := curUser.AddOrder([]byte(text))

	if err != nil {
		return false, "", errors.New("Не удалось распарсить текст: " + err.Error())
	}
	return ok, fmt.Sprintf(CorrectResponse, order.Id()), nil
}
