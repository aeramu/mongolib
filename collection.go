package mongolib

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
	*mongo.Collection
}

func (coll *Collection) FindByID(ctx context.Context, id primitive.ObjectID) Result {
	filter := bson.D{{"_id", id}}
	res := coll.FindOne(ctx, filter)

	return &SingleResult{
		SingleResult: res,
	}
}

func (coll *Collection) Save(ctx context.Context, id primitive.ObjectID, data interface{}) error {
	update := bson.D{{"$set", data}}
	filter := bson.D{{"_id", id}}
	opt := options.Update().SetUpsert(true)

	if _, err := coll.UpdateOne(ctx, filter, update, opt); err != nil {
		return err
	}
	return nil
}

func (coll *Collection) Query() Query {
	return Query{
		coll:   coll,
		filter: bson.A{},
		limit:  0,
		offset: 0,
		sort:   bson.D{},
	}
}