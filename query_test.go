package mongolib

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/strikesecurity/strikememongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
)

type person struct {
	ID    primitive.ObjectID `bson:"_id"`
	Name  string             `bson:"name"`
	Age   int                `bson:"age"`
	Car   car                `bson:"car"`
	Score []int              `bson:"score"`
	Alias []string           `bson:"alias"`
}

type car struct {
	Color string `bson:"color"`
	Speed int    `bson:"speed"`
}

func initTest(t *testing.T) *Database {
	const (
		version = "4.0.5"
		dbName = "db"
	)

	srv, err := strikememongo.Start(version)
	assert.NoError(t, err)
	uri := srv.URI()

	c, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	assert.NoError(t, err)

	err = c.Ping(context.Background(), readpref.Primary())
	assert.NoError(t, err)

	return NewDatabase(c, dbName)
}

func TestQuery_Set(t *testing.T) {
	const collName = "coll"
	db := initTest(t)
	var (
		ctx = context.Background()
		coll = db.Coll(collName)
		before = person{
			ID:   NewObjectID(),
			Name: "Trevor",
			Age:  27,
			Car: car{
				Color: "red",
			},
		}
		after = person{
			ID:   before.ID,
			Name: "Ali",
			Age:  27,
			Car: car{
				Color: "blue",
			},
		}
	)

	tests := []struct {
		name    string
		prepare func()
		action  func()
		assert  func()
		wantErr bool
	}{
		{
			name: "success: update set attribute",
			prepare: func() {
				err := coll.Save(ctx, before.ID, before)
				assert.NoError(t, err)
			},
			action: func() {
				err := coll.Query().
					Set("name", after.Name).
					Update(ctx)
				assert.NoError(t, err)
			},
			assert: func() {
				var result person
				err := coll.Query().FindOne(ctx).Consume(&result)
				assert.NoError(t, err)
				assert.Equal(t, after.Name, result.Name)
			},
			wantErr: false,
		},
		{
			name: "success: update set nested attribute",
			prepare: func() {
				err := coll.Save(ctx, before.ID, before)
				assert.NoError(t, err)
			},
			action: func() {
				err := coll.Query().
					Set("car.color", after.Car.Color).
					Update(ctx)
				assert.NoError(t, err)
			},
			assert: func() {
				var result person
				err := coll.Query().FindOne(ctx).Consume(&result)
				assert.NoError(t, err)
				assert.Equal(t, after.Car.Color, result.Car.Color)
			},
			wantErr: false,
		},
		{
			name: "success: update set multiple attribute",
			prepare: func() {
				err := coll.Save(ctx, before.ID, before)
				assert.NoError(t, err)
			},
			action: func() {
				err := coll.Query().
					Set("name", after.Name).
					Set("car.color", after.Car.Color).
					Update(ctx)
				assert.NoError(t, err)
			},
			assert: func() {
				var result person
				err := coll.Query().FindOne(ctx).Consume(&result)
				assert.NoError(t, err)
				assert.Equal(t, after.Name, result.Name)
				assert.Equal(t, after.Car.Color, result.Car.Color)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := coll.DeleteMany(context.Background(), options.Delete())
			assert.NoError(t, err)
			if tt.prepare != nil {
				tt.prepare()
			}
			if tt.action != nil {
				tt.action()
			}
			if tt.assert != nil {
				tt.assert()
			}
		})
	}
}

func TestQuery_Push(t *testing.T) {
	const collName = "coll"
	db := initTest(t)
	var (
		ctx = context.Background()
		coll = db.Coll(collName)
		newScore = 89
		before = person{
			ID:   NewObjectID(),
			Name: "Trevor",
			Age:  27,
			Score: []int{93,80,13},
			Alias: []string{"Joker"},
		}
		after = person{
			ID:   before.ID,
			Name: "Trevor",
			Age:  27,
			Score: []int{93,80,13,newScore},
			Alias: []string{"Joker", "Batman"},
		}
	)

	tests := []struct {
		name    string
		prepare func()
		action  func()
		assert  func()
		wantErr bool
	}{
		{
			name: "success: update push array attribute",
			prepare: func() {
				err := coll.Save(ctx, before.ID, before)
				assert.NoError(t, err)
			},
			action: func() {
				err := coll.Query().
					Push("score", newScore).
					Update(ctx)
				assert.NoError(t, err)
			},
			assert: func() {
				var result person
				err := coll.Query().FindOne(ctx).Consume(&result)
				assert.NoError(t, err)
				assert.Equal(t, after.Score, result.Score)
			},
			wantErr: false,
		},
		{
			name: "success: update push array attribute multiple value",
			prepare: func() {
				err := coll.Save(ctx, before.ID, before)
				assert.NoError(t, err)
			},
			action: func() {
				err := coll.Query().
					Push("score", newScore, 1).
					Update(ctx)
				assert.NoError(t, err)
			},
			assert: func() {
				var result person
				err := coll.Query().FindOne(ctx).Consume(&result)
				assert.NoError(t, err)
				assert.Equal(t, append(after.Score, 1), result.Score)
			},
			wantErr: false,
		},
		{
			name: "success: update push array multiple attribute",
			prepare: func() {
				err := coll.Save(ctx, before.ID, before)
				assert.NoError(t, err)
			},
			action: func() {
				err := coll.Query().
					Push("score", newScore).
					Push("alias", "Batman").
					Update(ctx)
				assert.NoError(t, err)
			},
			assert: func() {
				var result person
				err := coll.Query().FindOne(ctx).Consume(&result)
				assert.NoError(t, err)
				assert.Equal(t, after.Score, result.Score)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := coll.DeleteMany(context.Background(), options.Delete())
			assert.NoError(t, err)
			if tt.prepare != nil {
				tt.prepare()
			}
			if tt.action != nil {
				tt.action()
			}
			if tt.assert != nil {
				tt.assert()
			}
		})
	}
}

