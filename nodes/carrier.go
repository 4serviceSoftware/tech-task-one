package nodes

import "io"

type Carrier interface {
	NextRow() (*Node, error)
}

func NewCarrier(source io.Reader) (Carrier, error) {
	return NewCarrierExcelize(source)
}
