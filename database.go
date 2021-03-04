package mongolib

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func NewDatabase(client *mongo.Client, dbname string) *Database {
	return &Database{
		Database: client.Database(dbname),
	}
}

type Database struct {
	*mongo.Database
}

func (d *Database) Coll(collectionName string) *Collection {
	return &Collection{
		Collection: d.Collection(collectionName),
	}
}
