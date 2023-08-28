package nodes

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/4serviceSoftware/tech-task/internal/models"
)

type Service struct {
	repo      Repository
	cacheFile *CacheFile
}

func NewService(repo Repository, cacheFile *CacheFile) *Service {
	return &Service{repo: repo, cacheFile: cacheFile}
}

func (s *Service) SaveFromMultipartReader(multipartReader *multipart.Reader) error {
	// initializing new saving session
	err := s.startSaving()
	defer s.rollbackSaving()
	if err != nil {
		return err
	}

	// Loop through each part of the request body
	for {
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// if part is a file than handle this part
		if len(part.FileName()) > 0 {
			err = s.saveFromCarrier(part)
			if err != nil {
				return err
			}
		}
	}

	err = s.finishSaving()
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) WriteCachedJsonNodesTree(w io.Writer) error {
	cacheFile, err := s.cacheFile.GetFileReader()
	if err != nil {
		return s.writeJsonNodesTree(w, 0)
	}
	_, err = io.Copy(w, cacheFile)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) saveFromCarrier(r io.Reader) error {
	c, err := NewCarrier(r)
	if err != nil {
		return err
	}
	for {
		node, err := c.NextNode()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		err = s.checkNodeCircularRef(node)
		if err != nil {
			return err
		}
		_, err = s.repo.SaveNode(node)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) writeJsonNodesTree(w io.Writer, id int) error {
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
		s.writeJsonNodesTree(w, node.Id)
		fmt.Fprint(w, "}")
	}
	fmt.Fprint(w, "]")
	return nil
}

func (s *Service) checkNodeCircularRef(node *models.Node) error {
	// checking this node for circular dependency in parents
	parents, err := s.repo.GetNodeParents(node.ParentId)
	if err != nil {
		return err
	}
	for _, p := range parents {
		if p.ParentId == node.Id {
			return fmt.Errorf("Node validation error. Circular reference. Node {id: %d, name: %s, parent_id: %d}",
				node.Id, node.Name, node.ParentId)
		}
	}
	return nil
}

func (s *Service) startSaving() error {
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

func (s *Service) finishSaving() error {
	err := s.repo.CommitTransaction()
	if err != nil {
		return err
	}
	go func() {
		cacheFileWriter, closeFileFunc, err := s.cacheFile.GetNewFileWriter()
		if err != nil {
			return
		}
		s.writeJsonNodesTree(cacheFileWriter, 0)
		closeFileFunc()
	}()
	return nil
}

func (s *Service) rollbackSaving() error {
	return s.repo.RollbackTransaction()
}
