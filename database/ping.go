package database

import (
	"Go_redis/interface/resp"
	"Go_redis/resp/reply"
)

func Ping(db *DB, args [][]byte) resp.Reply {
	return reply.MakePongReply()
}

func init() {
	RegisterCommand("ping", Ping, 1)
}
