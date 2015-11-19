-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE comments(
  id SERIAL PRIMARY KEY,
  article_id INTEGER REFERENCES articles(id) ON DELETE CASCADE,
  author_name TEXT NOT NULL,
  author_email TEXT NOT NULL,
  content TEXT NOT NULL,
  answer TEXT,
  published BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);


-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back
DROP TABLE comments;

