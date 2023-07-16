package domain

import "time"

type Post struct {
	Id            string     `json:"id,omitempty"`
	Uid           string     `json:"uid,omitempty"`
	Title         string     `json:"title"`
	Type          string     `json:"type,omitempty"`
	AuthorId      string     `json:"authorId,omitempty"`
	PostText      string     `json:"postText"`
	CommentsCount int        `json:"commentsCount"`
	MyLike        int        `json:"myLike,omitempty"`
	LikeAmount    int        `json:"likeAmount,omitempty"`
	ImagePath     string     `json:"imagePath"`
	IsDeleted     bool       `json:"isDeleted"`
	DType         []string   `json:"dgraph.type,omitempty"`
	Tags          []Tag      `json:"tags,omitempty"` //TODO проверить на null бек
	Likes         []Like     `json:"likes,omitempty"`
	Comments      []Comment  `json:"comments,omitempty"`
	CreatedOn     *time.Time `json:"time" time_format:"2006-01-02 15:04:05.99Z07:00"`
	UpdateOn      *time.Time `json:"timeChanged" time_format:"2006-01-02 15:04:05.99Z07:00"`
	PublishDate   *time.Time `json:"publishDate,omitempty" time_format:"2006-01-02 15:04:05.99Z07:00"`
}

type PostList struct {
	List []Post `json:"postList,omitempty"`
}
