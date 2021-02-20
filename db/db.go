package db

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	*mongo.Client
}

func NewClient(ctx context.Context, URI string) (Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	return Client{client}, err
}

func (c Client) DisconnectDB(ctx context.Context) error {
	return c.Disconnect(ctx)
}

func (c Client) NewDB(name string) DB {
	return DB{c.Database(name)}
}

type DB struct {
	*mongo.Database
}

func (db DB) NewCollection(name string, doc interface{}) Collection {
	return Collection{db.Collection(name), reflect.TypeOf(doc)}
}

type Collection struct {
	*mongo.Collection
	docType reflect.Type
}

func (c Collection) Insert(ctx context.Context, docs interface{}) error {
	if reflect.TypeOf(docs) != reflect.SliceOf(c.docType) {
		return fmt.Errorf("Error: type of docs is invalid, %#v\n", docs)
	}

	switch reflect.TypeOf(docs).Kind() {
	case reflect.Slice:
		v := reflect.ValueOf(docs)
		for i := 0; i < v.Len(); i++ {
			if _, err := c.InsertOne(ctx, v.Index(i).Convert(c.docType).Interface()); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("Error: docs is not a slice, %#v\n", docs)
	}

	return nil
}

func (c Collection) Read(ctx context.Context, filter interface{}, docs interface{}) error {
	if reflect.TypeOf(docs) != reflect.PtrTo(reflect.SliceOf(c.docType)) {
		return fmt.Errorf("Error: type of docs is invalid, %#v\n", docs)
	}
	cursor, err := c.Find(ctx, filter)
	if err != nil {
		return err
	}
	if err := cursor.All(ctx, docs); err != nil {
		return err
	}
	return nil
}
