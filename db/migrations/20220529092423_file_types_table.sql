-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE ebook.file_types_id_seq AS BIGINT;

CREATE TABLE ebook.file_types
(
    id         BIGINT    default nextval('ebook.file_types_id_seq'::regclass) NOT NULL,
    name       VARCHAR(255)                                                     NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ebook.file_types;
DROP SEQUENCE IF EXISTS ebook.file_types_id_seq;
-- +goose StatementEnd
