package nodes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SaveFromCarrier(r io.Reader) error {
	c, err := NewCarrier(r)
	if err != nil {
		return err
	}
	for {
		row, err := c.NextRow()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		err = row.Validate(s.repo)
		if err != nil {
			return err
		}
		_, err = s.repo.SaveNode(row)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) StartSaving() error {
	err := s.repo.StartTransaction()
	if err != nil {
		return err
	}
	err = s.repo.DeleteAllNodes()
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) FinishSaving() error {
	return s.repo.CommitTransaction()
}

func (s *Service) RollbackSaving() error {
	return s.repo.RollbackTransaction()
}

func (s *Service) WriteJsonNodesTree(w http.ResponseWriter, id int) error {
	nodes, err := s.repo.GetNodeChildren(id)
	if err != nil {
		return err
	}
	fmt.Fprint(w, "[")
	for i, node := range nodes {
		if i > 0 {
			fmt.Fprint(w, ",")
		}
		name, err := json.Marshal(node.Name)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "{\"id\":%d,\"name\":%s,\"children\":", node.Id, name)
		s.WriteJsonNodesTree(w, node.Id)
		fmt.Fprint(w, "}")
	}
	fmt.Fprint(w, "]")
	return nil
}
