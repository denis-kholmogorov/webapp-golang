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

func (r FriendRepository) RequestFriend(currentUserId string, friendId string) error {
	friendshipTo := domain.Friendship{Uid: "_:friendTo", Status: domain.REQUEST_TO, FriendId: friendId, ReverseStatus: domain.REQUEST_FROM, DType: []string{"Friendship"}}
	friendshipFrom := domain.Friendship{Uid: "_:friendFrom", Status: domain.REQUEST_FROM, FriendId: currentUserId, ReverseStatus: domain.REQUEST_TO, DType: []string{"Friendship"}}
	ctx := context.Background()
	txn := r.conn.NewTxn()

	friendshipsm, err := json.Marshal([]domain.Friendship{friendshipTo, friendshipFrom})
	if err != nil {
		txn.Discard(ctx)
		log.Printf("FriendRepository:Request() Error marhalling friend %s", err)
		return fmt.Errorf("FriendRepository:Request() Error marhalling friend %s", err)
	}
	mutate, err := txn.Mutate(ctx, &api.Mutation{SetJson: friendshipsm})
	if err != nil {
		txn.Discard(ctx)
		log.Printf("FriendRepository:Request() Error mutate %s", err)
		return fmt.Errorf("FriendRepository:Request() Error mutate %s", err)
	}
	edges := []dto.Edge{
		{currentUserId, "friends", mutate.GetUids()["friendTo"]},  // текущий добавляет дружбу TO
		{mutate.GetUids()["friendTo"], "friend", friendId},        // дружба от текущего добавляет друга
		{friendId, "friends", mutate.GetUids()["friendFrom"]},     // другу добавляет дружбу FROM
		{mutate.GetUids()["friendFrom"], "friend", currentUserId}, // дружба друга добавляет текущего
	}
	err = AddEdges(txn, ctx, edges, true)

	return err
}

func (r FriendRepository) ApproveFriend(currentUserId string, friendId string) error {
	ctx := context.Background()
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	variables["$friendId"] = friendId

	mu := &api.Mutation{
		SetNquads: []byte(updateToFriend),
	}

	req := &api.Request{
		Query:     getFriendship,
		Mutations: []*api.Mutation{mu},
		Vars:      variables,
		CommitNow: true,
	}

	// Update email only if matching uid found.
	resp, err := r.conn.NewTxn().Do(ctx, req)
	if err != nil {
		return err
	}
	log.Printf(string(rune(len(resp.Uids))))

	return nil
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

func (r FriendRepository) Delete(currentUserId string, friendId string) error {
	ctx := context.Background()
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	variables["$friendId"] = friendId

	mu := &api.Mutation{
		DelNquads: []byte(deleteFromFriend),
	}

	req := &api.Request{
		Query:     deleteFriendship,
		Mutations: []*api.Mutation{mu},
		Vars:      variables,
		CommitNow: true,
	}

	// Update email only if matching uid found.
	resp, err := r.conn.NewTxn().Do(ctx, req)
	if err != nil {
		return err
	}
	log.Printf(string(rune(len(resp.Uids))))

	return nil
}

func (r FriendRepository) Count(currentUserId string) (int, error) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	var vars *api.Response
	var err error
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	variables["$status"] = domain.REQUEST_FROM
	vars, err = txn.QueryWithVars(ctx, getCount, variables)

	if err != nil {
		log.Printf("FriendRepository:FindAll() Error query %s", err)
		return 0, fmt.Errorf("FriendRepository:FindAll() Error query %s", err)
	}
	response := domain.CountRequest{}
	err = json.Unmarshal(vars.Json, &response)
	if err != nil {
		log.Printf("FriendRepository:FindAll() Error Unmarshal %s", err)
		return 0, fmt.Errorf("FriendRepository:FindAll() Error Unmarshal %s", err)
	}

	if len(response.CountRequest) > 0 {
		return response.CountRequest[0].Count, nil
	} else {
		return 0, nil
	}
}

