-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE ebook.tags_id_seq AS BIGINT;

CREATE TABLE ebook.tags
(
    id         BIGINT    default nextval('ebook.tags_id_seq'::regclass) NOT NULL,
    name       VARCHAR(255)                                             NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ebook.tags;
DROP SEQUENCE IF EXISTS ebook.tags_id_seq;
-- +goose StatementEnd
