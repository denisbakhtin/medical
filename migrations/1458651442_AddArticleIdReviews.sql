-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE reviews
  ADD COLUMN article_id INTEGER REFERENCES articles(id) ON DELETE SET NULL;


-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back
ALTER TABLE reviews DROP COLUMN article_id;
