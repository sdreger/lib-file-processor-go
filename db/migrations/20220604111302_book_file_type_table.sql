-- +goose Up
-- +goose StatementBegin
CREATE TABLE ebook.book_file_type
(
    book_id   BIGINT NOT NULL,
    file_type_id BIGINT NOT NULL,
    PRIMARY KEY (book_id, file_type_id)
);

ALTER TABLE ebook.book_file_type
    ADD CONSTRAINT fk_book
        FOREIGN KEY (book_id)
            REFERENCES ebook.books (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;

ALTER TABLE ebook.book_file_type
    ADD CONSTRAINT fk_file_type
        FOREIGN KEY (file_type_id)
            REFERENCES ebook.file_types (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ebook.book_file_type
    DROP CONSTRAINT fk_file_type;
ALTER TABLE ebook.book_file_type
    DROP CONSTRAINT fk_book;
DROP TABLE ebook.book_file_type;
-- +goose StatementEnd
