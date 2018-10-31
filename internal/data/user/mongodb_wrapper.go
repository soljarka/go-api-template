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

// Client is a wrapper for mongo.Client
type Client interface {
	Database(name string, opts ...dbopt.Option) Database
	Connect(ctx context.Context) error
}

// Database is a wrapper for mongo.Database
type Database interface {
	Collection(name string, opts ...collectionopt.Option) Collection
}

// Collection is a wrapper for mocking mongo.Collection
type Collection interface {
	DeleteOne(ctx context.Context, filter interface{}, opts ...deleteopt.Delete) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...deleteopt.Delete) (*mongo.DeleteResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...findopt.One) DocumentResult
	InsertOne(ctx context.Context, document interface{}, opts ...insertopt.One) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, filter interface{}, opts ...findopt.Find) (mongo.Cursor, error)
}

// DocumentResult is a wrapper for mongo.DocumentResult
type DocumentResult interface {
	Decode(v interface{}) error
}

type mongoClient struct {
	*mongo.Client
}

// NewMongoClient wraps mongo.NewClient and returns a MongoDB client instance conveniently wrapped into Client interface
func NewMongoClient(uri string) (Client, error) {
	client, err := mongo.NewClient(uri)
	if err != nil {
		return nil, err
	}
	return &mongoClient{client}, nil
}

func (c mongoClient) Connect(ctx context.Context) error {
	return c.Client.Connect(ctx)
}

func (c mongoClient) Database(name string, opts ...dbopt.Option) Database {
	return &mongoDatabase{ Database: c.Client.Database(name, opts...)}
}

type mongoDatabase struct {
	*mongo.Database
}

func (c *mongoDatabase) Collection(name string, opts ...collectionopt.Option) Collection {
	return &mongoCollection{ Collection: c.Database.Collection(name, opts...)}
}

type mongoCollection struct {
	*mongo.Collection
}

func (c *mongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...deleteopt.Delete) (*mongo.DeleteResult, error) {
	return c.Collection.DeleteOne(ctx, filter, opts...)
}

func (c *mongoCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...deleteopt.Delete) (*mongo.DeleteResult, error) {
	return c.Collection.DeleteMany(ctx, filter, opts...)
}
func (c *mongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...findopt.One) DocumentResult {
	return c.Collection.FindOne(ctx, filter, opts...)
}
func (c *mongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...insertopt.One) (*mongo.InsertOneResult, error) {
	return c.Collection.InsertOne(ctx, document, opts...)
}
func (c *mongoCollection) Find(ctx context.Context, filter interface{}, opts ...findopt.Find) (mongo.Cursor, error) {
	return c.Collection.Find(ctx, filter, opts...)
}
