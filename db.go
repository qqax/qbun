package qbun

import "github.com/uptrace/bun"

var (
	db *bun.DB
)

func Init(dataBase *bun.DB) {
	db = dataBase
}
