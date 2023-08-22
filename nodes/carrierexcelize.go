package nodes

import (
	"errors"
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

type CarrierExcelize struct {
	rows *excelize.Rows
}

func NewCarrierExcelize(source io.Reader) (Carrier, error) {
	f, err := excelize.OpenReader(source)
	fmt.Println(source)
	if err != nil {
		return nil, err
	}
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, nil
	}
	rows, err := f.Rows(sheets[0])
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, errors.New("Empty data carrier")
	}
	row, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	if len(row) < 3 || row[0] != "id" || row[1] != "name" || row[2] != "parent_id" {
		return nil, errors.New("Wrong data format")
	}

	return &CarrierExcelize{rows: rows}, nil
}

func (nc *CarrierExcelize) NextRow() (*Node, error) {
	if !nc.rows.Next() {
		return nil, io.EOF
	}
	n := Node{}
	row, err := nc.rows.Columns()
	if err != nil {
		return nil, err
	}
	n.Id = row[0]
	n.Name = row[1]
	n.ParentId = row[2]
	return &n, nil
}
