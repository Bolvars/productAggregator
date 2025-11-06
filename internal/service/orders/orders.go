package service

import (
	"errors"
	"productsParser/internal/domain"
	service "productsParser/internal/service/product"
	"sync"
)

type Parser interface {
	ParseOrder([]byte) (*domain.Order, error)
}

type Orders struct {
	rw     *sync.RWMutex
	user   *domain.User
	aggr   map[string]*domain.Product
	result []*domain.Product
	orders map[string]*domain.Order
	p      Parser
}

func NewOrderService(user *domain.User, p Parser) *Orders {
	return &Orders{
		user:   user,
		aggr:   make(map[string]*domain.Product),
		orders: make(map[string]*domain.Order),
		result: make([]*domain.Product, 0),
		rw:     &sync.RWMutex{},
		p:      p,
	}
}

func (os *Orders) AddOrder(b []byte) (*domain.Order, error) {
	order, err := os.p.ParseOrder(b)
	if err != nil {
		return nil, err
	}
	os.rw.Lock()
	defer os.rw.Unlock()
	os.orders[order.Id()] = order
	return order, nil
}

func (os *Orders) Compute() (*service.Products, error) {
	os.rw.Lock()
	defer os.rw.Unlock()

	if len(os.orders) == 0 {
		return nil, errors.New("no orders to compute")
	}

	for _, order := range os.orders {
		for _, product := range order.Products() {
			if err := os.aggrProduct(product); err != nil {
				return nil, err
			}
		}
	}
	return service.NewService(os.result), nil
}

func (os *Orders) aggrProduct(product *domain.Product) error {
	aggrProdutct, ok := os.aggr[product.Name()]
	if !ok {
		os.aggr[product.Name()] = product
		os.result = append(os.result, product)
		return nil
	}
	if err := aggrProdutct.AddProduct(product); err != nil {
		return err
	}

	return nil
}
