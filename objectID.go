package mongolib

import "go.mongodb.org/mongo-driver/bson/primitive"

func NewObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func ObjectID(hex string) (objectID primitive.ObjectID) {
	objectID, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return
	}
	return
}
