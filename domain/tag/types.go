package tag

import "context"

//go:generate mockgen -destination=./store_mock.go -package=tag github.com/sdreger/lib-file-processor-go/domain/tag Store
type Store interface {
	UpsertAll(ctx context.Context, tags []string) ([]int64, error)
	ReplaceBookTags(ctx context.Context, bookID int64, tagIDs []int64) error
}
