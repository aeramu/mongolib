package mongolib

import (
	"context"
	"errors"
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

var (
	ErrNotFound = errors.New("result not found")
)

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

func (repo *Repository) FindOneByIndex(ctx context.Context, indexName string, indexValue interface{}, v interface{}) error {
	filter := bson.D{{indexName, indexValue}}
	err := repo.FindOne(ctx, filter).Decode(v)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (repo *Repository) FindByIndex(ctx context.Context, indexName string, indexValue interface{}, v interface{}) error {
	filter := bson.D{{indexName, indexValue}}
	cur, err := repo.Find(ctx, filter)
	if err != nil {
		return err
	}
	if err := cur.All(ctx, v); err != nil {
		return err
	}
	return nil
}
