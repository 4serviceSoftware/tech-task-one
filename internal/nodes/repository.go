package nodes

import "github.com/4serviceSoftware/tech-task/internal/models"

type Repository interface {
	DeleteAllNodes() error
	SaveNode(n *models.Node) (int, error)
	GetNodeParents(id int) ([]*models.Node, error)
	GetNodeChildren(id int) ([]*models.Node, error)
	StartTransaction() error
	CommitTransaction() error
	RollbackTransaction() error
}
