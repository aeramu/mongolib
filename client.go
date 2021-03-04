package mongolib

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

var (
	dbLock sync.Once
	client *mongo.Client
)

func NewSingletonClient(ctx context.Context, uri string, opts ...*options.ClientOptions) (*mongo.Client, error) {
	opt := options.MergeClientOptions(opts...).ApplyURI(uri)
	var err error
	dbLock.Do(func() {
		client, err = mongo.Connect(ctx, opt)
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}
