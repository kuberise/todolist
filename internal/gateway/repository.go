package gateway

import "context"

type Respository interface {
	NewTODO(ctx context.Context, todo string) error
	RemoveTODO(ctx context.Context, todo string) error
	ListTODOS(ctx context.Context) ([]string, error)
}
