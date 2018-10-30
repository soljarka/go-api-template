package user

import (
	"context"
	"errors"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/collectionopt"
	"github.com/mongodb/mongo-go-driver/mongo/dbopt"
	"github.com/mongodb/mongo-go-driver/mongo/deleteopt"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/insertopt"
	"testing"
)

type clientMock struct {}

func (c *clientMock) Database(name string, opts ...dbopt.Option) Database {
	return &databaseMock{}
}

type databaseMock struct {}

func (d *databaseMock) Collection(name string, opts ...collectionopt.Option) Collection {
	return &collectionMock{}
}

type collectionMock struct {}

func (c *collectionMock) DeleteOne(ctx context.Context, filter interface{}, opts ...deleteopt.Delete) (*mongo.DeleteResult, error) {
	if ctx.Value("error") == true {
		return nil, errors.New("error")
	}
	return nil, nil
}

func (c *collectionMock) DeleteMany(ctx context.Context, filter interface{}, opts ...deleteopt.Delete) (*mongo.DeleteResult, error) {
	if ctx.Value("error") == true {
		return nil, errors.New("error")
	}
	return nil, nil
}

func (c *collectionMock) FindOne(ctx context.Context, filter interface{}, opts ...findopt.One) DocumentResult {
	return &documentResultMock{ctx}
}

func (c *collectionMock) InsertOne(ctx context.Context, document interface{}, opts ...insertopt.One) (*mongo.InsertOneResult, error) {
	if ctx.Value("error") == true {
		return nil, errors.New("error")
	}
	return &mongo.InsertOneResult{InsertedID: objectid.New()}, nil
}

func (c *collectionMock) Find(ctx context.Context, filter interface{}, opts ...findopt.Find) (mongo.Cursor, error) {
	if ctx.Value("error") == true {
		return nil, errors.New("error")
	}
	return &cursorMock{2}, nil
}

type cursorMock struct {
	Count int
}

func (c *cursorMock) ID() int64 {return 1}

func (c *cursorMock) Next(context.Context) bool {
	if c.Count > 0 {
		c.Count = c.Count - 1
		return true
	}
	return false
}

func (c *cursorMock) Decode(data interface{}) error {
	data = mongoUser{ID: objectid.New(), Name: "Foo", Surname: "Bar"}
	return nil
}

func (c *cursorMock) DecodeBytes() (bson.Reader, error) {return nil, nil}

func (c *cursorMock) Err() error {return nil}

func (c *cursorMock) Close(context.Context) error {return nil}

type documentResultMock struct {
	ctx context.Context
}

func (d *documentResultMock) Decode(v interface{}) error {
	if d.ctx.Value("error") == true {
		return errors.New("error")
	}
	v = mongoUser{ID: objectid.New(), Name: "Foo", Surname: "Bar"}
	return nil
}

func Test_All_Ok(t *testing.T) {
	client := clientMock{}
	repo := NewMongoRepository(&client, "", "")
	ctx := context.WithValue(context.Background(), "error", false)
	users, err := repo.All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if users == nil {
		t.Error("should be slice of users, got nil")
	}

	if len(users) != 2 {
		t.Errorf("should be 2 users, got %v", len(users))
	}
}

func Test_All_Err(t *testing.T) {
	client := clientMock{}
	repo := NewMongoRepository(&client, "", "")
	ctx := context.WithValue(context.Background(), "error", true)
	_, err := repo.All(ctx)
	if err == nil {
		t.Error("should return error, returned nil")
	}
}

func Test_Delete_Ok(t *testing.T) {
	client := clientMock{}
	repo := NewMongoRepository(&client, "", "")
	ctx := context.WithValue(context.Background(), "error", false)
	err := repo.Delete(ctx, objectid.New().Hex())
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Delete_Err(t *testing.T) {
	client := clientMock{}
	repo := NewMongoRepository(&client, "", "")
	ctx := context.WithValue(context.Background(), "error", true)
	err := repo.Delete(ctx, "1")
	if err == nil {
		t.Error("should return error, returned nil")
	}
}

func Test_DeleteAll_Ok(t *testing.T) {
	client := clientMock{}
	repo := NewMongoRepository(&client, "", "")
	ctx := context.WithValue(context.Background(), "error", false)
	err := repo.DeleteAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_DeleteAll_Err(t *testing.T) {
	client := clientMock{}
	repo := NewMongoRepository(&client, "", "")
	ctx := context.WithValue(context.Background(), "error", true)
	err := repo.DeleteAll(ctx)
	if err == nil {
		t.Error("should return error, returned nil")
	}
}

func Test_Find_Ok(t *testing.T) {
	client := clientMock{}
	repo := NewMongoRepository(&client, "", "")
	ctx := context.WithValue(context.Background(), "error", false)
	user, err := repo.Find(ctx, objectid.New().Hex())
	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Error("should return user, returned nil")
	}
}

func Test_Find_Err(t *testing.T) {
	client := clientMock{}
	repo := NewMongoRepository(&client, "", "")
	ctx := context.WithValue(context.Background(), "error", true)
	_, err := repo.Find(ctx, objectid.New().Hex())
	if err == nil {
		t.Error("should return error, returned nil")
	}
}

func Test_Save_Ok(t *testing.T) {
	client := clientMock{}
	repo := NewMongoRepository(&client, "", "")
	ctx := context.WithValue(context.Background(), "error", false)
	newUser := User{Name: "foo", Surname: "bar"}
	id, err := repo.Save(ctx, &newUser)
	if err != nil {
		t.Fatal(err)
	}

	_, err = objectid.FromHex(id)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Save_Err(t *testing.T) {
	client := clientMock{}
	repo := NewMongoRepository(&client, "", "")
	ctx := context.WithValue(context.Background(), "error", true)
	newUser := User{Name: "foo", Surname: "bar"}
	_, err := repo.Save(ctx, &newUser)
	if err == nil {
		t.Error("should return error, returned nil")
	}
}