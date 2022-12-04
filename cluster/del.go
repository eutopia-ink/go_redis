package cluster

import (
	"Go_redis/interface/resp"
	"Go_redis/resp/reply"
)

func del(cluster *ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcast(conn, cmdArgs)
	var errReply reply.ErrorReply
	var deleted int64
	for _, r := range replies {
		if reply.IsErrReply(r) {
			errReply = r.(reply.ErrorReply)
			break
		}
		intReply, ok := r.(*reply.IntReply)
		if !ok {
			errReply = reply.MakeErrReply("Err")
		}
		deleted += intReply.Code
	}
	if errReply != nil {
		return reply.MakeErrReply("error: " + errReply.Error())
	} else {
		return reply.MakeIntReply(deleted)
	}
}
