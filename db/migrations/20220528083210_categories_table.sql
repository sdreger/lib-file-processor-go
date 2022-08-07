-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE ebook.categories_id_seq AS BIGINT;

CREATE TABLE ebook.categories
(
    id         BIGINT    default nextval('ebook.categories_id_seq'::regclass) NOT NULL,
    name       VARCHAR(255)                                                   NOT NULL,
    parent_id  BIGINT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ebook.categories;
DROP SEQUENCE IF EXISTS ebook.categories_id_seq;
-- +goose StatementEnd