func (r FriendRepository) getMyFriends(currentUserId string) ([]string, error) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	var vars *api.Response
	var err error
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	variables["$statusFriend"] = domain.FRIEND
	variables["$statusSubscribe"] = domain.SUBSCRIBED
	vars, err = txn.QueryWithVars(ctx, getMyFriends, variables)

	if err != nil {
		log.Printf("FriendRepository:FindAll() Error query %s", err)
		return nil, fmt.Errorf("FriendRepository:FindAll() Error query %s", err)
	}
	response := domain.IdsFriends{}
	err = json.Unmarshal(vars.Json, &response)
	if err != nil {
		log.Printf("FriendRepository:FindAll() Error Unmarshal %s", err)
		return nil, fmt.Errorf("FriendRepository:FindAll() Error Unmarshal %s", err)
	}
	var ids []string
	for _, friend := range response.Ids {
		ids = append(ids, friend.Id)
	}
	return ids, nil
}

func (r FriendRepository) Recommendations(currentUserId string) ([]domain.Account, error) {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	var vars *api.Response
	var err error
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	variables["$statusFriend"] = domain.FRIEND
	vars, err = txn.QueryWithVars(ctx, getRecommendations, variables)

	if err != nil {
		log.Printf("FriendRepository:FindAll() Error query %s", err)
		return nil, fmt.Errorf("FriendRepository:FindAll() Error query %s", err)
	}
	response := dto.PageResponse[domain.Account]{}
	err = json.Unmarshal(vars.Json, &response)
	if err != nil {
		log.Printf("FriendRepository:FindAll() Error Unmarshal %s", err)
		return nil, fmt.Errorf("FriendRepository:FindAll() Error Unmarshal %s", err)
	}
	return response.Content, nil
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
	photo: photo
  }
  count(func: uid(A)){
		totalElement:count(uid)
  }
}
`

var getFriendship = `query setFriend($currentUserId: string, $friendId: string)  {
	fr1(func: uid($currentUserId)){
		friends @filter(eq(friendId,$friendId)){
			A as uid
		}
	}
	fr2(func: uid($friendId)){
		friends @filter(eq(friendId,$currentUserId)){
			B as uid
		}
	}
}`

var updateToFriend = `uid(A) <status> "FRIEND" .
                      uid(B) <status> "FRIEND" .`

var deleteFriendship = `query setFriend($currentUserId: string, $friendId: string)  {
	fr1(func: uid($currentUserId)){
		A as uid
		friends @filter(eq(friendId,$friendId)){
			B as uid
		}
	}
	fr2(func: uid($friendId)){
		C as uid
		friends @filter(eq(friendId,$currentUserId)){
			D as uid
		}
	}
}`

var deleteFromFriend = `uid(A) <friends> uid(B) .
						   uid(C) <friends> uid(D) .
                           uid(B) * * .
						   uid(D) * * .`

// TODO переделать на var и посмотреть count
var getCount = `query count($currentUserId: string, $status: string)  { 
	countRequest(func: uid($currentUserId)) @normalize {
		friends @filter(eq(status, $status)){
			count :count(uid)
		}
	}
}`

var getMyFriends = `query count($currentUserId: string, $statusFriend: string, $statusSubscribe: string)  { 
	ids(func: uid($currentUserId)) @normalize {
      friends @filter(eq(status, $statusFriend) or eq(status, $statusSubscribe)){
        friend{
		  id : uid
        }		
      }
	}
}`

var getRecommendations = `query count($currentUserId: string, $statusFriend: string)  {
	var(func: uid($currentUserId)) @normalize {
      friends @filter(eq(status, $statusFriend)){
        friend{
          friends @filter(eq(status, $statusFriend)){
      		friend@filter(not uid($currentUserId)){
						A as uid
      		}
          }
		}
  	  }
    }
   content(func: uid(A)){
     id:uid
	 firstName:firstName
	 lastName:lastName
	 city:city
	 country:country
	 birthDate:birthDate
	 isOnline:isOnline
	 photo: photo
   }
}`
