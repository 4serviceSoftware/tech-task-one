package carrier

import (
	"io"

	"github.com/4serviceSoftware/tech-task/internal/models"
)

type Carrier interface {
	NextNode() (*models.Node, error)
}

func NewCarrier(source io.Reader) (Carrier, error) {
	return NewCarrierExcelize(source)
}
