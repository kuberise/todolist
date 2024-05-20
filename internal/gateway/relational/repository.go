package relational

import (
	"context"
	"database/sql"

	"github.com/kuberise/todolist/internal/gateway"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) gateway.Respository {

	return &repository{db: db}
}

func (r *repository) Index(ctx context.Context) ([]string, error) {

	var res string
	var todos []string

	rows, err := r.db.QueryContext(ctx, "SELECT * FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&res)
		todos = append(todos, res)
	}

	return todos, nil
}
