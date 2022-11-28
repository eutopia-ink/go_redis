package database

import (
	"Go_redis/datastruct/dict"
	"Go_redis/interface/database"
	"Go_redis/interface/resp"
	"Go_redis/resp/reply"
	"strings"
)

type DB struct {
	index int
	data  dict.Dict
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply

type CmdLine = [][]byte

func makeDB() *DB {
	return &DB{
		data: dict.MakeSyncDict(),
	}
}
func (db *DB) Exec(c resp.Connection, cmdline CmdLine) resp.Reply {
	cmdName := strings.ToLower(string(cmdline[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.MakeErrReply("ERR unknown command " + cmdName)
	}
	if !validateArity(cmd.arity, cmdline) {
		return reply.MakeArgNumErrReply(cmdName)
	}
	return cmd.exector(db, cmdline[1:])
}

func validateArity(arity int, cmdArgs [][]byte) bool {
	argNum := len(cmdArgs)
	if arity >= 0 {
		return argNum == arity
	}
	return argNum >= -arity
}

func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	entity, _ := raw.(*database.DataEntity)
	return entity, true
}

func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

// Remove the given key from db
func (db *DB) Remove(key string) {
	db.data.Remove(key)
}

// Removes the given keys from db
func (db *DB) Removes(keys ...string) (deleted int) {
	deleted = 0
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.Remove(key)
			deleted++
		}
	}
	return deleted
}

// Flush clean database
func (db *DB) Flush() {
	db.data.Clear()

}
