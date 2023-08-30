package repos

import (
	"context"

	"github.com/4serviceSoftware/tech-task/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type NodesRepositoryPostgres struct {
	ctx  context.Context
	conn *pgxpool.Pool
	tx   pgx.Tx
}

func NewNodesRepositoryPostgres(ctx context.Context, conn *pgxpool.Pool) *NodesRepositoryPostgres {
	return &NodesRepositoryPostgres{ctx: ctx, conn: conn, tx: nil}
}

func (r *NodesRepositoryPostgres) DeleteAllNodes() error {
	if r.tx != nil {

		_, err := r.tx.Exec(r.ctx, "DELETE FROM nodes")
		return err
	} else {
		_, err := r.conn.Exec(r.ctx, "DELETE FROM nodes")
		return err
	}
}

func (r *NodesRepositoryPostgres) SaveNode(n *models.Node) (int, error) {
	var row pgx.Row
	query := "INSERT INTO nodes (id, name, parent_id) VALUES ($1, $2, $3) RETURNING id"
	if r.tx != nil {
		row = r.tx.QueryRow(r.ctx, query, n.Id, n.Name, n.ParentId)
	} else {
		row = r.conn.QueryRow(r.ctx, query, n.Id, n.Name, n.ParentId)
	}
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *NodesRepositoryPostgres) GetNodeParents(id int) ([]*models.Node, error) {
	query := `WITH RECURSIVE  Parents
				AS
				(
						-- anchor
					SELECT  id, name, parent_id, 1 as depth
					FROM    nodes t1
					WHERE   id=$1

					UNION ALL
						--recursive member
					SELECT  t2.id, t2.name, t2.parent_id, depth+1
					FROM    nodes AS t2
							JOIN Parents AS M ON t2.id = M.parent_id
					where depth<1000
				)

				SELECT id, name, parent_id FROM Parents
				`
	var rows pgx.Rows
	var err error
	if r.tx != nil {
		rows, err = r.tx.Query(r.ctx, query, id)
	} else {
		rows, err = r.conn.Query(r.ctx, query, id)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var nodes []*models.Node
	for rows.Next() {
		n := models.Node{}
		err = rows.Scan(&n.Id, &n.Name, &n.ParentId)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, &n)
	}
	return nodes, nil
}

func (r *NodesRepositoryPostgres) GetNodeChildren(id int) ([]*models.Node, error) {
	query := `SELECT id,name,parent_id 
				FROM nodes
				WHERE parent_id=$1
				ORDER BY id
				`
	var rows pgx.Rows
	var err error
	if r.tx != nil {
		rows, err = r.tx.Query(r.ctx, query, id)
	} else {
		rows, err = r.conn.Query(r.ctx, query, id)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var nodes []*models.Node
	for rows.Next() {
		n := models.Node{}
		err = rows.Scan(&n.Id, &n.Name, &n.ParentId)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, &n)
	}
	return nodes, nil
}

func (r *NodesRepositoryPostgres) StartTransaction() error {
	var err error
	r.tx, err = r.conn.Begin(r.ctx)
	if err != nil {
		r.tx = nil
		return err
	}
	return nil
}

func (r *NodesRepositoryPostgres) CommitTransaction() error {
	if r.tx != nil {
		err := r.tx.Commit(r.ctx)
		r.tx = nil
		return err
	}
	return nil
}

func (r *NodesRepositoryPostgres) RollbackTransaction() error {
	if r.tx != nil {
		err := r.tx.Rollback(r.ctx)
		r.tx = nil
		return err
	}
	return nil
}
