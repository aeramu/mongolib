package mongolib

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewRepository(collection *mongo.Collection) *Repository {
	return &Repository{
		Collection: collection,
	}
}

type Repository struct{
	*mongo.Collection
}

func (repo *Repository) Save(ctx context.Context, id primitive.ObjectID, data interface{}) error {
	update := bson.D{{"$set", data}}
	filter := bson.D{{"_id", id}}
	opt := options.Update().SetUpsert(true)

	if _, err := repo.UpdateOne(ctx, filter, update, opt); err != nil {
		return err
	}
	return nil
}

func (repo *Repository) Query() Query {
	return Query{
		coll:   repo.Collection,
		filter: bson.A{},
		limit:  0,
		sort:   bson.D{},
	}
}