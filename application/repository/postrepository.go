package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"log"
	"math"
	"strconv"
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

func (r PostRepository) Create(post *domain.Post, authorId string) (*string, error) {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	timeNow := time.Now().UTC()
	post.DType = []string{"Post"}
	post.CreatedOn = &timeNow
	post.UpdateOn = &timeNow
	author := domain.Account{Uid: authorId, Posts: []domain.Post{*post}}
	authorm, err := json.Marshal(author)
	if err != nil {
		log.Printf("PostRepository:save() Error marhalling post %s", err)
		return nil, fmt.Errorf("PostRepository:Create() Error marhalling post %s", err)
	}
	mutate, err := txn.Mutate(ctx, &api.Mutation{SetJson: authorm, CommitNow: true})
	if err != nil {
		log.Printf("PostRepository:save() Error mutate %s", err)
		return nil, fmt.Errorf("PostRepository:Create() Error mutate %s", err)
	}
	postId := mutate.Uids[""]
	//if len(postId) == 0 {
	//	log.Printf("PostRepository:save() capthcaId not found")
	//	return nil, fmt.Errorf("PostRepository:Create() capthcaId not found")
	//}
	return &postId, nil
}

func (r PostRepository) GetAll(dto dto.PostSearchDto) (post *domain.Posts, e error) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	var vars *api.Response
	var err error
	variables := make(map[string]string)
	variables["$first"] = strconv.Itoa(dto.Size)
	variables["$offset"] = strconv.Itoa(dto.Size * utils.GetPageNumber(&dto))
	if len(dto.Text) == 0 {
		variables["$accountId"] = dto.AuthorId
		vars, err = txn.QueryWithVars(ctx, getAllPosts, variables)
	} else {
		variables["$text"] = strconv.Itoa(dto.Size)
		vars, err = txn.QueryWithVars(ctx, getAllPostsByText, variables)
	}

	if err != nil {
		log.Printf("GeoRepository:FindAll() Error query %s", err)
		return nil, fmt.Errorf("GeoRepository:FindAll() Error query %s", err)
	}

	response := domain.PostResponse{}
	err = json.Unmarshal(vars.Json, &response)
	response.Posts[0].Size = dto.Size
	response.Posts[0].TotalPages = int(math.Ceil(float64(response.Posts[0].TotalElement) / float64(dto.Size)))
	response.Posts[0].Number = dto.Page + 1
	for i := range response.Posts[0].List {
		response.Posts[0].List[i].AuthorId = dto.AuthorId
	}
	return &response.Posts[0], nil
}

var getAllPostsByText = `query Posts($text: string, $first: int, $offset: int)
{
    posts(func: anyoftext(postText, $text) ) {
    content:posts(orderdesc: time)(first: $first, offset: $offset){
	id:uid
	postText
	title
	time
	tags
	timeChanged
	type
	isDeleted
	isBlocked
	imagePath
    }
    totalElement:count(posts)
  }
}`

var getAllPosts = `query Posts($accountId: string, $first: int, $offset: int)
{
    posts(func: uid($accountId) ) {
    content:posts(orderdesc: time)(first: $first, offset: $offset){
	id:uid
	postText
	title
	time
	tags
	timeChanged
	type
	isDeleted
	isBlocked
	imagePath
    }
    totalElement:count(posts)
  }
}`
