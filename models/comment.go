package models

import "time"

//Comment type contains article comments
type Comment struct {
	ID          uint      `json:"id" form:"id"`
	ArticleID   uint      `json:"article_id" form:"article_id"`
	AuthorCity  string    `json:"author_city" form:"author_city"`
	AuthorName  string    `json:"author_name" form:"author_name"`
	AuthorEmail string    `json:"author_email" form:"author_email"`
	Content     string    `json:"content" form:"content"`
	Answer      string    `json:"answer" form:"answer"`
	Published   bool      `json:"published" form:"published"`
	Captcha     string    `json:"captcha" form:"captcha" gorm:"-" db:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

//Title returns comment excerpt, 100 char long
func (comment *Comment) Title() string {
	return truncate(comment.Content, 50) + "..."
}

//Excerpt returns comment excerpt, 100 char long
func (comment *Comment) Excerpt() string {
	return truncate(comment.Content, 20)
}

//GetCommentsByArticleID returns a slice of published comments, associated with the given article
func GetComments(articleID uint) (comments []Comment) {
	var list []Comment
	//published == false is required ;D
	db.Where("published = ? AND answer <> ? AND article_id = ? AND author_city = ?",
		false, "", articleID, "Москва").Order("id desc").Find(&comments)
	db.Where("published = ? AND answer <> ? AND article_id = ? AND author_city <> ?",
		false, "", articleID, "Москва").Order("id desc").Find(&list)
	comments = append(comments, list...)
	return
}

//GetTopCommentsByArticleID returns a slice of top (latest, or rated) published comments, associated with given article
func GetTopComments(articleID uint) (comments []Comment) {
	db.Where("published = ? AND article_id = ? AND author_city = ?",
		true, articleID, "Москва").Order("id desc").Limit(10).Find(&comments)
	return
}
