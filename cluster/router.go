package cluster

import "Go_redis/interface/resp"

func makeRouter() map[string]CmdFunc {
	routerMap := make(map[string]CmdFunc)
	routerMap["exists"] = defaultFunc // exists key
	routerMap["type"] = defaultFunc   // type key
	routerMap["set"] = defaultFunc    // set key value
	routerMap["setnx"] = defaultFunc  // setnx key value
	routerMap["get"] = defaultFunc    // get key
	routerMap["getset"] = defaultFunc // getset key value
	routerMap["ping"] = ping
	routerMap["rename"] = rename
	routerMap["renamenx"] = rename
	routerMap["flushdb"] = flushDB
	routerMap["del"] = del
	routerMap["select"] = execSelect

	return routerMap
}

// GET key // SET k1 v1
func defaultFunc(cluster *ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply {
	key := string(cmdArgs[1])
	peer := cluster.peerPicker.PickNode(key)
	return cluster.relay(peer, conn, cmdArgs)
}
