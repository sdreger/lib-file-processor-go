package category

import "context"

//go:generate mockgen -destination=./store_mock.go -package=category github.com/sdreger/lib-file-processor-go/domain/category Store
type Store interface {
	UpsertAll(ctx context.Context, categories []string) ([]int64, error)
	ReplaceBookCategories(ctx context.Context, bookID int64, categoryIDs []int64) error
}
