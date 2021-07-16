package mongolib

import "go.mongodb.org/mongo-driver/bson"

type filter bson.A

func Filter() filter {
	return filter{}
}

func (f filter) Equal(key string, value interface{}) filter {
	filter := append(f, bson.D{{Key: key, Value: value}})
	return filter
}

func (f filter) Regex(key string, value interface{}) filter {
	filter := append(f, bson.D{{Key: key, Value: bson.D{{"$regex", value}}}})
	return filter
}

func (f filter) NotEqual(key string, value interface{}) filter {
	filter := append(f, bson.D{{Key: key, Value: bson.D{{"$ne", value}}}})
	return filter
}

func (f filter) GreaterThan(key string, value interface{}) filter {
	filter := append(f, bson.D{{Key: key, Value: bson.D{{"$gt", value}}}})
	return filter
}

func (f filter) GreaterThanEqual(key string, value interface{}) filter {
	filter := append(f, bson.D{{Key: key, Value: bson.D{{"$gte", value}}}})
	return filter
}

func (f filter) LessThan(key string, value interface{}) filter {
	filter := append(f, bson.D{{Key: key, Value: bson.D{{"$lt", value}}}})
	return filter
}

func (f filter) LessThanEqual(key string, value interface{}) filter {
	filter := append(f, bson.D{{Key: key, Value: bson.D{{"$lte", value}}}})
	return filter
}

func (f filter) In(key string, value interface{}) filter {
	filter := append(f, bson.D{{Key: key, Value: bson.D{{"$in", value}}}})
	return filter
}

func (f filter) NotIn(key string, value interface{}) filter {
	filter := append(f, bson.D{{Key: key, Value: bson.D{{"$nin", value}}}})
	return filter
}
