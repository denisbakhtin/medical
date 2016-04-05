-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE comments
  ADD COLUMN author_city TEXT NOT NULL DEFAULT '';


-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back
ALTER TABLE comments DROP COLUMN author_city;
