package tilda

import (
	"errors"
	"fmt"
	"io"
	"productsParser/internal/domain"
	"regexp"
	"strconv"
	"strings"
)

const (
	pcStr       = "pc"
	pcRuStr     = "шт."
	gramStr     = "гр."
	kgStrRu     = "кг."
	kgStr       = "kg"
	rePcStr     = `^\d+\.\s*(.+?),.*\(\s*(\d+)\s*([a-zA-Z]+)`
	reWeightstr = `^\d+\.\s*(.+?):\s*[\d\.,]+\s*\(([\d\s\w]+?)x[\d\s\.,]+?\)\s*Вес:\s*(\d+)\s*(\S+)`
	reOrderStr  = `Order\s+#(\d+)`
)

var (
	reWeight = regexp.MustCompile(reWeightstr)
	rePc     = regexp.MustCompile(rePcStr)
	reOrder  = regexp.MustCompile(reOrderStr)
)

var unitInit = map[string]func(value int) domain.Uniter{
	pcStr:   func(value int) domain.Uniter { return domain.NewPieceUnit(value) },
	pcRuStr: func(value int) domain.Uniter { return domain.NewPieceUnit(value) },
	kgStr:   func(value int) domain.Uniter { return domain.NewGramUnit(1000) },
	gramStr: func(value int) domain.Uniter { return domain.NewGramUnit(value) },
	kgStrRu: func(value int) domain.Uniter { return domain.NewGramUnit(value) },
}

type Parser struct {
	reWeight *regexp.Regexp
	rePc     *regexp.Regexp
}

func New() *Parser {
	return &Parser{
		reWeight: reWeight,
		rePc:     rePc,
	}
}

func parseInt(s string) int {
	s = strings.TrimSpace(strings.Split(s, " ")[0])
	s = strings.ReplaceAll(s, ",", ".")
	f, _ := strconv.ParseFloat(s, 64)
	return int(f)
}

func (p *Parser) ParseOrder(b []byte) (*domain.Order, error) {
	strB := string(b)
	lines := strings.Split(strB, "\n")
	var id string

	m := reOrder.FindStringSubmatch(lines[0])
	if m == nil {
		return nil, errors.New("invalid message - orderId is empty")
	}
	id = m[1]

	products := make([]*domain.Product, 0)

	for _, line := range lines[1:] {
		line := strings.TrimSpace(line)
		if len(line) != 0 {
			p, err := p.pareseProduct(line)
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}

			products = append(products, p)
		}
	}

	if len(products) == 0 {
		return nil, errors.New("parse error: empty products find")
	}
	return domain.NewOrder(id, products), nil
}

func (p *Parser) pareseProduct(str string) (*domain.Product, error) {
	// Попробуем разные шаблоны
	if m := reWeight.FindStringSubmatch(str); m != nil {
		name := strings.TrimSpace(m[1])
		qty := parseInt(m[2])
		weight := parseInt(m[3])
		unit := strings.TrimSpace(m[4])

		return p.validProduct(name, unit, qty, weight)
	}

	if m := rePc.FindStringSubmatch(str); m != nil {
		name := strings.TrimSpace(m[1])
		qty := parseInt(m[2])
		pc := strings.TrimSpace(m[3])
		return p.validProduct(name, pc, qty, 1)
	}

	return nil, io.EOF
}

func (p *Parser) validProduct(name, unit string, qty, weight int) (*domain.Product, error) {
	if len(name) == 0 {
		return nil, errors.New("parse error: name product is empty")
	}

	if qty == 0 {
		return nil, errors.New("parse error: qty is zero")
	}

	if weight == 0 {
		return nil, errors.New("parse error: weight is empty")
	}

	unitInitFunc, ok := unitInit[unit]
	if !ok {
		return nil, fmt.Errorf("parse error: unknown unit %s", unit)
	}

	product := domain.NewProduct(name)
	product.Compute(unitInitFunc(weight), qty)
	return product, nil
}
