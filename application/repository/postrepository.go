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

func (r PostRepository) GetAll(searchDto dto.PostSearchDto) (post *dto.PageResponse, e error) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	var vars *api.Response
	var err error
	variables := make(map[string]string)
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

var getAllPostsByText = `query Posts($text: string, $first: int, $offset: int)
{
    content(func: anyoftext(postText, $text), first: $first, offset: $offset, orderdesc: time) {
	id:uid
	postText
	authorId
	title
	time
	tags
	timeChanged
	type
	isDeleted
	isBlocked
	imagePath
    }
    count(func: anyoftext(postText, $text)){
		totalElement:count(uid)
  }
}`

var getAllPosts = `query Posts($accountId: string, $first: int, $offset: int)
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
	tags
	timeChanged
	type
	isDeleted
	isBlocked
	imagePath
  }
  count(func: uid(A)){
		totalElement:count(uid)
  }
}
`
