package nodes

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/4serviceSoftware/tech-task/internal/models"
	"github.com/4serviceSoftware/tech-task/internal/nodes/cachefile"
	"github.com/4serviceSoftware/tech-task/internal/nodes/carrier"
	"github.com/4serviceSoftware/tech-task/internal/repos"
)

type Service struct {
	repo      repos.NodesRepository
	cacheFile *cachefile.CacheFile
}

func NewService(repo repos.NodesRepository, cacheFile *cachefile.CacheFile) *Service {
	return &Service{repo: repo, cacheFile: cacheFile}
}

// SaveFromMultipartReader takes multipart.Reader, reads parts (files) from it
// and calls saveFromCarrier() method to process this parts.
// Returns an error or nil.
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

// WriteCachedJsonNodesTree checks if response is stored in cache file or
// it needs to be generated. In any case responce is writing to a givan io.Writer
// directly without storing responce in a memory.
// Response is a json formatted nodes tree.
// Returns error or nil
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

// saveFromCarrier takes io.Reader as a handler of file with stored nodes (carrier),
// reads nodes from it, checks for errors and saves to repository
// Returns error or nil
func (s *Service) saveFromCarrier(r io.Reader) error {
	carrier, err := carrier.NewCarrier(r)
	if err != nil {
		return err
	}
	for {
		node, err := carrier.NextNode()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// check nodes parents chain for circular references
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

// writeJsonNodesTree builds a json formatted nodes tree and writes it
// to a given io.Writer directly
// Returns error or nil
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

// checkNodeCircularRef checks node parents chain for circular references
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

// startSaving does a preparations before saving all nodes tree to repository.
// It starts a transaction and clears all previous nodes
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

// finishSaving does all needed work after saving all nodes tree to repository.
// It commits a transaction and refreshes cache file
func (s *Service) finishSaving() error {
	err := s.repo.CommitTransaction()
	if err != nil {
		return err
	}
	// refresh cache file in parallel
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

// finishSaving does rollback for saving nodes tree session
// in case of an error in the middle of saving session.
func (s *Service) rollbackSaving() error {
	return s.repo.RollbackTransaction()
}
