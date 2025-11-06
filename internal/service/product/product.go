package product

import (
	"productsParser/internal/domain"
	"sort"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

type Products struct {
	p []*domain.Product
}

func NewService(products []*domain.Product) *Products {
	return &Products{
		p: products,
	}
}

func (ps *Products) ProductsSortByName() []*domain.Product {
	collator := collate.New(language.Russian, collate.IgnoreCase)
	sort.Slice(ps.p, func(i, j int) bool {
		return collator.CompareString(ps.p[i].Name(), ps.p[j].Name()) < 0
	})
	return ps.p
}

func (ps *Products) Products() []*domain.Product {
	return ps.p
}
