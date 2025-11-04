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
	value int
}

func NewCommonUnit(value int) *CommonUnit {
	return &CommonUnit{
		value: value,
	}
}

func (c *CommonUnit) Compute(quantity int) int64 {
	return int64(c.value * quantity)
}

type GramUnit struct {
	*CommonUnit
}

func NewGramUnit(value int) *GramUnit {
	return &GramUnit{CommonUnit: NewCommonUnit(value)}
}

func (g *GramUnit) ToString(totalSize int64) string {
	grams := big.NewRat(totalSize, 1)                  // g.size — например, 1234
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

func (p *PieceUnit) ToString(totalSize int64) string {
	return fmt.Sprintf("%d шт.", totalSize)
}

func (p *PieceUnit) Code() string {
	return pieceCode
}
