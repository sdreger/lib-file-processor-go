-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE ebook.books_id_seq AS BIGINT;

CREATE TABLE ebook.books
(
    id              BIGINT        default nextval('ebook.books_id_seq'::regclass) NOT NULL,
    title           VARCHAR(1024)                                                 NOT NULL,
    subtitle        VARCHAR(1024) DEFAULT NULL,
    description     TEXT                                                          NOT NULL,
    isbn10          VARCHAR(10)   DEFAULT NULL,
    isbn13          BIGINT        DEFAULT NULL,
    asin            VARCHAR(10)   DEFAULT NULL,
    pages           SMALLINT                                                      NOT NULL,
    language_id     BIGINT                                                        NOT NULL,
    publisher_id    BIGINT                                                        NOT NULL,
    publisher_url   VARCHAR(255)                                                  NOT NULL,
    edition         SMALLINT                                                      NOT NULL,
    pub_date        DATE                                                          NOT NULL,
    book_file_name  VARCHAR(255)                                                  NOT NULL,
    book_file_size  BIGINT                                                        NOT NULL,
    cover_file_name VARCHAR(255)                                                  NOT NULL,
    created_at      TIMESTAMP     DEFAULT now(),
    updated_at      TIMESTAMP     DEFAULT now(),
    PRIMARY KEY (id)
);

ALTER TABLE ebook.books
    ADD CONSTRAINT fk_language
        FOREIGN KEY (language_id)
            REFERENCES ebook.languages (id)
            ON DELETE SET NULL
            ON UPDATE CASCADE;

ALTER TABLE ebook.books
    ADD CONSTRAINT fk_publisher
        FOREIGN KEY (publisher_id)
            REFERENCES ebook.publishers (id)
            ON DELETE SET NULL
            ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE ebook.books
    DROP CONSTRAINT fk_publisher;
ALTER TABLE ebook.books
    DROP CONSTRAINT fk_language;
DROP TABLE IF EXISTS ebook.books;
DROP SEQUENCE IF EXISTS ebook.books_id_seq;
-- +goose StatementEnd
