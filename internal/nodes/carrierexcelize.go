package nodes

import (
	"errors"
	"io"
	"strconv"

	"github.com/4serviceSoftware/tech-task/internal/models"
	"github.com/xuri/excelize/v2"
)

type CarrierExcelize struct {
	rows *excelize.Rows
}

func NewCarrierExcelize(source io.Reader) (Carrier, error) {
	f, err := excelize.OpenReader(source)
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

func (nc *CarrierExcelize) NextNode() (*models.Node, error) {
	if !nc.rows.Next() {
		return nil, io.EOF
	}
	n := models.Node{}
	row, err := nc.rows.Columns()
	if err != nil {
		return nil, err
	}
	n.Id, err = strconv.Atoi(row[0])
	if err != nil {
		return nil, err
	}
	n.Name = row[1]
	n.ParentId, err = strconv.Atoi(row[2])
	if err != nil {
		return nil, err
	}
	return &n, nil
}
