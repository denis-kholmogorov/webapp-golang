package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"log"
	"strconv"
	"strings"
	"time"
	"web/application/domain"
	"web/application/dto"
	"web/application/utils"
)

var postRepo *PostRepository
var isInitializedpostRepo bool

type PostRepository struct {
	conn *dgo.Dgraph
}

func GetPostRepository() *PostRepository {
	if !isInitializedpostRepo {
		postRepo = &PostRepository{}
		postRepo.conn = GetDGraphConn().connection
		isInitializedpostRepo = true
	}
	return postRepo
}

func (r PostRepository) GetAll(searchDto dto.PostSearchDto, currentUserId string) (post *dto.PageResponse, e error) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	var vars *api.Response
	var err error
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	variables["$first"] = strconv.Itoa(searchDto.Size)
	variables["$offset"] = strconv.Itoa(searchDto.Size * utils.GetPageNumber(&searchDto))
	if len(searchDto.Text) == 0 {
		variables["$accountId"] = strings.Join(searchDto.AccountIds, ",")
		vars, err = txn.QueryWithVars(ctx, getAllPosts, variables)
	} else {
		variables["$text"] = searchDto.Text
		vars, err = txn.QueryWithVars(ctx, getAllPostsByText, variables)
	}

	if err != nil {
		log.Printf("PostRepository:FindAll() Error query %s", err)
		return nil, fmt.Errorf("PostRepository:FindAll() Error query %s", err)
	}
	response := dto.PageResponse{}
	err = json.Unmarshal(vars.Json, &response)
	if err != nil {
		log.Printf("PostRepository:FindAll() Error Unmarshal %s", err)
		return nil, fmt.Errorf("PostRepository:FindAll() Error Unmarshal %s", err)
	}

	response.SetPage(searchDto.Size, searchDto.Page)
	return &response, nil

}

func (r PostRepository) Create(post *domain.Post, authorId string) (*string, error) {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	timeNow := time.Now().UTC()
	post.DType = []string{"Post"}
	post.CreatedOn = &timeNow
	post.UpdateOn = &timeNow
	post.AuthorId = authorId
	author := domain.Account{Uid: authorId, Posts: []domain.Post{*post}}
	authorm, err := json.Marshal(author)
	if err != nil {
		log.Printf("PostRepository:save() Error marhalling post %s", err)
		return nil, fmt.Errorf("PostRepository:Create() Error marhalling post %s", err)
	}
	_, err = txn.Mutate(ctx, &api.Mutation{SetJson: authorm, CommitNow: true})
	if err != nil {
		log.Printf("PostRepository:save() Error mutate %s", err)
		return nil, fmt.Errorf("PostRepository:Create() Error mutate %s", err)
	}

	return nil, nil
}

func (r PostRepository) Update(post *domain.Post, tagIds []string) error {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	timeNow := time.Now().UTC()
	post.DType = []string{"Post"}
	post.Uid = post.Id
	fields := map[string]string{
		"timeChanged": timeNow.String(),
		"title":       post.Title,
		"postText":    post.PostText,
	}
	err := UpdateNodeFields(txn, ctx, post.Id, fields, false)
	if err != nil {
		return err
	}
	return AddNodesToEdge(txn, ctx, post.Id, "tags", tagIds, true)
}

func (r PostRepository) GetAllComments(pageDto dto.PageRequest, parentId string, currentUserId string) (post *dto.PageResponse, e error) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	var vars *api.Response
	var err error
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	variables["$first"] = strconv.Itoa(pageDto.Size)
	variables["$offset"] = strconv.Itoa(pageDto.Size * utils.GetPageNumber(&pageDto))
	variables["$parentId"] = parentId
	vars, err = txn.QueryWithVars(ctx, getAllComments, variables)

	if err != nil {
		log.Printf("PostRepository:FindAll() Error query %s", err)
		return nil, fmt.Errorf("PostRepository:FindAll() Error query %s", err)
	}
	response := dto.PageResponse{}
	err = json.Unmarshal(vars.Json, &response)
	if err != nil {
		log.Printf("PostRepository:FindAll() Error Unmarshal %s", err)
		return nil, fmt.Errorf("PostRepository:FindAll() Error Unmarshal %s", err)
	}
	response.SetPage(pageDto.Size, pageDto.Page)
	return &response, nil
}

