package mongolib

import "go.mongodb.org/mongo-driver/bson"

func Regex(value interface{}) bson.D {
	return bson.D{{"$regex", value}}
}
