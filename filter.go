package mongolib

import "go.mongodb.org/mongo-driver/bson"

type Filter bson.A

func (f Filter) Equal(key string, value interface{}) Filter {
	filter := append(f, bson.D{{Key: key, Value: value}})
	return filter
}

func (f Filter) Regex(key string, value interface{}) Filter {
	filter := append(f, bson.D{{Key: key, Value: bson.D{{"$regex", value}}}})
	return filter
}

func (f Filter) NotEqual(key string, value interface{}) Filter {
	filter := append(f, bson.D{{Key: key, Value: bson.D{{"$ne", value}}}})
	return filter
}

func (f Filter) GreaterThan(key string, value interface{}) Filter {
	filter := append(f, bson.D{{Key: key, Value: bson.D{{"$gt", value}}}})
	return filter
}

func (f Filter) GreaterThanEqual(key string, value interface{}) Filter {
	filter := append(f, bson.D{{Key: key, Value: bson.D{{"$gte", value}}}})
	return filter
}

func (f Filter) LessThan(key string, value interface{}) Filter {
	filter := append(f, bson.D{{Key: key, Value: bson.D{{"$lt", value}}}})
	return filter
}

func (f Filter) LessThanEqual(key string, value interface{}) Filter {
	filter := append(f, bson.D{{Key: key, Value: bson.D{{"$lte", value}}}})
	return filter
}

func (f Filter) In(key string, value interface{}) Filter {
	filter := append(f, bson.D{{Key: key, Value: bson.D{{"$in", value}}}})
	return filter
}

func (f Filter) NotIn(key string, value interface{}) Filter {
	filter := append(f, bson.D{{Key: key, Value: bson.D{{"$nin", value}}}})
	return filter
}
