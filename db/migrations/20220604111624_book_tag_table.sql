-- +goose Up
-- +goose StatementBegin
CREATE TABLE ebook.book_tag
(
    book_id   BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    PRIMARY KEY (book_id, tag_id)
);

ALTER TABLE ebook.book_tag
    ADD CONSTRAINT fk_book
        FOREIGN KEY (book_id)
            REFERENCES ebook.books (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;

ALTER TABLE ebook.book_tag
    ADD CONSTRAINT fk_tag
        FOREIGN KEY (tag_id)
            REFERENCES ebook.tags (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ebook.book_tag
    DROP CONSTRAINT fk_tag;
ALTER TABLE ebook.book_tag
    DROP CONSTRAINT fk_book;
DROP TABLE ebook.book_tag;
-- +goose StatementEnd
