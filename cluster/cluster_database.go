package cluster

import (
	"Go_redis/config"
	database2 "Go_redis/database"
	"Go_redis/interface/database"
	"Go_redis/interface/resp"
	"Go_redis/lib/consistenthash"
	"Go_redis/lib/logger"
	"Go_redis/resp/reply"
	"context"
	pool "github.com/jolestar/go-commons-pool/v2"
	"strings"
)

type ClusterDatabase struct {
	self       string
	nodes      []string
	peerPicker *consistenthash.NodeMap

	peerConnection map[string]*pool.ObjectPool
	db             database.Database
}

func MakeClusterDatabase() *ClusterDatabase {
	cluster := &ClusterDatabase{
		self:           config.Properties.Self,
		db:             database2.NewStandaloneDatabase(),
		peerPicker:     consistenthash.NewNodeMap(nil),
		peerConnection: make(map[string]*pool.ObjectPool),
	}
	nodes := make([]string, 0, len(config.Properties.Peers)+1)
	for _, peer := range config.Properties.Peers {
		nodes = append(nodes, peer)
	}
	nodes = append(nodes, config.Properties.Self)
	cluster.peerPicker.AddNode(nodes...)
	ctx := context.Background()
	for _, peer := range config.Properties.Peers {
		cluster.peerConnection[peer] = pool.NewObjectPoolWithDefaultConfig(ctx, &connectionFactory{
			Peer: peer,
		})
	}
	cluster.nodes = nodes
	return cluster
}

type CmdFunc func(cluster *ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply

var router = makeRouter()

func (cluster *ClusterDatabase) Exec(client resp.Connection, args [][]byte) (result resp.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			result = &reply.UnknowErrReply{}
		}
	}()

	cmdName := strings.ToLower(string(args[0]))
	cmdFunc, ok := router[cmdName]
	if !ok {
		reply.MakeErrReply("Err: Not support command")
	}
	result = cmdFunc(cluster, client, args)
	return
}

func (cluster *ClusterDatabase) AfterClientClose(conn resp.Connection) {
	cluster.db.AfterClientClose(conn)
}

func (cluster *ClusterDatabase) Close() {
	cluster.db.Close()
}
