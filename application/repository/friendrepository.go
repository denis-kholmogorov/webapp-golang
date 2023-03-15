package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"log"
	"strconv"
	"web/application/domain"
	"web/application/dto"
	"web/application/utils"
)

var friendRepo *FriendRepository
var isInitializedFriendRepo bool

type FriendRepository struct {
	conn *dgo.Dgraph
}

func GetFriendRepository() *FriendRepository {
	if !isInitializedFriendRepo {
		friendRepo = &FriendRepository{}
		friendRepo.conn = GetDGraphConn().connection
		isInitializedFriendRepo = true
	}
	return friendRepo
}

func (r FriendRepository) RequestFriend(currentId string, friendId string) error {
	friendshipTo := domain.Friendship{Uid: "_:friendTo", Status: domain.REQUEST_TO, ReverseStatus: domain.REQUEST_FROM, DType: []string{"Friendship"}}
	friendshipFrom := domain.Friendship{Uid: "_:friendFrom", Status: domain.REQUEST_FROM, ReverseStatus: domain.REQUEST_TO, DType: []string{"Friendship"}}
	ctx := context.Background()
	txn := r.conn.NewTxn()

	friendshipsm, err := json.Marshal([]domain.Friendship{friendshipTo, friendshipFrom})
	if err != nil {
		txn.Discard(ctx)
		log.Printf("FriendRepository:RequestFriend() Error marhalling friend %s", err)
		return fmt.Errorf("FriendRepository:RequestFriend() Error marhalling friend %s", err)
	}
	mutate, err := txn.Mutate(ctx, &api.Mutation{SetJson: friendshipsm})
	if err != nil {
		txn.Discard(ctx)
		log.Printf("FriendRepository:RequestFriend() Error mutate %s", err)
		return fmt.Errorf("FriendRepository:RequestFriend() Error mutate %s", err)
	}
	edges := []dto.Edge{
		{currentId, "friends", mutate.GetUids()["friendTo"]},  // текущий добавляет дружбу TO
		{mutate.GetUids()["friendTo"], "friend", friendId},    // дружба от текущего добавляет друга
		{friendId, "friends", mutate.GetUids()["friendFrom"]}, // другу добавляет дружбу FROM
		{mutate.GetUids()["friendFrom"], "friend", currentId}, // дружба друга добавляет текущего
	}
	err = AddEdges(txn, ctx, edges, true)

	return err
}

func (r FriendRepository) FindAll(currentUserId string, statusCode dto.StatusCode, page dto.PageRequest) (interface{}, interface{}) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	var vars *api.Response
	var err error
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	variables["$first"] = strconv.Itoa(page.Size)
	variables["$offset"] = strconv.Itoa(page.Size * utils.GetPageNumber(&page))
	variables["$statusCode"] = statusCode.StatusCode

	vars, err = txn.QueryWithVars(ctx, getAllFriends, variables)

	if err != nil {
		log.Printf("FriendRepository:FindAll() Error query %s", err)
		return nil, fmt.Errorf("FriendRepository:FindAll() Error query %s", err)
	}
	response := dto.PageResponse[domain.Account]{}
	err = json.Unmarshal(vars.Json, &response)
	for i := range response.Content {
		response.Content[i].StatusCode = statusCode.StatusCode
	}
	if err != nil {
		log.Printf("FriendRepository:FindAll() Error Unmarshal %s", err)
		return nil, fmt.Errorf("FriendRepository:FindAll() Error Unmarshal %s", err)
	}

	response.SetPage(page.Size, page.Page)
	return &response, nil
}

var getAllFriends = `query Posts($currentUserId: string, $statusCode: string, $first: int, $offset: int)
{
  var(func: uid($currentUserId)) @filter(eq(isDeleted, false))  {
    friends @filter(eq(status,$statusCode)) {
      friend {
        A as uid
      }
    }
  }
  content(func: uid(A), orderdesc: time, first: $first, offset: $offset)  {
    id:uid
	firstName:firstName
	lastName:lastName
	city:city
	country:country
	birthDate:birthDate
	isOnline:isOnline
	imagePath: photo
  }
  count(func: uid(A)){
		totalElement:count(uid)
  }
}
`
