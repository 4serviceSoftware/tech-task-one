package nodes

import "fmt"

type Node struct {
	Id       int
	Name     string
	ParentId int
}

func (n *Node) Validate(repo Repository) error {
	// checking this node for circular dependency in parents
	parents, err := repo.GetNodeParents(n.ParentId)
	if err != nil {
		return err
	}
	for _, p := range parents {
		if p.ParentId == n.Id {
			return fmt.Errorf("Node validation error. Circular reference. Node {id: %d, name: %s, parent_id: %d",
				n.Id, n.Name, n.ParentId)
		}
	}
	return nil
}
