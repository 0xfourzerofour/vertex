package schemas

import (
	"context"
)

type SchemaRepository interface {
	GetMaster(ctx context.Context) error
	ListSubSchemas(ctx context.Context) ([][]byte, error)
	Merge(ctx context.Context, schemaArr [][]byte) error
}
