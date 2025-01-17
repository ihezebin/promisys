package storage

import (
	"context"
	"testing"

	"github.com/ihezebin/promisys/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMongo(t *testing.T) {
	ctx := context.Background()
	err := InitMongoStorageClient(ctx, "mongodb://root:root@localhost:27017/promisys?authSource=admin")
	if err != nil {
		t.Fatal(err)
	}

	collection := MongoStorageDatabase().Collection("example")
	_, err = collection.InsertOne(ctx, &entity.Example{
		Id:       primitive.NewObjectID().Hex(),
		Username: "admin",
		Password: "123456",
		Email:    "6wqz8@example.com",
		Salt:     "123456",
	})
	if err != nil {
		t.Fatal(err)
	}

	examples := make([]*entity.Example, 0)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		t.Fatal(err)
	}

	err = cursor.All(ctx, &examples)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(examples)
}
