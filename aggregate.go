package mongolib

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Aggregate struct {
	coll *Collection
	pipeline mongo.Pipeline
}

func (a Aggregate) Match(filter filter) Aggregate {
	a.pipeline = append(a.pipeline, bson.D{{"$match", bson.D{{"$and", filter}}}})
	return a
}

func (a Aggregate) Sort(key string, order int) Aggregate {
	a.pipeline = append(a.pipeline, bson.D{{"$sort", bson.D{{key, order}}}})
	return a
}

func (a Aggregate) Limit(limit int) Aggregate {
	a.pipeline = append(a.pipeline, bson.D{{"$limit", limit}})
	return a
}

func (a Aggregate) Offset(offset int) Aggregate {
	a.pipeline = append(a.pipeline, bson.D{{"$skip", offset}})
	return a
}

func (a Aggregate) Lookup(from, localField, foreignField, as string) Aggregate {
	a.pipeline = append(a.pipeline, bson.D{{"$lookup", bson.D{
		{"from", from},
		{"localField", localField},
		{"foreignField", foreignField},
		{"as", as},
	}}})
	return a
}

func (a Aggregate) Unwind(field string) Aggregate {
	a.pipeline = append(a.pipeline, bson.D{{"$unwind", field}})
	return a
}

func (a Aggregate) Exec(ctx context.Context) Result {
	cur, err := a.coll.Collection.Aggregate(ctx, a.pipeline)
	if err != nil {
		return &MultipleResult{
			Cursor: nil,
			Error: err,
		}
	}

	return &MultipleResult{
		Cursor: cur,
		Error:  nil,
	}
}
