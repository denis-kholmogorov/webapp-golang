package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"log"
	"web/application/domain"
)

var likeRepo *LikeRepository
var isInitializedlikeRepo bool

type LikeRepository struct {
	conn *dgo.Dgraph
}

func GetLikeRepository() *LikeRepository {
	if !isInitializedlikeRepo {
		likeRepo = &LikeRepository{}
		likeRepo.conn = GetDGraphConn().connection
		isInitializedlikeRepo = true
	}
	return likeRepo
}

func (r LikeRepository) CreateLike(parentId string, authorId string) (bool, error) {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	exist, err := isLikeExist(ctx, txn, parentId, authorId)
	if err != nil {
		txn.Discard(ctx)
		return false, err
	}
	if exist {
		txn.Discard(ctx)
		return false, fmt.Errorf("like has been installed")
	}

	likem, err := json.Marshal(domain.Like{Uid: "_:like", DType: []string{"Like"}, AuthorId: authorId})
	if err != nil {
		txn.Discard(ctx)
		log.Printf("LikeRepository:save() Error marhalling post %s", err)
		return false, fmt.Errorf("LikeRepository:Create() Error marhalling post %s", err)
	}
	mutate, err := txn.Mutate(ctx, &api.Mutation{SetJson: likem})
	if err != nil {
		txn.Discard(ctx)
		log.Printf("LikeRepository:save() Error mutate %s", err)
		return false, fmt.Errorf("LikeRepository:Create() Error mutate %s", err)
	}
	likeId := mutate.GetUids()["like"]
	err = AddEdge(txn, ctx, parentId, "likes", likeId, true)
	if err != nil {
		return false, err
	}
	return exist, err
}

func (r LikeRepository) DeleteLike(parentId string, authorId string) (bool, error) {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	like, err := getLikeByAuthorAndPostId(ctx, txn, parentId, authorId)
	if err != nil {
		txn.Discard(ctx)
		return false, err
	}

	likem, err := json.Marshal(like)
	if err != nil {
		txn.Discard(ctx)
		return false, fmt.Errorf("LikeRepository:Create() Error marhalling post %s", err)
	}
	err = RemoveEdge(txn, ctx, parentId, "likes", like.Uid, false)
	if err != nil {
		txn.Discard(ctx)
		return false, fmt.Errorf("LikeRepository:DeleteLike() Error mutate DeleteEdges %s", err)
	}
	_, err = txn.Mutate(ctx, &api.Mutation{DeleteJson: likem, CommitNow: true})
	if err != nil {
		log.Printf("LikeRepository:DeleteLike() Error mutate delete like %s", err)
		return false, fmt.Errorf("LikeRepository:Create() Error mutate  delete like %s", err)
	}
	return true, err
}

func getLikeByAuthorAndPostId(ctx context.Context, txn *dgo.Txn, parentId string, authorId string) (*domain.Like, error) {
	variables := make(map[string]string)
	variables["$parentId"] = parentId
	variables["$authorId"] = authorId
	vars, err := txn.QueryWithVars(ctx, getLikeByPostIdAndAuthorId, variables)
	if err != nil {
		return nil, err
	}
	like := domain.LikeList{}
	err = json.Unmarshal(vars.Json, &like)
	if err != nil {
		return nil, err
	}
	if len(like.List) == 0 {
		return nil, fmt.Errorf("like of authorId=%s not found for delete parentId=%s", authorId, parentId)
	}
	return &like.List[0], err
}

func isLikeExist(ctx context.Context, txn *dgo.Txn, parentId string, authorId string) (bool, error) {
	variables := make(map[string]string)
	variables["$parentId"] = parentId
	variables["$authorId"] = authorId
	vars, err := txn.QueryWithVars(ctx, existLike, variables)
	if err != nil {
		return false, err
	}
	exist := domain.LikeExist{}
	err = json.Unmarshal(vars.Json, &exist)

	return exist.Exists[0].Count > 0, err
}

var existLike = `query Exists($parentId: string, $authorId: string)
{
exists(func: uid($parentId)){
	count:count(likes @filter(eq(authorId,$authorId)))
}
}
`

var getLikeByPostIdAndAuthorId = `query Exists($parentId: string, $authorId: string)
{
var(func: uid($parentId)){
  likes @filter(eq(authorId, $authorId)){
      A as uid
    }
  }
}
  {
    likeList(func: uid(A)){
    uid
    authorId
  }
}
`
