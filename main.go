package main

import (
	"context"
	"log"
	"mongodb-go-sample/db"
	"time"

	"github.com/k0kubun/pp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const mongodbURI string = "mongodb://localhost:27017"

type Document struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Field1 string             `bson:"field1,omitempty"`
	Field2 string             `bson:"field2,omitempty"`
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	client, err := db.NewClient(ctx, mongodbURI)
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if err = client.DisconnectDB(ctx); err != nil {
			log.Panic(err)
		}
	}()

	sample := client.NewDB("db").NewCollection("sample", Document{})

	// write
	insert(ctx, sample, []Document{
		Document{
			Field1: "111",
			Field2: "aaa",
		},
		Document{
			Field1: "111",
			Field2: "bbb",
		},
	})

	// read
	var docs []Document
	read(ctx, sample, docs, bson.M{})

	update(ctx, sample,
		bson.M{"field1": "111"},
		bson.D{
			{"$set", bson.D{{"field1", "xxx"}}},
		})

	read(ctx, sample, docs, bson.M{"field1": "xxx"})
}

func insert(ctx context.Context, c db.Collection, docs []Document) {
	if err := c.Insert(ctx, docs); err != nil {
		log.Panic(err)
	}

}
func update(ctx context.Context, c db.Collection, filter, update interface{}) {
	c.Update(ctx, filter, update)
}

func read(ctx context.Context, c db.Collection, docs []Document, filter interface{}) {
	if err := c.Read(ctx, filter, &docs); err != nil {
		log.Panic(err)
	}
	pp.Println(docs)
}
