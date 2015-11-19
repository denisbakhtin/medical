-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE articles(
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  content TEXT NOT NULL,
  published BOOLEAN NOT NULL DEFAULT true,
  user_id INTEGER REFERENCES users (id) ON DELETE SET NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);


-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back
DROP TABLE articles;

