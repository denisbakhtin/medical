-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE reviews
  ADD COLUMN video TEXT NOT NULL DEFAULT '';
UPDATE reviews SET video='';

-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back
ALTER TABLE reviews DROP COLUMN video;
