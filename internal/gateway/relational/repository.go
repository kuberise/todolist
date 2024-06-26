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

func (r *repository) RemoveTODO(ctx context.Context, todo string) error {
	_, err := r.db.ExecContext(ctx, "DELETE from todos WHERE title=$1", todo)
	return err
}

func (r *repository) NewTODO(ctx context.Context, todo string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO todos VALUES ($1)", todo)
	return err
}

func (r *repository) ListTODOS(ctx context.Context) ([]string, error) {

	var res string
	todos := []string{}

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
