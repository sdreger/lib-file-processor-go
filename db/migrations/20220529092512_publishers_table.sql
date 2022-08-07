-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE ebook.publishers_id_seq AS BIGINT;

CREATE TABLE ebook.publishers
(
    id         BIGINT    default nextval('ebook.publishers_id_seq'::regclass) NOT NULL,
    name       VARCHAR(255)                                                   NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ebook.publishers;
DROP SEQUENCE IF EXISTS ebook.publishers_id_seq;
-- +goose StatementEnd
