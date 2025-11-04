package domain

import (
	"errors"
	"strings"
)

type Uniter interface {
	Compute(int) int64
	Add(int64) int64
	TotalSize() int64
	Code() string
	ToString() string
}

type Product struct {
	name    string
	uniters map[string]Uniter
}

func NewProduct(name string) *Product {
	return &Product{
		name:    name,
		uniters: make(map[string]Uniter),
	}
}

func (p *Product) Name() string {
	return p.name
}

func (p *Product) ToString() string {
	units := make([]string, 0, len(p.uniters))
	for _, unit := range p.uniters {
		units = append(units, unit.ToString())
	}
	return strings.Join(units, ", ")
}

func (p *Product) Uniters() map[string]Uniter {
	return p.uniters
}

func (p *Product) Compute(uniter Uniter, quantity int) int64 {
	u, ok := p.uniters[uniter.Code()]
	if !ok {
		u = uniter
		p.uniters[uniter.Code()] = uniter
	}

	return u.Compute(quantity)
}

func (p *Product) AddUniters(uniters map[string]Uniter) {
	for _, uniter := range uniters {
		_, ok := p.uniters[uniter.Code()]
		if !ok {
			p.uniters[uniter.Code()] = uniter
		}

		if ok {
			p.uniters[uniter.Code()].Add(uniter.TotalSize())
		}
	}
}

func (p *Product) AddProduct(someProduct *Product) error {
	if p.name != someProduct.name {
		return errors.New("name is not equal")
	}
	p.AddUniters(someProduct.uniters)
	return nil
}
