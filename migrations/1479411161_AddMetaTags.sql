-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE pages ADD COLUMN meta_keywords TEXT NOT NULL DEFAULT '';
ALTER TABLE pages ADD COLUMN meta_description TEXT NOT NULL DEFAULT '';
UPDATE pages SET meta_keywords='';
UPDATE pages SET meta_description='';
ALTER TABLE reviews ADD COLUMN meta_keywords TEXT NOT NULL DEFAULT '';
ALTER TABLE reviews ADD COLUMN meta_description TEXT NOT NULL DEFAULT '';
UPDATE reviews SET meta_keywords='';
UPDATE reviews SET meta_description='';
ALTER TABLE articles ADD COLUMN meta_keywords TEXT NOT NULL DEFAULT '';
ALTER TABLE articles ADD COLUMN meta_description TEXT NOT NULL DEFAULT '';
UPDATE articles SET meta_keywords='';
UPDATE articles SET meta_description='';

-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back
ALTER TABLE pages DROP COLUMN meta_keywords;
ALTER TABLE pages DROP COLUMN meta_description;
ALTER TABLE reviews DROP COLUMN meta_keywords;
ALTER TABLE reviews DROP COLUMN meta_description;
ALTER TABLE articles DROP COLUMN meta_keywords;
ALTER TABLE articles DROP COLUMN meta_description;
