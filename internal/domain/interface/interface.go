package i

import "productsParser/internal/domain"

type Parser interface {
	ParseOrder([]byte) (*domain.Order, error)
}
