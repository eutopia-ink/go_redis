package cluster

import (
	"Go_redis/interface/resp"
	"Go_redis/lib/utils"
	"Go_redis/resp/client"
	"Go_redis/resp/reply"
	"context"
	"errors"
	"strconv"
)

func (cluster *ClusterDatabase) getPeerClient(peer string) (*client.Client, error) {
	pool, ok := cluster.peerConnection[peer]
	if !ok {
		return nil, errors.New("Connection not found")
	}
	object, err := pool.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}
	curClient, ok := object.(*client.Client)
	if !ok {
		return nil, errors.New("Wrong type in cliuster database connection")
	}
	return curClient, nil
}

func (cluster *ClusterDatabase) returnPeerClient(peer string, curClient *client.Client) error {
	pool, ok := cluster.peerConnection[peer]
	if !ok {
		return errors.New("Connection not found")
	}
	return pool.ReturnObject(context.Background(), curClient)
}

func (cluster *ClusterDatabase) relay(peer string, conn resp.Connection, args [][]byte) resp.Reply {
	if peer == cluster.self {
		return cluster.db.Exec(conn, args)
	}
	peerClient, err := cluster.getPeerClient(peer)
	if err != nil {
		return reply.MakeErrReply(err.Error())
	}
	defer cluster.returnPeerClient(peer, peerClient)
	peerClient.Send(utils.ToCmdLine("SELECT", strconv.Itoa(conn.GetDBIndex())))
	return peerClient.Send(args)
}

func (cluster *ClusterDatabase) broadcast(conn resp.Connection, args [][]byte) map[string]resp.Reply {
	results := make(map[string]resp.Reply)
	for _, node := range cluster.nodes {
		result := cluster.relay(node, conn, args)
		results[node] = result
	}
	return results
}
