package domain

import "time"

type Comment struct {
	Id            string     `json:"id,omitempty"`
	Uid           string     `json:"uid,omitempty"`
	CommentText   string     `json:"commentText,omitempty"`
	AuthorId      string     `json:"authorId,omitempty"`
	ParentId      string     `json:"parentId,omitempty"`
	PostId        string     `json:"postId,omitempty"`
	CommentType   string     `json:"commentType,omitempty"`
	CommentsCount int        `json:"commentsCount,omitempty"`
	MyLike        bool       `json:"myLike,omitempty"`
	LikeAmount    int        `json:"likeAmount,omitempty"`
	TimeChanged   *time.Time `json:"timeChanged,omitempty"`
	Time          *time.Time `json:"time,omitempty"`
	ImagePath     string     `json:"imagePath,omitempty"`
	IsBlocked     bool       `json:"isBlocked,omitempty"`
	IsDeleted     bool       `json:"isDeleted,omitempty"`
	DType         []string   `json:"dgraph.type,omitempty"`
	Comments      []Comment  `json:"comments,omitempty"`
}