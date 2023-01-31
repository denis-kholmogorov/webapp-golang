package repository

import (
	"github.com/dgraph-io/dgo/v210"
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
