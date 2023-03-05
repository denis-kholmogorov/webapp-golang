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

func (r TagRepository) CreateIfNotExists(post *domain.Post) error {

	var err error
	for i, _ := range post.Tags {
		err := r.getByNane(&post.Tags[i])
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

func (r TagRepository) getByNane(tag *domain.Tag) error {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
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

var getTagByName = `query Tag($name: string)
{
  tagList(func: eq(name, $name))  {
	id: uid
	name
  }
}
`
