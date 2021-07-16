package mongolib

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	Ascending = 1
	Descending = -1
)

type Query struct {
	coll *Collection
	filter bson.A
	limit int
	offset int
	sort bson.D
}

// Filter

func (q Query) Equal(key string, value interface{}) Query {
	q.filter = append(q.filter, bson.D{{Key: key, Value: value}})
	return q
}

func (q Query) Regex(key string, value interface{}) Query {
	q.filter = append(q.filter, bson.D{{Key: key, Value: bson.D{{"$regex", value}}}})
	return q
}

func (q Query) NotEqual(key string, value interface{}) Query {
	q.filter = append(q.filter, bson.D{{Key: key, Value: bson.D{{"$ne", value}}}})
	return q
}

func (q Query) GreaterThan(key string, value interface{}) Query {
	q.filter = append(q.filter, bson.D{{Key: key, Value: bson.D{{"$gt", value}}}})
	return q
}

func (q Query) GreaterThanEqual(key string, value interface{}) Query {
	q.filter = append(q.filter, bson.D{{Key: key, Value: bson.D{{"$gte", value}}}})
	return q
}

func (q Query) LessThan(key string, value interface{}) Query {
	q.filter = append(q.filter, bson.D{{Key: key, Value: bson.D{{"$lt", value}}}})
	return q
}

func (q Query) LessThanEqual(key string, value interface{}) Query {
	q.filter = append(q.filter, bson.D{{Key: key, Value: bson.D{{"$lte", value}}}})
	return q
}

func (q Query) In(key string, value interface{}) Query {
	q.filter = append(q.filter, bson.D{{Key: key, Value: bson.D{{"$in", value}}}})
	return q
}

func (q Query) NotIn(key string, value interface{}) Query {
	q.filter = append(q.filter, bson.D{{Key: key, Value: bson.D{{"$nin", value}}}})
	return q
}

func (q Query) Sort(key string, order int) Query {
	q.sort = append(q.sort, bson.E{Key: key, Value: order})
	return q
}

func (q Query) Limit(limit int) Query {
	q.limit = limit
	return q
}

func (q Query) Offset(offset int) Query {
	q.offset = offset
	return q
}

// Get Result

func (q Query) Find(ctx context.Context) Result {
	filter := bson.D{}
	if len(q.filter) > 0 {
		filter = bson.D{{Key: "$and", Value: q.filter}}
	}

	opt := options.Find().SetSort(q.sort)
	if q.limit > 0 {
		opt = opt.SetLimit(int64(q.limit))
	}
	if q.offset > 0 {
		opt = opt.SetSkip(int64(q.offset))
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
	filter := bson.D{}
	if len(q.filter) > 0 {
		filter = bson.D{{Key: "$and", Value: q.filter}}
	}

	opt := options.FindOne().SetSort(q.sort)
	if q.offset > 0 {
		opt = opt.SetSkip(int64(q.offset))
	}

	result := q.coll.FindOne(ctx, filter, opt)
	return &SingleResult{
		SingleResult: result,
	}
}

func (q Query) Count(ctx context.Context) (int, error) {
	filter := bson.D{}
	if len(q.filter) > 0 {
		filter = bson.D{{Key: "$and", Value: q.filter}}
	}

	opt := options.Count()
	if q.offset > 0 {
		opt = opt.SetSkip(int64(q.offset))
	}
	if q.limit > 0 {
		opt = opt.SetLimit(int64(q.limit))
	}

	count, err := q.coll.CountDocuments(ctx, filter, opt)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (q Query) DeleteOne(ctx context.Context) error {
	filter := bson.D{}
	if len(q.filter) > 0 {
		filter = bson.D{{Key: "$and", Value: q.filter}}
	}

	_, err := q.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (q Query) DeleteMany(ctx context.Context) error {
	filter := bson.D{}
	if len(q.filter) > 0 {
		filter = bson.D{{Key: "$and", Value: q.filter}}
	}

	_, err := q.coll.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (q Query) Save(ctx context.Context, data interface{}) error {
	filter := bson.D{}
	if len(q.filter) > 0 {
		filter = bson.D{{Key: "$and", Value: q.filter}}
	}

	opt := options.Update().SetUpsert(true)

	update := bson.D{{"$set", data}}
	if _, err := q.coll.UpdateOne(ctx, filter, update, opt); err != nil {
		return err
	}

	return nil
}