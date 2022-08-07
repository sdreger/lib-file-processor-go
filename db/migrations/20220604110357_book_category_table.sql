-- +goose Up
-- +goose StatementBegin
CREATE TABLE ebook.book_category
(
    book_id   BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    PRIMARY KEY (book_id, category_id)
);

ALTER TABLE ebook.book_category
    ADD CONSTRAINT fk_book
        FOREIGN KEY (book_id)
            REFERENCES ebook.books (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;

ALTER TABLE ebook.book_category
    ADD CONSTRAINT fk_category
        FOREIGN KEY (category_id)
            REFERENCES ebook.categories (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ebook.book_category
    DROP CONSTRAINT fk_category;
ALTER TABLE ebook.book_category
    DROP CONSTRAINT fk_book;
DROP TABLE ebook.book_category;
-- +goose StatementEnd
