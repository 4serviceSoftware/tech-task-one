package nodes

type Repository interface {
	DeleteAllNodes() error
	SaveNode(n *Node) (int, error)
	GetNodeParents(id int) ([]*Node, error)
	StartTransaction() error
	CommitTransaction() error
	RollbackTransaction() error
}
