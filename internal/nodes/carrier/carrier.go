// This is interface of carrier.
// Carrier is a concrete form of incoming data. For this project carrier is a xlsx file.
// This interface contains only one method - NextNode()
package carrier

import (
	"io"

	"github.com/4serviceSoftware/tech-task/internal/models"
)

type Carrier interface {
	// NextNode reads next node from carrier
	NextNode() (*models.Node, error)
}

func NewCarrier(source io.Reader) (Carrier, error) {
	return NewCarrierExcelize(source)
}
