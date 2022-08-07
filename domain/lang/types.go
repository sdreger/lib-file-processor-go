package lang

import "context"

//go:generate mockgen -destination=./store_mock.go -package=lang github.com/sdreger/lib-file-processor-go/domain/lang Store
type Store interface {
	Upsert(ctx context.Context, language string) (int64, error)
}
