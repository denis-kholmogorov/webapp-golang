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

func RemoveEdge(txn *dgo.Txn, ctx context.Context, objectID, edgeName, edgeUID string, commitNow bool) error {
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
	if err != nil {
		return err
	}
	return nil
}
