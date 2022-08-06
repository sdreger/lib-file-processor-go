package filetype

import "context"

//go:generate mockgen -destination=./store_mock.go -package=filetype github.com/sdreger/lib-file-processor-go/domain/filetype Store
type Store interface {
	UpsertAll(ctx context.Context, fileTypes []string) ([]int64, error)
	ReplaceBookFileTypes(ctx context.Context, bookID int64, fileTypeIDs []int64) error
}