func (r PostRepository) CreateComment(comment domain.Comment, postId string, authorId string) (c *string, e error) {
	var err error
	var marshal []byte
	ctx := context.Background()
	txn := r.conn.NewTxn()
	timeNow := time.Now().UTC()
	comment.Uid = "_:comment"
	comment.DType = []string{"Comment"}
	comment.AuthorId = authorId
	comment.Time = &timeNow
	comment.PostId = postId
	comment.TimeChanged = &timeNow
	comment.IsDeleted = false
	if len(comment.ParentId) > 0 {
		comment.CommentType = "COMMENT"
		postId = comment.ParentId
		marshal, err = json.Marshal(comment)
	} else {
		comment.CommentType = "POST"
		marshal, err = json.Marshal(comment)
	}

	if err != nil {
		log.Printf("PostRepository:save() Error marhalling post %s", err)
		return nil, fmt.Errorf("PostRepository:Create() Error marhalling post %s", err)
	}
	mutate, err := txn.Mutate(ctx, &api.Mutation{SetJson: marshal})
	if err != nil {
		log.Printf("PostRepository:save() Error mutate %s", err)
		return nil, fmt.Errorf("PostRepository:Create() Error mutate %s", err)
	}
	err = AddEdge(txn, ctx, postId, "comments", mutate.GetUids()["comment"], true)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r PostRepository) UpdateComment(comment domain.Comment) error {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	timeNow := time.Now().UTC()
	fields := map[string]string{"timeChanged": timeNow.String(), "commentText": comment.CommentText}
	return UpdateNodeFields(txn, ctx, comment.Id, fields, true)
}

func (r PostRepository) Delete(postId string) error {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	return Delete(txn, ctx, postId, true)
}

var getAllPostsByText = `query Posts($currentUserId: string, $text: string, $first: int, $offset: int)
{
    content(func: anyoftext(postText, $text), first: $first, offset: $offset, orderdesc: time) @filter(eq(isDeleted, false)) {
	id:uid
	postText
	authorId
	title
	time 
	timeChanged
	type
	isDeleted
	isBlocked
	imagePath
	myLike: count(likes @filter(eq(authorId,$currentUserId)))
	likeAmount: count(likes)
	commentsCount: count(comments @filter(eq(isDeleted, false)))
	tags {
		name
	}
    }
    count(func: anyoftext(postText, $text)){
		totalElement:count(uid)
  }
}`

var getAllPosts = `query Posts($currentUserId: string, $accountId: string, $first: int, $offset: int)
{
var(func: uid($accountId)) @filter(eq(isDeleted, false))  {
  content:posts @filter(eq(isDeleted, false)){
	A as uid
  }
}
  content(func: uid(A), orderdesc: time, first: $first, offset: $offset)  {
    id:uid
	postText
	authorId
	title
	time
	timeChanged
	type
	isDeleted
	isBlocked
	imagePath
	myLike: count(likes @filter(eq(authorId,$currentUserId)))
	likeAmount: count(likes)
	commentsCount: count(comments @filter(eq(isDeleted, false))) 
	tags {
	name
}
  }
  count(func: uid(A)){
		totalElement:count(uid)
  }
}
`
var getAllComments = `query Comments($parentId: string, $currentUserId: string, $first: int, $offset: int)
{
var(func: uid($parentId)) {
	A as comments
}
  content(func: uid(A), orderdesc: timeChanged, first: $first, offset: $offset) @filter(eq(isDeleted, false)) {
	id: uid
	commentText
	authorId
	parentId
	postId
	commentType
	commentsCount: count(comments  @filter(eq(isDeleted, false)))
    myLike: count(likes @filter(eq(authorId,$currentUserId)))
	likeAmount: count(likes)
	timeChanged
	time
	imagePath
	isBlocked
	isDeleted
  }
 count(func: uid(A), orderdesc: timeChanged){
	totalElement: count(uid)
  }
}
`
