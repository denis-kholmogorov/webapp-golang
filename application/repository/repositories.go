package repository

import (
	"context"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
)

var conn *DGraphConn

type DGraphConn struct {
	connection *dgo.Dgraph
}

func NewDGraphConn(connection *dgo.Dgraph) {
	conn = &DGraphConn{connection: connection}
}

func GetDGraphConn() *DGraphConn {
	return conn
}

func DeleteEdge(txn *dgo.Txn, ctx context.Context, objectID, edgeName, edgeUID string, commitNow bool) error {
	mu := &api.Mutation{CommitNow: commitNow,
		Del: []*api.NQuad{
			{
				Subject:   objectID,
				Predicate: edgeName,
				ObjectId:  edgeUID,
			},
		},
	}
	_, err := txn.Mutate(ctx, mu)
	if err != nil {
		return err
	}
	return nil
}

func Delete(txn *dgo.Txn, ctx context.Context, objectID string, commitNow bool) error {
	mu := &api.Mutation{CommitNow: commitNow,
		Set: []*api.NQuad{
			{
				Subject:     objectID,
				Predicate:   "isDeleted",
				ObjectValue: &api.Value{Val: &api.Value_StrVal{StrVal: "true"}},
			},
		},
	}
	_, err := txn.Mutate(ctx, mu)
	return err
}

func AddEdge(txn *dgo.Txn, ctx context.Context, objectID, edgeName, edgeUID string, commitNow bool) error {
	mu := &api.Mutation{CommitNow: commitNow,
		Set: []*api.NQuad{
			{
				Subject:   objectID,
				Predicate: edgeName,
				ObjectId:  edgeUID,
			},
		},
	}
	_, err := txn.Mutate(ctx, mu)
	return err
}

func UpdateNodeFields(txn *dgo.Txn, ctx context.Context, nodeID string, fields map[string]string, commitNow bool) error {

	var nquads []*api.NQuad
	for fieldName, fieldValue := range fields {
		nquads = append(nquads, &api.NQuad{
			Subject:     nodeID,
			Predicate:   fieldName,
			ObjectValue: &api.Value{Val: &api.Value_StrVal{StrVal: fieldValue}},
		})
	}
	_, err := txn.Mutate(ctx, &api.Mutation{Set: nquads, CommitNow: commitNow})
	return err
}

func AddNodesToEdge(txn *dgo.Txn, ctx context.Context, nodeID, edgeName string, fields []string, commitNow bool) error {
	nquads := make([]*api.NQuad, 0, len(fields))
	for _, targetNodeID := range fields {
		nquads = append(nquads, &api.NQuad{
			Subject:   nodeID,
			Predicate: edgeName,
			ObjectId:  targetNodeID,
		})
	}
	_, err := txn.Mutate(ctx, &api.Mutation{Set: nquads, CommitNow: commitNow})
	return err
}
