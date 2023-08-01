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

func (r PostRepository) GetAll(searchDto dto.PostSearchDto, currentUserId string) (post *dto.PageResponse[domain.Post], e error) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	var vars *api.Response
	var err error
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	variables["$first"] = strconv.Itoa(searchDto.Size)
	variables["$offset"] = strconv.Itoa(searchDto.Size * utils.GetPageNumber(&searchDto))
	variables["$dateFrom"] = utils.ConvSecToDateString(searchDto.DateFrom)
	variables["$dateTo"] = utils.GetCurrentTimeString()

	if len(searchDto.AccountIds) == 0 {
		searchDto.AccountIds = append(searchDto.AccountIds, currentUserId)
	}

	if searchDto.WithFriends {
		ids, err := friendRepo.GetMyFriends(currentUserId)
		if err != nil {
			log.Printf("PostRepository:FindAll() Error add friends %s", err)
			return nil, fmt.Errorf("PostRepository:FindAll() Error add friends %s", err)
		}
		for _, id := range ids {
			searchDto.AccountIds = append(searchDto.AccountIds, id)
		}
	}

	if len(searchDto.Text) == 0 {
		ids := strings.Join(searchDto.AccountIds, ",")
		vars, err = txn.QueryWithVars(ctx, fmt.Sprintf(getAllPosts, ids), variables)
	} else {
		variables["$text"] = searchDto.Text
		variables["$author"] = fmt.Sprintf(regexAuthor, searchDto.Author)
		vars, err = txn.QueryWithVars(ctx, getAllPostsByText, variables)
	}

	if err != nil {
		log.Printf("PostRepository:FindAll() Error query %s", err)
		return nil, fmt.Errorf("PostRepository:FindAll() Error query %s", err)
	}
	response := dto.PageResponse[domain.Post]{}
	replacer := strings.NewReplacer("@groupby", "reactions", "\"myLike\":0", "\"myLike\":false", "\"myLike\":1", "\"myLike\":true")
	jsonReplace := []byte(replacer.Replace(string(vars.Json)))
	err = json.Unmarshal(jsonReplace, &response)
	if err != nil {
		log.Printf("PostRepository:FindAll() Error Unmarshal %s", err)
		return nil, fmt.Errorf("PostRepository:FindAll() Error Unmarshal %s", err)
	}
	updateLikes(response.Content)
	response.SetPage(searchDto.Size, searchDto.Page)
	return &response, nil

}

func updateLikes(posts []domain.Post) {

	for i, _ := range posts {
		if len(posts[i].RowReactions) > 0 {
			posts[i].Reactions = posts[i].RowReactions[0].Reactions
			posts[i].RowReactions = nil
		}
		if len(posts[i].RowMyReaction) > 0 {
			posts[i].MyReaction = posts[i].RowMyReaction[0].MyReaction
			posts[i].RowMyReaction = nil
		}
	}
}

