package mongolib

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Query struct {
	coll *mongo.Collection
	filter bson.A
	limit int
	sort bson.D
}

func (q Query) Asc(key string) Query {
	q.sort = append(q.sort, bson.E{Key: key, Value: 1})
	return q
}

func (q Query) Desc(key string) Query {
	q.sort = append(q.sort, bson.E{Key: key, Value: -1})
	return q
}

func (q Query) Equal(key string, value interface{}) Query {
	q.filter = append(q.filter, bson.D{{Key: key, Value: value}})
	return q
}

func (q Query) GreaterThan(key string, value interface{}) Query {
	q.filter = append(q.filter, bson.D{{Key: key, Value: bson.D{{"$gt", value}}}})
	return q
}

func (q Query) LessThan(key string, value interface{}) Query {
	q.filter = append(q.filter, bson.D{{Key: key, Value: bson.D{{"$lt", value}}}})
	return q
}

func (q Query) Limit(limit int) Query {
	q.limit = limit
	return q
}

func (q Query) Find(ctx context.Context) Result {
	filter := bson.D{}
	if len(q.filter) > 0 {
		filter = bson.D{{Key: "$and", Value: q.filter}}
	}
	opt := options.Find().SetSort(q.sort)
	if q.limit > 0 {
		opt = opt.SetLimit(int64(q.limit))
	}

	cur, err := q.coll.Find(ctx, filter, opt)
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

func (q Query) FindOne(ctx context.Context) Result {
	filter := bson.D{{Key: "&and", Value: q.filter}}

	result := q.coll.FindOne(ctx, filter)
	return &SingleResult{
		SingleResult: result,
	}
}