package domain

import (
	"fmt"
	"math/big"
)

const (
	gramCode          = "gram"
	toStringGramUnit  = "кг"
	pieceCode         = "piece"
	toStringPieceUnit = "шт"
)

type CommonUnit struct {
	totalSize int64
	value     int
}

func NewCommonUnit(value int) *CommonUnit {
	return &CommonUnit{
		totalSize: 0,
		value:     value,
	}
}

func (c *CommonUnit) TotalSize() int64 {
	return c.totalSize
}

func (c *CommonUnit) Compute(quantity int) int64 {
	c.totalSize += int64(c.value * quantity)
	return c.totalSize
}

type GramUnit struct {
	*CommonUnit
}

func NewGramUnit(value int) *GramUnit {
	return &GramUnit{CommonUnit: NewCommonUnit(value)}
}

func (g *GramUnit) ToString() string {
	grams := big.NewRat(g.TotalSize(), 1)              // g.size — например, 1234
	kg := new(big.Rat).Quo(grams, big.NewRat(1000, 1)) // делим точно
	return fmt.Sprintf("%s кг.", kg.FloatString(3))    // 3 знака после запятой
}

func (g *GramUnit) Code() string {
	return gramCode
}

type PieceUnit struct {
	*CommonUnit
}

func NewPieceUnit(value int) *PieceUnit {
	return &PieceUnit{CommonUnit: NewCommonUnit(value)}
}

func (p *PieceUnit) ToString() string {
	return fmt.Sprintf("%d шт.", p.TotalSize())
}

func (p *PieceUnit) Code() string {
	return pieceCode
}
