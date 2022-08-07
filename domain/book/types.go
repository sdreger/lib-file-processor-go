package book

import "context"

//go:generate mockgen -destination=./store_mock.go -package=book github.com/sdreger/lib-file-processor-go/domain/book Store
type Store interface {
	Find(ctx context.Context, req SearchRequest) (*StoredData, error)
	Add(ctx context.Context, parsedData ParsedData) (int64, error)
	Update(ctx context.Context, existingData *StoredData, parsedData *ParsedData) error
}
