package domain

import "errors"

type Uniter interface {
	Compute(int) int64
	Code() string
	ToString(totalSize int64) string
}

type Product struct {
	totalSize int64
	unit      Uniter
	name      string
}

func NewProduct(name string, unit Uniter) *Product {
	return &Product{
		name:      name,
		unit:      unit,
		totalSize: int64(0),
	}
}

func (p *Product) Name() string {
	return p.name
}

func (p *Product) ToString() string {
	return p.unit.ToString(p.totalSize)
}

func (p *Product) Uniter() Uniter {
	return p.unit
}

func (p *Product) Compute(quantity int) int64 {
	p.totalSize = p.unit.Compute(quantity)
	return p.totalSize
}

func (p *Product) TotalSize() int64 {
	return p.totalSize
}

func (p *Product) AddProduct(someProduct *Product) error {
	if p.name != someProduct.name {
		return errors.New("name is not equal")
	}
	if p.unit.Code() != someProduct.unit.Code() {
		return errors.New("unit code is not equal")
	}

	p.totalSize += someProduct.totalSize
	return nil
}
