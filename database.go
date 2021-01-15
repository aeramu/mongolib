package mongolib

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbLock sync.Once
	client *mongo.Client
)

func NewDatabase(uri string, dbname string) *mongo.Database {
	var err error
	dbLock.Do(func() {
		client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(
			uri,
		))
	})

	if err != nil {
		log.WithError(err).Errorln("failed connect MongoDB")
		return nil
	}

	return client.Database(dbname)
}
