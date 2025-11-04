package tilda_test

import (
	"log"
	"productsParser/internal/domain"
	"productsParser/internal/infrastructure/parsers/tilda"
	"testing"

	"github.com/stretchr/testify/assert"
)

const messageTest1 = `Order #1898307952
1. Огурец среднеплодный: 749 (7 x 107) Вес: 1000 гр.
2. Томат красный: 135 (1 x 135) Вес: 1000 гр.
3. Перец чили красный: 35 (100 x 0.35) Вес: 1 гр.
4. Салат Афицион, pc: 58 (1 pc x 58)
5. Авокадо: 558 (1 x 558) Вес: 1000 гр.
Payment Amount: 1535 RUB
Payment system: (none)
Purchaser information:
name: Шаркова М.Б.
Адрес_доставки_: проспект ленина 126, ТЦ. Кондор, цокольный этаж, вход со стороны ленина,
Checkbox: yes
ma_name: ИП Шаркова М.Б. «Суши Кит»
ma_email: sushi-kit-tsk@yandex.ru
ma_phone: +79521629303
ma_id: 47259433
Additional information:
Transaction ID: 11091747:7807579192
Block ID: rec820069702
Form Name: Cart
https://bubon-horeca.ru/shop-spices
`

const messageTest2 = `Order #1024598132
1. Баклажан: 294 (3 x 98) Вес: 1000 гр.
2. Петрушка: 0.9 (3 x 0.3) Вес: 1 гр.
3. Картофель желтый: 160 (5 x 32) Вес: 1000 гр.
4. Морковь: 144 (4 x 36) Вес: 1000 гр.
5. Лук репчатый новый урожай: 90 (3 x 30) Вес: 1000 гр.
6. Салат Ромейн: 2 (4 x 0.5) Вес: 1 гр.
7. Орех пекан очищенный: 1450 (1 x 1450) Вес: 1000 гр.
Payment Amount: 2140.9 RUB
Payment system: (none)

Purchaser information:
name: ООО Тапас
Адрес_доставки_: Пер кооперативный 7
Checkbox: yes
ma_name: ООО «Тапос» - Хуанчо
ma_email: ooo.tapas.tomsk@gmail.com
ma_phone: +79138430719
ma_id: 44278449

Additional information:
Transaction ID: 11091747:7807310862
Block ID: rec820069702
Form Name: Cart
https://bubon-horeca.ru/shop-nuts
-----`

func TestParseMessage(t *testing.T) {
	tests := []struct {
		name    string
		message []byte
		order   *domain.Order
		err     error
	}{
		{
			name:    "Test Parse all",
			message: []byte(messageTest1),
			order: domain.NewOrder("1898307952", []*domain.Product{
				domain.NewProduct("Огурец среднеплодный", domain.NewGramUnit(1000)),
				domain.NewProduct("Томат красный", domain.NewGramUnit(1000)),
				domain.NewProduct("Перец чили красный", domain.NewGramUnit(1)),
				domain.NewProduct("Салат Афицион", domain.NewPieceUnit(1)),
				domain.NewProduct("Авокадо", domain.NewGramUnit(1000)),
			}),
			err: nil,
		},
		{
			name:    "Parse order 1024598132",
			message: []byte(messageTest2),
			order: domain.NewOrder("1024598132", []*domain.Product{
				domain.NewProduct("Баклажан", domain.NewGramUnit(1000)),
				domain.NewProduct("Петрушка", domain.NewGramUnit(1)),
				domain.NewProduct("Картофель желтый", domain.NewGramUnit(1000)),
				domain.NewProduct("Морковь", domain.NewGramUnit(1000)),
				domain.NewProduct("Лук репчатый новый урожай", domain.NewGramUnit(1000)),
				domain.NewProduct("Салат Ромейн", domain.NewGramUnit(1)),
				domain.NewProduct("Орех пекан очищенный", domain.NewGramUnit(1000)),
			}),
			err: nil,
		},
	}

	tests[0].order.Products()[0].Compute(7)
	tests[0].order.Products()[1].Compute(1)
	tests[0].order.Products()[2].Compute(100)
	tests[0].order.Products()[3].Compute(1)
	tests[0].order.Products()[4].Compute(1)

	tests[1].order.Products()[0].Compute(3)
	tests[1].order.Products()[1].Compute(3)
	tests[1].order.Products()[2].Compute(5)
	tests[1].order.Products()[3].Compute(4)
	tests[1].order.Products()[4].Compute(3)
	tests[1].order.Products()[5].Compute(4)
	tests[1].order.Products()[6].Compute(1)

	p := tilda.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := p.ParseOrder(tt.message)
			assert.Equal(t, err, tt.err)
			for _, p := range order.Products() {
				log.Println(order.Id(), p.Name(), p.TotalSize())
			}
			assert.Equal(t, order, tt.order)
		})
	}
}
