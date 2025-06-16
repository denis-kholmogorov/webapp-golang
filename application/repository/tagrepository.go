package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"log"
	"web/application/domain"
	"web/application/dto"
	"web/application/errorhandler"
)

const regexTagName = "/%s.*/"

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

func (r TagRepository) FindAll(searchDto dto.TagSearchDto) *[]domain.Tag {
	variables := make(map[string]string)
	variables["$search"] = fmt.Sprintf(regexTagName, searchDto.Name)
	txn := r.conn.NewReadOnlyTxn()
	vars, err := txn.QueryWithVars(context.Background(), findAllTag, variables)
	if err != nil {
		panic(errorhandler.DbError{Message: fmt.Sprintf("TagRepository:FindAll() Error query %s", err)})
	}
	response := domain.TagList{}
	err = json.Unmarshal(vars.Json, &response)
	return &response.List
}

func (r TagRepository) Create(post *domain.Post) error {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	var err error
	for i, _ := range post.Tags {
		err = r.getByName(ctx, txn, &post.Tags[i])
		if err != nil {
			break
		}
	}
	if err != nil {
		err = txn.Discard(ctx)
		log.Printf("TagRepository:FindAll() Error query %s", err)
		return fmt.Errorf("TagRepository:FindAll() Error query %s", err)
	}
	err = txn.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r TagRepository) Update(post *domain.Post) ([]string, error) {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	var err error
	err = deleteEdgesTags(ctx, txn, post.Id)
	if err != nil {
		err = txn.Discard(ctx)
		return nil, err
	}
	var tagIds []string
	for i, _ := range post.Tags {
		id, err := r.getIdByName(ctx, txn, post.Tags[i].Name)
		if err != nil {
			break
		}
		tagIds = append(tagIds, id)
	}
	if err != nil {
		err = txn.Discard(ctx)
		log.Printf("TagRepository:FindAll() Error query %s", err)
		return nil, fmt.Errorf("TagRepository:FindAll() Error query %s", err)
	}
	err = txn.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return tagIds, err
}

func (r TagRepository) getIdByName(ctx context.Context, txn *dgo.Txn, tagName string) (string, error) {
	var vars *api.Response
	var err error
	vars, err = txn.QueryWithVars(ctx, getTagByName, map[string]string{"$name": tagName}) //TODO посмотреть другой вариант
	response := domain.TagList{}
	err = json.Unmarshal(vars.Json, &response)
	if err != nil {
		log.Printf("TagRepository:getIdByName() Error Unmarshal %s", err)
		return "", fmt.Errorf("TagRepository:getIdByName() Error Unmarshal %s", err)
	}
	if len(response.List) != 0 {
		return response.List[0].Id, nil
	}
	tag := domain.Tag{Uid: "_:tag", Name: tagName, DType: []string{"Tag"}}
	tagm, err := json.Marshal(tag)
	if err != nil {
		log.Printf("TagRepository:getIdByName() Error marhalling tagName %s", err)
		return "", fmt.Errorf("TagRepository:getIdByName() Error marhalling tagName %s", err)
	}
	mutate, err := txn.Mutate(ctx, &api.Mutation{SetJson: tagm})
	if err != nil {
		log.Printf("TagRepository:getIdByName() Error mutate %s", err)
		return "", fmt.Errorf("TagRepository:getIdByName() Error mutate %s", err)
	}
	tagId := mutate.GetUids()["tag"]
	return tagId, err

}

func (r TagRepository) getByName(ctx context.Context, txn *dgo.Txn, tag *domain.Tag) error {
	var vars *api.Response
	var err error
	vars, err = txn.QueryWithVars(ctx, getTagByName, map[string]string{"$name": tag.Name}) //TODO посмотреть другой вариант
	response := domain.TagList{}
	err = json.Unmarshal(vars.Json, &response)
	if err != nil {
		log.Printf("TagRepository:getByName() Error Unmarshal %s", err)
		return fmt.Errorf("TagRepository:getByName() Error Unmarshal %s", err)
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
var findAllTag = `query findAllTag($search: string)
{
	var(func: type(Tag)) @filter(regexp(name, $search)){
	    tags as uid
        countTag as count(~posts)
  }
	tagList(func: uid(tags), orderdesc: val(countTag)){
	     id: uid
         name
  }
}`
