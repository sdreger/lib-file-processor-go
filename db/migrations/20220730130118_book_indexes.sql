-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS isbn10_unique ON ebook.books (isbn10);
CREATE UNIQUE INDEX IF NOT EXISTS isbn13_unique ON ebook.books (isbn13);
CREATE UNIQUE INDEX IF NOT EXISTS asin_unique ON ebook.books (asin);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS ebook.isbn10_unique;
DROP INDEX IF EXISTS ebook.isbn13_unique;
DROP INDEX IF EXISTS ebook.asin_unique;
-- +goose StatementEnd
