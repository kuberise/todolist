package gateway

import "context"

type Respository interface {
	Index(ctx context.Context) ([]string, error)
}
