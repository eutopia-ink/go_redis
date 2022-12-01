package database

import (
	"Go_redis/interface/database"
	"Go_redis/interface/resp"
	"Go_redis/resp/reply"
)

// GET
func (db *DB) getAsString(key string) ([]byte, reply.ErrorReply) {
	entity, ok := db.GetEntity(key)
	if !ok {
		return nil, nil
	}
	str, ok := entity.Data.(string)
	bytes := []byte(str)
	if !ok {
		return nil, &reply.WrongTypeReply{}
	}
	return bytes, nil
}

func execGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	bytes, err := db.getAsString(key)
	if err != nil {
		return err
	}
	if bytes == nil {
		return &reply.NullBulkReply{}
	}
	return reply.MakeBulkReply(bytes)
}

// SET
func execSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := string(args[1])
	entity := &database.DataEntity{
		Data: value,
	}
	db.PutEntity(key, entity)
	return reply.MakeOkReply()
}

// SETNX
func execSetnx(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := string(args[1])
	entity := &database.DataEntity{
		Data: value,
	}

	result := db.PutIfAbsent(key, entity)
	return reply.MakeIntReply(int64(result))
}

// GETSET
func execGetset(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := string(args[1])
	old, exist := db.GetEntity(key)
	entity := &database.DataEntity{
		Data: value,
	}
	db.PutEntity(key, entity)
	if !exist {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeBulkReply(old.Data.([]byte))
}

// STRLEN
func execStrlen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exist := db.GetEntity(key)
	if !exist {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeIntReply(int64(len(entity.Data.([]byte))))
}

func init() {
	RegisterCommand("GET", execGet, 2)
	RegisterCommand("SET", execSet, -3)
	RegisterCommand("SETNX", execSetnx, 3)
	RegisterCommand(" GETSET", execGetset, 3)
	RegisterCommand(" STRLEN", execStrlen, 2)
}
