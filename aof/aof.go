package aof

import (
	"Go_redis/config"
	"Go_redis/interface/database"
	"Go_redis/lib/logger"
	"Go_redis/lib/utils"
	"Go_redis/resp/connection"
	"Go_redis/resp/parser"
	"Go_redis/resp/reply"
	"io"
	"os"
	"strconv"
)

type CmdLine = [][]byte

const (
	AOF_BUFFER_SIZE = 1 << 16
)

type payload struct {
	cmdLine CmdLine
	dbIndex int
}

type AofHandler struct {
	database    database.Database
	aofChan     chan *payload
	aofFile     *os.File
	aofFilename string
	currentDB   int // 当前工作的db索引
}

// NewAofHandler
func NewAofHandler(database database.Database) (*AofHandler, error) {
	handler := &AofHandler{}
	handler.aofFilename = config.Properties.AppendFilename
	handler.database = database
	handler.LoadAof()
	aofile, err := os.OpenFile(handler.aofFilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handler.aofFile = aofile
	handler.aofChan = make(chan *payload, AOF_BUFFER_SIZE)
	// 处理放入channel的指令
	go func() {
		handler.handleAof()
	}()
	return handler, nil
}

// Add payload(set k v ) -> aofChan
func (handler AofHandler) AddAof(dbIndex int, cmd CmdLine) {
	if config.Properties.AppendOnly == true && handler.aofChan != nil {
		handler.aofChan <- &payload{
			cmdLine: cmd,
			dbIndex: dbIndex,
		}
	}

}

// handleAof payload(set k v ) <- aofChan (落盘)
func (handler *AofHandler) handleAof() {
	handler.currentDB = 0
	for p := range handler.aofChan {
		// 有切换数据库的操作
		if p.dbIndex != handler.currentDB {
			data := reply.MakeMultiBulkReply(utils.ToCmdLine("select", strconv.Itoa(p.dbIndex))).ToBytes()
			_, err := handler.aofFile.Write(data)
			if err != nil {
				logger.Error(err)
				continue
			}
			handler.currentDB = p.dbIndex
		}
		data := reply.MakeMultiBulkReply(p.cmdLine).ToBytes()
		_, err := handler.aofFile.Write(data)
		if err != nil {
			logger.Error(err)
		}

	}
}

// LoadAof
func (handler *AofHandler) LoadAof() {
	file, err := os.Open(handler.aofFilename)
	if err != nil {
		logger.Error(err)
		return
	}
	defer file.Close()
	fakeConn := &connection.Connection{}
	ch := parser.ParseStream(file)
	for p := range ch {
		if p.Err != nil {
			if p.Err == io.EOF {
				break
			}
			logger.Error(p.Err)
			continue
		}
		if p.Data == nil {
			logger.Error("empty payload")
			continue
		}
		r, ok := p.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("need multi bulk")
			continue
		}
		rep := handler.database.Exec(fakeConn, r.Args)
		if reply.IsErrReply(rep) {
			logger.Error(rep)
		}
	}

}
