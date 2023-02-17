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
	DType         []string   `json:"dgraph.type,omitempty"`
	Tags          []string   `json:"tags,omitempty"`
	CreatedOn     *time.Time `json:"time" time_format:"2006-01-02 15:04:05.99Z07:00"`
	UpdateOn      *time.Time `json:"timeChanged" time_format:"2006-01-02 15:04:05.99Z07:00"`
	PublishDate   *time.Time `json:"publishDate,omitempty" time_format:"2006-01-02 15:04:05.99Z07:00"`
}

type Posts struct {
	List         []Post `json:"content,omitempty"`
	TotalElement int    `json:"totalElement,omitempty"`
	TotalPages   int    `json:"totalPages,omitempty"`
	Number       int    `json:"number"`
	Size         int    `json:"size,omitempty"`
}

type PostResponse struct {
	Posts []Posts `json:"posts,omitempty"`
}
