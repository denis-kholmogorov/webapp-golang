package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"log"
	"web/application/domain"
	"web/application/errorhandler"
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

func (r LikeRepository) CreateLike(parentId string, like *domain.Like) {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	likes := LikeByParentIdAndAuthorId(ctx, txn, parentId, like.AuthorId)

	if len(likes.List) > 0 {
		if likes.List[0].ReactionType != like.ReactionType {
			updateLikeReaction(ctx, txn, likes.List[0].Uid, like)
		}
		return
	}

	like.Uid = "_:like"
	like.DType = []string{"Like"}
	likem, err := json.Marshal(like)
	if err != nil {
		txn.Discard(ctx)
		log.Printf("LikeRepository:save() Error marhalling post %s", err)
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("LikeRepository:CreateLike() Error marshal like %s", err)})
	}
	mutate, err := txn.Mutate(ctx, &api.Mutation{SetJson: likem})
	if err != nil {
		txn.Discard(ctx)
		log.Printf("LikeRepository:save() Error mutate %s", err)
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("likeRepository:CreateLike() Error update parent %s", err)})
	}
	likeId := mutate.GetUids()["like"]
	err = AddEdge(txn, ctx, parentId, "likes", likeId, true)
	if err != nil {
		txn.Discard(ctx)
		panic(errorhandler.DbError{Message: fmt.Sprintf("likeRepository:CreateLike() Error update parent %s", err)})
	}
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
	err = DeleteEdge(txn, ctx, parentId, "likes", like.Uid, false)
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

func updateLikeReaction(ctx context.Context, txn *dgo.Txn, idLike string, like *domain.Like) {
	err := UpdateNodeFields(txn, ctx, idLike, map[string]string{"reactionType": like.ReactionType}, true)
	if err != nil {
		txn.Discard(ctx)
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("LikeRepository:LikeByParentIdAndAuthorId() Error get likes %s", err)})
	}
}

func getLikeByAuthorAndPostId(ctx context.Context, txn *dgo.Txn, parentId string, authorId string) (*domain.Like, error) {
	variables := make(map[string]string)
	variables["$parentId"] = parentId
	variables["$authorId"] = authorId
	vars, err := txn.QueryWithVars(ctx, existLike, variables)
	if err != nil {
		panic(errorhandler.DbError{Message: fmt.Sprintf("LikeRepository:getLikeByAuthorAndPostId() Error get likes %s", err)})
	}
	like := domain.LikeList{}
	err = json.Unmarshal(vars.Json, &like)
	if err != nil {
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("LikeRepository:getLikeByAuthorAndPostId() Error get likes %s", err)})
	}
	if len(like.List) == 0 {
		return nil, fmt.Errorf("like of authorId=%s not found for delete parentId=%s", authorId, parentId)
	}
	return &like.List[0], err
}

func LikeByParentIdAndAuthorId(ctx context.Context, txn *dgo.Txn, parentId string, authorId string) *domain.LikeList {
	variables := make(map[string]string)
	variables["$parentId"] = parentId
	variables["$authorId"] = authorId
	vars, err := txn.QueryWithVars(ctx, existLike, variables)
	if err != nil {
		txn.Discard(ctx)
		panic(errorhandler.DbError{Message: fmt.Sprintf("LikeRepository:LikeByParentIdAndAuthorId() Error get likes %s", err)})
	}
	likes := domain.LikeList{}
	err = json.Unmarshal(vars.Json, &likes)
	if err != nil {
		txn.Discard(ctx)
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("LikeRepository:LikeByParentIdAndAuthorId() Error Unmarshal %s", err)})
	}
	return &likes
}

var existLike = `query Exists($parentId: string, $authorId: string)
{
likes(func: uid($parentId)) @normalize{
	likes: likes @filter(eq(authorId,$authorId)){
		uid:uid
		reactionType: reactionType
	}
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
    reactionType
  }
}
`
