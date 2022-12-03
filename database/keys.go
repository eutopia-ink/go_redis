package database

import (
	"Go_redis/interface/resp"
	"Go_redis/lib/utils"
	"Go_redis/lib/wildcard"
	"Go_redis/resp/reply"
)

// DEL
func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	deleted := db.Removes(keys...)
	if deleted > 0 {
		db.addAof(utils.ToCmdLine2("DEL", args...))
	}
	return reply.MakeIntReply(int64(deleted))
}

// EXISTS
func execExists(db *DB, args [][]byte) resp.Reply {
	var result int64
	for _, v := range args {
		if _, exist := db.GetEntity(string(v)); exist {
			result++
		}
	}
	return reply.MakeIntReply(result)
}

// KEYS
func execKeys(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.data.ForEach(func(key string, val interface{}) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})
	return reply.MakeMultiBulkReply(result)
}

// FLUSHDB
func execFlushDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	db.addAof(utils.ToCmdLine2("FLUSHDB", args...))
	return reply.MakeOkReply()
}

// TYPE
func execType(db *DB, args [][]byte) resp.Reply {
	var result string
	key := string(args[0])
	entity, exist := db.GetEntity(key)
	if !exist {
		return reply.MakeStatusReply("none")
	}
	switch entity.Data.(type) {
	case []byte:
		result = "string"
	}
	if result == "" {
		return &reply.UnknowErrReply{}
	}
	return reply.MakeStatusReply(result)
}

// RENAME
func execRename(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])
	entity, exist := db.GetEntity(src)
	if !exist {
		return reply.MakeErrReply("no such key")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)
	db.addAof(utils.ToCmdLine2("RENAME", args...))
	return reply.MakeOkReply()
}

// RENAMENX
func execRenamenx(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])
	_, exist := db.GetEntity(dest)
	if exist {
		return reply.MakeIntReply(0)
	}
	entity, exist := db.GetEntity(src)
	if !exist {
		return reply.MakeErrReply("no such key")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)
	db.addAof(utils.ToCmdLine2("RENAMENX", args...))
	return reply.MakeOkReply()
}

func init() {
	RegisterCommand("DEL", execDel, -2)
	RegisterCommand("EXISTS", execExists, -2)
	RegisterCommand("KEYS", execKeys, 2)
	RegisterCommand("FLUSHDB", execFlushDB, -1)
	RegisterCommand("TYPE", execType, 2)
	RegisterCommand("RENAME", execRename, 3)
	RegisterCommand("RENAMENX", execRenamenx, 3)
}
