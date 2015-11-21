-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE reviews(
  id SERIAL PRIMARY KEY,
  content TEXT NOT NULL,
  author_name TEXT NOT NULL,
  author_email TEXT NOT NULL,
  image TEXT NOT NULL,
  published BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);


-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back
DROP TABLE reviews;

