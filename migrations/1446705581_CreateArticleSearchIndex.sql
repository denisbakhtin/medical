-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE INDEX articles_search_idx ON articles USING gin(to_tsvector('russian', name || ' ' || content));

-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back
DROP INDEX articles_search_idx;

