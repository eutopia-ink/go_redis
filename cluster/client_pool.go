package cluster

import (
	"Go_redis/resp/client"
	"context"
	"errors"
	pool "github.com/jolestar/go-commons-pool/v2"
)

// 连接工厂
type connectionFactory struct {
	Peer string
}

func (c connectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	client, err := client.MakeClient(c.Peer)
	if err != nil {
		return nil, err
	}
	client.Start()
	return pool.NewPooledObject(client), nil
}

func (c connectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	curClient, ok := object.Object.(*client.Client)
	if !ok {
		return errors.New("type mismatch")
	}
	curClient.Close()
	return nil
}

func (c connectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return true
}

func (c connectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (c connectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
