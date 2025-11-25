package client

import (
	"zero-admin/rpc/sys/db/mysql/query"
)

type MysqlDB struct {
	q *query.Query
}

func NewMysqlDB(q *query.Query) (*MysqlDB, error) {
	return &MysqlDB{q}, nil
}
