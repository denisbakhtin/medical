package models

import ()

//Tag struct contains post tag info
type Tag struct {
	Name string `json:"name"`
	//calculated fields
	PostCount int64 `json:"post_count" db:"post_count"`
}

//Insert stores Tag in db
func (tag *Tag) Insert() error {
	_, err := db.Exec("INSERT INTO tags(name) VALUES($1)", tag.Name)
	return err
}

//Delete removes tag record from db
func (tag *Tag) Delete() error {
	_, err := db.Exec("DELETE FROM tags WHERE name=$1", tag.Name)
	return err
}

//GetPosts returns a slice of associated published posts
func (tag *Tag) GetPosts() ([]Post, error) {
	return GetPostsByTag(tag.Name)
}

//GetTag loads tag record by its name
func GetTag(name interface{}) (*Tag, error) {
	tag := &Tag{}
	err := db.Get(tag, "SELECT * FROM tags WHERE name=$1", name)
	return tag, err
}

//GetTags returns a slice of tags
func GetTags() ([]Tag, error) {
	var list []Tag
	err := db.Select(
		&list,
		`SELECT tags.name, count(poststags.tag_name) as post_count FROM tags 
		LEFT OUTER JOIN poststags ON tags.name=poststags.tag_name 
		GROUP BY tags.name 
		ORDER BY name ASC`,
	)
	return list, err
}

//GetNotEmptyTags returns a slice of tags that have at least one associated blog post
func GetNotEmptyTags() ([]Tag, error) {
	var list []Tag
	err := db.Select(
		&list,
		`SELECT * FROM tags 
		WHERE EXISTS 
		(SELECT null FROM poststags WHERE poststags.tag_name=tags.name) 
		ORDER BY name ASC`,
	)
	return list, err
}
