-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE articles
  ADD COLUMN selling_preface TEXT NOT NULL DEFAULT '';


-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back
ALTER TABLE articles DROP COLUMN selling_preface;
