package mongolib

import "go.mongodb.org/mongo-driver/bson/primitive"

func NewObjectID() string {
	return primitive.NewObjectID().Hex()
}

func ObjectID(hex string) (objectID primitive.ObjectID) {
	objectID, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return
	}
	return
}
