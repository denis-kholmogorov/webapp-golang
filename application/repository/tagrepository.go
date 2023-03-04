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

func (r TagRepository) CreateIfNotExists(names []string) (t []domain.Tag, e error) {

	var err error
	var tags []domain.Tag
	for _, name := range names {
		tag, err := r.getByNane(name)
		if err != nil {
			break
		}
		tags = append(tags, tag)
	}
	if err != nil {
		log.Printf("TagRepository:FindAll() Error query %s", err)
		return nil, fmt.Errorf("TagRepository:FindAll() Error query %s", err)
	}
	return tags, nil

}

func (r TagRepository) getByNane(name string) (t domain.Tag, e error) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	var vars *api.Response
	var err error
	vars, err = txn.QueryWithVars(ctx, getTagByName, map[string]string{"$name": name}) //TODO посмотреть другой вариант
	response := domain.TagList{}
	err = json.Unmarshal(vars.Json, &response)
	if err != nil {
		log.Printf("TagRepository:getByNane() Error Unmarshal %s", err)
		return t, fmt.Errorf("TagRepository:getByNane() Error Unmarshal %s", err)
	}
	if len(response.List) == 0 {
		tag := domain.Tag{Name: name, DType: []string{"Tag"}}
		return tag, nil
	}
	return response.List[0], nil

}

var getTagByName = `query Tag($name: string)
{
  tagList(func: eq(name, $name))  {
	id: uid
	name
  }
}
`
