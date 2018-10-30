package user

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/collectionopt"
	"github.com/mongodb/mongo-go-driver/mongo/dbopt"
	"github.com/mongodb/mongo-go-driver/mongo/deleteopt"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/insertopt"

	"github.com/mongodb/mongo-go-driver/mongo"
)

type Client interface {
	Database(name string, opts ...dbopt.Option) Database
}

type Database interface {
	Collection(name string, opts ...collectionopt.Option) Collection
}

type Collection interface {
	DeleteOne(ctx context.Context, filter interface{}, opts ...deleteopt.Delete) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...deleteopt.Delete) (*mongo.DeleteResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...findopt.One) DocumentResult
	InsertOne(ctx context.Context, document interface{}, opts ...insertopt.One) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, filter interface{}, opts ...findopt.Find) (mongo.Cursor, error)
}

type DocumentResult interface {
	Decode(v interface{}) error
}

type MongoClient struct {
	*mongo.Client
}

func (c MongoClient) Database(name string, opts ...dbopt.Option) Database {
	return &MongoDatabase{ Database: c.Client.Database(name, opts...)}
}

type MongoDatabase struct {
	*mongo.Database
}

func (c MongoDatabase) Collection(name string, opts ...collectionopt.Option) Collection {
	return &MongoCollection{ Collection: c.Database.Collection(name, opts...)}
}

type MongoCollection struct {
	*mongo.Collection
}

func (c *MongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...deleteopt.Delete) (*mongo.DeleteResult, error) {
	return c.Collection.DeleteOne(ctx, filter, opts...)
}
func (c *MongoCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...deleteopt.Delete) (*mongo.DeleteResult, error) {
	return c.Collection.DeleteMany(ctx, filter, opts...)
}
func (c *MongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...findopt.One) DocumentResult {
	return c.Collection.FindOne(ctx, filter, opts...)
}
func (c *MongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...insertopt.One) (*mongo.InsertOneResult, error) {
	return c.Collection.InsertOne(ctx, document, opts...)
}
func (c *MongoCollection) Find(ctx context.Context, filter interface{}, opts ...findopt.Find) (mongo.Cursor, error) {
	return c.Collection.Find(ctx, filter, opts...)
}
