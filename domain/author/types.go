package author

import "context"

//go:generate mockgen -destination=./store_mock.go -package=author github.com/sdreger/lib-file-processor-go/domain/author Store
type Store interface {
	UpsertAll(ctx context.Context, authors []string) ([]int64, error)
	ReplaceBookAuthors(ctx context.Context, bookID int64, authorIDs []int64) error
}