func TestQuery_Pull(t *testing.T) {
	const collName = "coll"
	db := initTest(t)
	var (
		ctx = context.Background()
		coll = db.Coll(collName)
		deletedScore = 89
		before = person{
			ID:   NewObjectID(),
			Name: "Trevor",
			Age:  27,
			Score: []int{93, 80, deletedScore, 13, deletedScore},
			Alias: []string{"Joker", "Batman"},
		}
		after = person{
			ID:   before.ID,
			Name: "Trevor",
			Age:  27,
			Score: []int{93,80,13},
			Alias: []string{"Joker", "Batman"},
		}
	)

	tests := []struct {
		name    string
		prepare func()
		action  func()
		assert  func()
		wantErr bool
	}{
		{
			name: "success: update pull array attribute",
			prepare: func() {
				err := coll.Save(ctx, before.ID, before)
				assert.NoError(t, err)
			},
			action: func() {
				err := coll.Query().
					Pull("score", deletedScore).
					Update(ctx)
				assert.NoError(t, err)
			},
			assert: func() {
				var result person
				err := coll.Query().FindOne(ctx).Consume(&result)
				assert.NoError(t, err)
				assert.Equal(t, after.Score, result.Score)
			},
			wantErr: false,
		},
		{
			name: "success: update pull array attribute multiple value",
			prepare: func() {
				err := coll.Save(ctx, before.ID, before)
				assert.NoError(t, err)
			},
			action: func() {
				err := coll.Query().
					Pull("score", deletedScore, 13, 93).
					Update(ctx)
				assert.NoError(t, err)
			},
			assert: func() {
				var result person
				err := coll.Query().FindOne(ctx).Consume(&result)
				assert.NoError(t, err)
				assert.Equal(t, []int{80}, result.Score)
			},
			wantErr: false,
		},
		{
			name: "success: update pull array multiple attribute",
			prepare: func() {
				err := coll.Save(ctx, before.ID, before)
				assert.NoError(t, err)
			},
			action: func() {
				err := coll.Query().
					Pull("score", deletedScore).
					Pull("alias", "Joker").
					Update(ctx)
				assert.NoError(t, err)
			},
			assert: func() {
				var result person
				err := coll.Query().FindOne(ctx).Consume(&result)
				assert.NoError(t, err)
				assert.Equal(t, after.Score, result.Score)
				assert.Equal(t, []string{"Batman"}, result.Alias)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := coll.DeleteMany(context.Background(), options.Delete())
			assert.NoError(t, err)
			if tt.prepare != nil {
				tt.prepare()
			}
			if tt.action != nil {
				tt.action()
			}
			if tt.assert != nil {
				tt.assert()
			}
		})
	}
}

func TestQuery_Inc(t *testing.T) {
	const collName = "coll"
	db := initTest(t)
	var (
		ctx = context.Background()
		coll = db.Coll(collName)
		before = person{
			ID:   NewObjectID(),
			Name: "Trevor",
			Age:  27,
			Score: []int{93,80,13},
			Car: car{
				Color: "red",
				Speed: 10,
			},
		}
		after = person{
			ID:   before.ID,
			Name: "Trevor",
			Age:  28,
			Score: []int{93,80,13},
			Car: car{
				Color: "red",
				Speed: 15,
			},
		}
	)

	tests := []struct {
		name    string
		prepare func()
		action  func()
		assert  func()
		wantErr bool
	}{
		{
			name: "success: update increase attribute",
			prepare: func() {
				err := coll.Save(ctx, before.ID, before)
				assert.NoError(t, err)
			},
			action: func() {
				err := coll.Query().
					Inc("age", 1).
					Update(ctx)
				assert.NoError(t, err)
			},
			assert: func() {
				var result person
				err := coll.Query().FindOne(ctx).Consume(&result)
				assert.NoError(t, err)
				assert.Equal(t, after.Age, result.Age)
			},
			wantErr: false,
		},
		{
			name: "success: update decrease attribute",
			prepare: func() {
				err := coll.Save(ctx, before.ID, before)
				assert.NoError(t, err)
			},
			action: func() {
				err := coll.Query().
					Inc("age", -2).
					Update(ctx)
				assert.NoError(t, err)
			},
			assert: func() {
				var result person
				err := coll.Query().FindOne(ctx).Consume(&result)
				assert.NoError(t, err)
				assert.Equal(t, before.Age-2, result.Age)
			},
			wantErr: false,
		},
		{
			name: "success: update increase multiple attribute",
			prepare: func() {
				err := coll.Save(ctx, before.ID, before)
				assert.NoError(t, err)
			},
			action: func() {
				err := coll.Query().
					Inc("age", 1).
					Inc("car.speed", 5).
					Update(ctx)
				assert.NoError(t, err)
			},
			assert: func() {
				var result person
				err := coll.Query().FindOne(ctx).Consume(&result)
				assert.NoError(t, err)
				assert.Equal(t, after.Age, result.Age)
				assert.Equal(t, after.Car.Speed, result.Car.Speed)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := coll.DeleteMany(context.Background(), options.Delete())
			assert.NoError(t, err)
			if tt.prepare != nil {
				tt.prepare()
			}
			if tt.action != nil {
				tt.action()
			}
			if tt.assert != nil {
				tt.assert()
			}
		})
	}
}