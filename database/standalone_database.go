package database

import (
	"Go_redis/aof"
	"Go_redis/config"
	"Go_redis/interface/resp"
	"Go_redis/lib/logger"
	"Go_redis/resp/reply"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
)

type StandaloneDatabase struct {
	dbSet      []*DB           // 连接池
	aofHandler *aof.AofHandler // aof记录
}

func NewStandaloneDatabase() *StandaloneDatabase {
	mdb := &StandaloneDatabase{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	mdb.dbSet = make([]*DB, config.Properties.Databases)
	for i := range mdb.dbSet {
		singleDB := makeDB()
		singleDB.index = i
		mdb.dbSet[i] = singleDB
	}
	if config.Properties.AppendOnly == true {
		aofHandler, err := aof.NewAofHandler(mdb)
		if err != nil {
			panic(err)
		}
		mdb.aofHandler = aofHandler
		for _, db := range mdb.dbSet {
			ddb := db
			ddb.addAof = func(line CmdLine) {
				mdb.aofHandler.AddAof(ddb.index, line)
			}
		}
	}

	return mdb
}

func (mdb *StandaloneDatabase) Exec(c resp.Connection, cmdLine [][]byte) (result resp.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
		}
	}()
	cmdName := strings.ToLower(string(cmdLine[0]))
	if cmdName == "select" {
		if len(cmdLine) != 2 {
			return reply.MakeArgNumErrReply("select")
		}
		return execSelect(c, mdb, cmdLine[1:])
	}

	dbIndex := c.GetDBIndex()
	selectedDB := mdb.dbSet[dbIndex]
	return selectedDB.Exec(c, cmdLine)
}

func (mdb *StandaloneDatabase) Close() {

}

func (mdb *StandaloneDatabase) AfterClientClose(c resp.Connection) {
}

func execSelect(c resp.Connection, mdb *StandaloneDatabase, args [][]byte) resp.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeErrReply("ERR invalid DB index")
	}
	if dbIndex >= len(mdb.dbSet) {
		return reply.MakeErrReply("ERR DB index is out of range")
	}
	c.SelectDB(dbIndex)
	return reply.MakeOkReply()
}
