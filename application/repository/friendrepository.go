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
	"web/application/errorhandler"
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

func (r FriendRepository) RequestFriend(currentUserId string, friendId string) {
	friendshipTo := domain.Friendship{Uid: "_:friendTo", Status: domain.REQUEST_TO, DType: []string{"Friendship"}}
	friendshipFrom := domain.Friendship{Uid: "_:friendFrom", Status: domain.REQUEST_FROM, DType: []string{"Friendship"}}
	ctx := context.Background()
	txn := r.conn.NewTxn()

	friendshipsm, err := json.Marshal([]domain.Friendship{friendshipTo, friendshipFrom})
	if err != nil {
		txn.Discard(ctx)
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("FriendRepository:RequestFriend() Marshal error %s", err)})
	}
	mutate, err := txn.Mutate(ctx, &api.Mutation{SetJson: friendshipsm})
	if err != nil {
		txn.Discard(ctx)
		panic(errorhandler.DbError{Message: fmt.Sprintf("FriendRepository:Request() Error mutate %s", err)})
	}
	edges := []dto.Edge{
		{currentUserId, "friends", mutate.GetUids()["friendTo"]},  // текущий добавляет дружбу TO
		{mutate.GetUids()["friendTo"], "friend", friendId},        // дружба от текущего добавляет друга
		{friendId, "friends", mutate.GetUids()["friendFrom"]},     // другу добавляет дружбу FROM
		{mutate.GetUids()["friendFrom"], "friend", currentUserId}, // дружба друга добавляет текущего
	}
	err = AddEdges(txn, ctx, edges, true)
	if err != nil {
		txn.Discard(ctx)
		panic(errorhandler.DbError{Message: fmt.Sprintf("FriendRepository:Request() AddEdges Error mutate %s", err)})
	}
}

func (r FriendRepository) ApproveFriend(currentUserId string, friendId string) {
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

	_, err := r.conn.NewTxn().Do(ctx, req)
	if err != nil {
		panic(errorhandler.DbError{Message: fmt.Sprintf("FriendRepository:ApproveFriend() ApproveFriend Error mutate %s", err)})
	}
}

func (r FriendRepository) FindAll(currentUserId string, statusCode dto.StatusCode, page dto.PageRequest) *dto.PageResponse[dto.FriendDto] {
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
		panic(errorhandler.DbError{Message: fmt.Sprintf("FriendRepository:FindAll() Error query getAllFriends %s", err)})
	}

	response := dto.PageResponse[dto.FriendDto]{}

	err = json.Unmarshal(vars.Json, &response)
	if err != nil {
		panic(errorhandler.DbError{Message: fmt.Sprintf("FriendRepository:FindAll() Error Unmarshal getAllFriends %s", err)})
	}

	for i := range response.Content {
		response.Content[i].StatusCode = statusCode.StatusCode
	}

	response.SetPage(page.Size, page.Page)
	return &response
}

func (r FriendRepository) Delete(currentUserId string, friendId string) {
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

	// UpdateSettings email only if matching uid found.
	_, err := r.conn.NewTxn().Do(ctx, req)
	if err != nil {
		panic(errorhandler.DbError{Message: fmt.Sprintf("FriendRepository:Delete() Error query deleteFriendship %s", err)})
	}
}

func (r FriendRepository) Count(currentUserId string) int {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	var vars *api.Response
	var err error
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	variables["$status"] = domain.REQUEST_FROM

	vars, err = txn.QueryWithVars(ctx, getCount, variables)
	if err != nil {
		panic(errorhandler.DbError{Message: fmt.Sprintf("FriendRepository:Count() Error query getCount %s", err)})
	}

	response := domain.CountRequest{}

	err = json.Unmarshal(vars.Json, &response)
	if err != nil {
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("FriendRepository:Count() Error Unmarshal getCount %s", err)})
	}

	if len(response.CountRequest) > 0 {
		return response.CountRequest[0].Count
	} else {
		return 0
	}
}

