package mongolib

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNotFound = errors.New("result not found")
)

type Result interface {
	Consume(v interface{}) error
}

type SingleResult struct {
	*mongo.SingleResult
}

func (r *SingleResult) Consume(v interface{}) error {
	err := r.Decode(v)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}
		return err
	}
	return nil
}

type MultipleResult struct {
	*mongo.Cursor
	Error error
}

func (r *MultipleResult) Consume(v interface{}) error {
	if r.Error != nil {
		return r.Error
	}
	err := r.All(context.TODO(), v)
	if err != nil {
		return err
	}
	return nil
}