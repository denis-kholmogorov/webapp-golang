package domain

import "time"

type Post struct {
	Uid           string     `json:"id,omitempty"`
	Title         string     `json:"title"`
	Type          string     `json:"type,omitempty"`
	AuthorId      string     `json:"authorId,omitempty"`
	PostText      string     `json:"postText"`
	CommentsCount int        `json:"commentsCount,omitempty"`
	LikeAmount    int        `json:"likeAmount,omitempty"`
	ImagePath     string     `json:"imagePath"`
	IsDeleted     bool       `json:"isDeleted"`
	DType         []string   `json:"dgraph.type,omitempty"`
	Tags          []string   `json:"tags,omitempty"`
	CreatedOn     *time.Time `json:"time" time_format:"2006-01-02 15:04:05.99Z07:00"`
	UpdateOn      *time.Time `json:"timeChanged" time_format:"2006-01-02 15:04:05.99Z07:00"`
	PublishDate   *time.Time `json:"publishDate,omitempty" time_format:"2006-01-02 15:04:05.99Z07:00"`
}