func (r FriendRepository) GetMyFriends(currentUserId string) ([]string, error) {
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

func (r FriendRepository) Recommendations(currentUserId string) []dto.FriendDto {
	ctx := context.Background()
	txn := r.conn.NewReadOnlyTxn()
	var vars *api.Response
	var err error
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	variables["$statusFriend"] = domain.FRIEND

	vars, err = txn.QueryWithVars(ctx, getRecommendations, variables)
	if err != nil {
		panic(errorhandler.DbError{Message: fmt.Sprintf("FriendRepository:Recommendations() Error query getRecommendations %s", err)})
	}

	response := dto.PageResponse[dto.FriendDto]{}

	err = json.Unmarshal(vars.Json, &response)

	if err != nil {
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("FriendRepository:Recommendations() Error marshal getRecommendations %s", err)})
	}
	return response.Content
}

func (r FriendRepository) UnBlock(currentUserId string, friendId string) {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	variables["$friendId"] = friendId
	vars, err := txn.QueryWithVars(ctx, getFriendshipStatus, variables)

	if err != nil {
		txn.Discard(ctx)
		panic(errorhandler.DbError{Message: fmt.Sprintf("FriendRepository:Block() Error query getFriendshipStatus %s", err)})
	}

	response := domain.Friendships{}

	err = json.Unmarshal(vars.Json, &response)
	if err != nil {
		txn.Discard(ctx)
		panic(errorhandler.DbError{Message: fmt.Sprintf("FriendRepository:Block() Error getFriendshipStatus %s", err)})
	}

	var mu *api.Mutation
	if response.Friendships[0].Status == domain.BLOCKED {
		mu = &api.Mutation{
			SetNquads: []byte(fmt.Sprintf(blockFriendship, domain.BLOCKED, domain.NONE, response.Friendships[0].PreviousStatus, response.Friendships[1].PreviousStatus)),
		}
	} else {
		if response.Friendships[0].PreviousStatus == "" {
			txn.Commit(ctx)
			r.Delete(currentUserId, friendId)
		}
		mu = &api.Mutation{
			SetNquads: []byte(fmt.Sprintf(blockFriendship, response.Friendships[0].PreviousStatus, response.Friendships[1].PreviousStatus, "", "")),
		}
	}

	req := &api.Request{
		Query:     getFriendship,
		Mutations: []*api.Mutation{mu},
		Vars:      variables,
		CommitNow: true,
	}

	// UpdateSettings email only if matching uid found.
	_, err = r.conn.NewTxn().Do(ctx, req)
	if err != nil {
		panic(errorhandler.DbError{Message: fmt.Sprintf("FriendRepository:Block() Error mutate getFriendship %s", blockFriendship)})
	}

}

func (r FriendRepository) Block(currentUserId string, friendId string) {
	ctx := context.Background()
	txn := r.conn.NewTxn()
	variables := make(map[string]string)
	variables["$currentUserId"] = currentUserId
	variables["$friendId"] = friendId

	vars, err := txn.QueryWithVars(ctx, getFriendshipStatus, variables)
	if err != nil {
		txn.Discard(ctx)
		panic(errorhandler.DbError{Message: fmt.Sprintf("FriendRepository:Block() Error query getFriendshipStatus %s", err)})
	}

	response := domain.Friendships{}

	err = json.Unmarshal(vars.Json, &response)
	if err != nil {
		txn.Discard(ctx)
		panic(errorhandler.DbError{Message: fmt.Sprintf("FriendRepository:Block() Error getFriendshipStatus %s", err)})
	}
	if len(response.Friendships) == 0 {
		r.createFriendship(ctx, txn, currentUserId, friendId, domain.BLOCKED, domain.NONE)
	} else {
		var mu *api.Mutation
		if response.Friendships[0].Status != domain.BLOCKED {
			mu = &api.Mutation{
				SetNquads: []byte(fmt.Sprintf(blockFriendship, domain.BLOCKED, domain.NONE, response.Friendships[0].Status, response.Friendships[1].Status)),
			}
		} else {
			if response.Friendships[0].PreviousStatus == "" {
				txn.Commit(ctx)
				r.Delete(currentUserId, friendId)
			}
			mu = &api.Mutation{
				SetNquads: []byte(fmt.Sprintf(blockFriendship, response.Friendships[0].PreviousStatus, response.Friendships[1].PreviousStatus, "", "")),
			}
		}

		req := &api.Request{
			Query:     getFriendship,
			Mutations: []*api.Mutation{mu},
			Vars:      variables,
			CommitNow: true,
		}

		// UpdateSettings email only if matching uid found.
		_, err = r.conn.NewTxn().Do(ctx, req)
		if err != nil {
			panic(errorhandler.DbError{Message: fmt.Sprintf("FriendRepository:Block() Error mutate getFriendship %s", blockFriendship)})
		}

	}
}

