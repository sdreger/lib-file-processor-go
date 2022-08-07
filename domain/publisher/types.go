package publisher

import "context"

//go:generate mockgen -destination=./store_mock.go -package=publisher github.com/sdreger/lib-file-processor-go/domain/publisher Store
type Store interface {
	Upsert(ctx context.Context, publisher string) (int64, error)
}
