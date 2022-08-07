-- +goose Up
-- +goose StatementBegin
CREATE TABLE ebook.book_author
(
    book_id   BIGINT NOT NULL,
    author_id BIGINT NOT NULL,
    PRIMARY KEY (book_id, author_id)
);

ALTER TABLE ebook.book_author
    ADD CONSTRAINT fk_book
        FOREIGN KEY (book_id)
            REFERENCES ebook.books (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;

ALTER TABLE ebook.book_author
    ADD CONSTRAINT fk_author
        FOREIGN KEY (author_id)
            REFERENCES ebook.authors (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ebook.book_author
    DROP CONSTRAINT fk_author;
ALTER TABLE ebook.book_author
    DROP CONSTRAINT fk_book;
DROP TABLE ebook.book_author;
-- +goose StatementEnd