func (r FriendRepository) createFriendship(ctx context.Context, txn *dgo.Txn, currentUserId, friendId, statusToFriend, statusFromFriend string) {
	friendshipTo := domain.Friendship{Uid: "_:friendTo", Status: statusToFriend, DType: []string{"Friendship"}}
	friendshipFrom := domain.Friendship{Uid: "_:friendFrom", Status: statusFromFriend, DType: []string{"Friendship"}}
	friendshipsm, err := json.Marshal([]domain.Friendship{friendshipTo, friendshipFrom})
	if err != nil {
		txn.Discard(ctx)
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("FriendRepository:createFriendship() Error marshal createFriendship %s", err)})
	}

	mutate, err := txn.Mutate(ctx, &api.Mutation{SetJson: friendshipsm})
	if err != nil {
		txn.Discard(ctx)
		panic(errorhandler.DbError{Message: fmt.Sprintf("FriendRepository:createFriendship() Error mutate createFriendship %s", err)})
	}
	edges := []dto.Edge{
		{currentUserId, "friends", mutate.GetUids()["friendTo"]},  // текущий добавляет дружбу TO
		{mutate.GetUids()["friendTo"], "friend", friendId},        // дружба от текущего добавляет друга
		{friendId, "friends", mutate.GetUids()["friendFrom"]},     // другу добавляет дружбу FROM
		{mutate.GetUids()["friendFrom"], "friend", currentUserId}, // дружба друга добавляет текущего
	}
	err = AddEdges(txn, ctx, edges, true)
	if err != nil {
		txn.Discard(ctx)
		panic(errorhandler.DbError{Message: fmt.Sprintf("FriendRepository:createFriendship() Error AddEdges createFriendship %s", err)})
	}
}

var getAllFriends2 = `query Posts($currentUserId: string, $statusCode: string, $first: int, $offset: int)
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
    friendId:uid
  }
  count(func: uid(A)){
		totalElement:count(uid)
  }
}
`

var getFriendship = `query setFriend($currentUserId: string, $friendId: string)  {
	fr1(func: uid($currentUserId)){
		friends @filter(uid_in(friend, $friendId)){
			A as uid
		}
	}
	fr2(func: uid($friendId)){
		friends @filter(uid_in(friend,$currentUserId)){
			B as uid
		}
	}
}`

var getFriendshipStatus = `query setFriend($currentUserId: string, $friendId: string)  {
	fr1(func: uid($currentUserId)){
		friends @filter(uid_in(friend,$friendId)){
			A as uid
		}
	}
	fr2(func: uid($friendId)){
		friends @filter(uid_in(friend,$currentUserId)){
			B as uid
		}
	}
	 friendships(func: uid(A,B)){
      uid
      status
      previousStatus
    }
}`

var updateToFriend = `uid(A) <status> "FRIEND" .
                      uid(B) <status> "FRIEND" .`

var blockFriendship = `uid(A) <status> "%s" .
                       uid(B) <status> "%s" .
                       uid(A) <previousStatus> "%s" .
                       uid(B) <previousStatus> "%s" .`

var deleteFriendship = `query setFriend($currentUserId: string, $friendId: string)  {
	fr1(func: uid($currentUserId)){
		A as uid
		friends @filter(uid_in(friend,$friendId)){
			B as uid
		}
	}
	fr2(func: uid($friendId)){
		C as uid
		friends @filter(uid_in(friend,$currentUserId)){
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
      friends{
        friend{
		  A as uid
		}
  	  }
    }
	content(func: uid(A)) @normalize {
      friends @filter(eq(status, $statusFriend)){
        statusCode:status
        friend @filter(not uid($currentUserId) and not uid(A)){
          friendId:uid
          friends @filter(uid_in(friend,uid(A))){
      		rating:count(friend)
          }
		}
  	  }
    }
}`
