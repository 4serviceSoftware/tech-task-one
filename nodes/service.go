package nodes

import (
	"fmt"
	"io"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (ns *Service) SaveFromCarrier(r io.Reader) error {
	nc, err := NewCarrier(r)
	if err != nil {
		return err
	}
	for row, err := nc.NextRow(); err != io.EOF; {
		fmt.Println(row)
		row, err = nc.NextRow()
	}
	return nil
}
