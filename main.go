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
	primitive.ObjectID `bson:"_id,omitempty"`
	Field1             string `bson:"field1,omitempty"`
	Field2             string `bson:"field2,omitempty"`
}

type Aaa struct{}

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
	if err := sample.Insert(ctx, []Document{
		Document{
			Field1: "alice",
			Field2: "bob",
		},
	}); err != nil {
		log.Panic(err)
	}

	// read
	var docs []Document
	if err := sample.Read(ctx, bson.M{"field1": "alice"}, &docs); err != nil {
		log.Panic(err)
	}

	// print
	pp.Println(docs)
}
