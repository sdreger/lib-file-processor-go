-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE ebook.languages_id_seq AS BIGINT;

CREATE TABLE ebook.languages
(
    id         BIGINT    default nextval('ebook.languages_id_seq'::regclass) NOT NULL,
    name       VARCHAR(255)                                             NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ebook.languages;
DROP SEQUENCE IF EXISTS ebook.languages_id_seq;
-- +goose StatementEnd
