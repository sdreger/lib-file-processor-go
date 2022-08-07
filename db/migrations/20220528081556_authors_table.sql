-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE ebook.authors_id_seq AS BIGINT;

CREATE TABLE ebook.authors
(
    id         BIGINT    default nextval('ebook.authors_id_seq'::regclass) NOT NULL,
    name       VARCHAR(255)                                                NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ebook.authors;
DROP SEQUENCE IF EXISTS ebook.authors_id_seq;
-- +goose StatementEnd
