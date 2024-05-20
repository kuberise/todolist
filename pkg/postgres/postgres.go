package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func New(c *Config) (*sql.DB, error) {

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)

	return sql.Open("postgres", connStr)

}