func (r PostRepository) Create(post *domain.Post, authorId string, tagIds []string) (*domain.Post, error) {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	timeNow := time.Now().UTC()
	post.Uid = "_:post"
	post.DType = []string{"Post"}
	post.CreatedOn = &timeNow
	post.UpdateOn = &timeNow
	post.AuthorId = authorId
	if post.PublishDate != nil {
		post.Type = "QUEUED"
	} else {
		post.Type = "POSTED"
	}
	author := domain.Account{Uid: authorId, Posts: []domain.Post{*post}}
	authorm, err := json.Marshal(author)
	if err != nil {
		log.Printf("PostRepository:save() Error marhalling post %s", err)
		return nil, fmt.Errorf("PostRepository:Create() Error marhalling post %s", err)
	}
	mu, err := txn.Mutate(ctx, &api.Mutation{SetJson: authorm, CommitNow: true})
	if err != nil {
		log.Printf("PostRepository:save() Error mutate %s", err)
		return nil, fmt.Errorf("PostRepository:Create() Error mutate %s", err)
	}
	post.Id = mu.Uids["post"]
	AddNodesToEdge(txn, ctx, post.Id, "tags", tagIds, true)
	return post, nil
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

func (r PostRepository) GetAllComments(pageDto dto.PageRequest, parentId string, currentUserId string) (post *dto.PageResponse[domain.Comment], e error) {
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
	response := dto.PageResponse[domain.Comment]{}
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

func (r PostRepository) UpdatePostQueue() {
	ctx := context.Background()
	variables := make(map[string]string)
	variables["$type"] = "QUEUED"
	variables["$currentTime"] = utils.GetCurrentTimeString()

	mu := &api.Mutation{
		SetNquads: []byte(updateStatusPost),
	}

	req := &api.Request{
		Query:     getAllQueuePosts,
		Mutations: []*api.Mutation{mu},
		Vars:      variables,
		CommitNow: true,
	}
	_, err := r.conn.NewTxn().Do(ctx, req)
	if err != nil {
		log.Printf("PostRepository: UpdatePostQueue() Error Unmarshal %s", err)
	}
}

var updateStatusPost = `uid(A) <status> "POSTED" .`

var getAllPostsByText = `query Posts($currentUserId: string, $text: string, $author: string, $dateFrom: string, $dateTo: string, $first: int, $offset: int)
{
  var(func: anyoftext(postText, $text)) @filter(eq(isDeleted, false) and between(time, $dateFrom, $dateTo))  {
    A as uid
  	}

  var(func: type(Account)) @filter(eq(isDeleted, false) and regexp(firstName, $author) or regexp(lastName, $author)){
    posts @filter(eq(isDeleted, false) and between(time, $dateFrom, $dateTo)) {
 		B as uid
    	}
	}

    content(func: uid(A,B), first: $first, offset: $offset, orderdesc: time) @normalize  {
	id:uid
	postText: postText
	authorId: authorId
	title: title
	time: time 
	timeChanged: timeChanged
	type: type
	isDeleted: isDeleted
	isBlocked: isBlocked
	imagePath: imagePath
	myLike: count(likes @filter(eq(authorId,$currentUserId)))
    rowMyReaction: likes @filter(eq(authorId,$currentUserId)){
		myReaction:reactionType
    }
    rowReactions: likes @groupby(reactionType){
        count(uid)
    }
	likeAmount: count(likes)
	commentsCount: count(comments @filter(eq(isDeleted, false)))
	tags: tags {
        uid
		name
	}
    }
    count(func: anyoftext(postText, $text)){
		totalElement:count(uid)
  }
}`

var getAllPosts = `query Posts($currentUserId: string, $dateFrom: string, $dateTo: string, $first: int, $offset: int)
{
var(func: uid(%s)) @filter(eq(isDeleted, false))  {
  content:posts @filter(eq(isDeleted, false) and between(time, $dateFrom, $dateTo)){
	A as uid
  }
}
  content(func: uid(A), orderdesc: time, first: $first, offset: $offset) {
    id:uid
	postText: postText
	authorId: authorId
	title: title
	time: time
	timeChanged: timeChanged
	type: type
	isDeleted: isDeleted
	isBlocked: isBlocked
	imagePath: imagePath
	myLike: count(likes @filter(eq(authorId,$currentUserId)))
    rowMyReaction: likes @filter(eq(authorId,$currentUserId)){
		myReaction:reactionType
    }
    rowReactions: likes @groupby(reactionType){
       count(uid)
    }
    tags {
	  id: uid
      name
    }
	likeAmount: count(likes)
	commentsCount: count(comments @filter(eq(isDeleted, false)))
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

var getAllQueuePosts = `query getQueuePosts($currentTime: string, $type: string )
{
  posts(func: type(Post)) @filter(eq(type,$type) and lt(publishDate, $currentTime))  {
    A as uid
  }
}
`
