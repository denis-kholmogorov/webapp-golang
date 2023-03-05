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

var tagRepo *TagRepository
var isInitializedtagRepo bool

type TagRepository struct {
	conn *dgo.Dgraph
}

func GetTagRepository() *TagRepository {
	if !isInitializedtagRepo {
		tagRepo = &TagRepository{}
		tagRepo.conn = GetDGraphConn().connection
		isInitializedtagRepo = true
	}
	return tagRepo
}

func (r TagRepository) Update(post *domain.Post) error {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	var err error
	err = deleteEdgesTags(ctx, txn, post.Id)
	if err != nil {
		err = txn.Discard(ctx)
		return err
	}
	for i, _ := range post.Tags {
		err = r.getByNane(ctx, txn, &post.Tags[i])
		if err != nil {
			err = txn.Discard(ctx)
			break
		}
	}
	if err != nil {
		log.Printf("TagRepository:FindAll() Error query %s", err)
		return fmt.Errorf("TagRepository:FindAll() Error query %s", err)
	}
	err = txn.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r TagRepository) Create(post *domain.Post) error {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	var err error
	for i, _ := range post.Tags {
		err := r.getByNane(ctx, txn, &post.Tags[i])
		if err != nil {
			break
		}
	}
	if err != nil {
		log.Printf("TagRepository:FindAll() Error query %s", err)
		return fmt.Errorf("TagRepository:FindAll() Error query %s", err)
	}
	return nil
}

func (r TagRepository) getByNane(ctx context.Context, txn *dgo.Txn, tag *domain.Tag) error {
	var vars *api.Response
	var err error
	vars, err = txn.QueryWithVars(ctx, getTagByName, map[string]string{"$name": tag.Name}) //TODO посмотреть другой вариант
	response := domain.TagList{}
	err = json.Unmarshal(vars.Json, &response)
	if err != nil {
		log.Printf("TagRepository:getByNane() Error Unmarshal %s", err)
		return fmt.Errorf("TagRepository:getByNane() Error Unmarshal %s", err)
	}
	tag.DType = []string{"Tag"}
	if len(response.List) != 0 {
		tag.Id = response.List[0].Id
		return nil
	}
	return nil

}

func deleteEdgesTags(ctx context.Context, txn *dgo.Txn, id string) error {
	mu := &api.Mutation{}
	dgo.DeleteEdges(mu, id, "tags")
	_, err := txn.Mutate(ctx, mu)
	return err
}

var getTagByName = `query Tag($name: string)
{
  tagList(func: eq(name, $name))  {
	id: uid
	name
  }
}
`
