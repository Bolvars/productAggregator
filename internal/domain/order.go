package domain

type Order struct {
	id       string
	products []*Product
}

func NewOrder(id string, products []*Product) *Order {
	return &Order{
		id:       id,
		products: products,
	}
}

func (o *Order) Id() string {
	return o.id
}

func (o *Order) Products() []*Product {
	return o.products
}
