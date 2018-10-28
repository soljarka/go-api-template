package user

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type mongoUser struct {
	ID      objectid.ObjectID `bson:"_id"`
	Name    string            `bson:"name"`
	Surname string            `bson:"surname"`
}

type mongoUserInput struct {
	Name    string `bson:"name"`
	Surname string `bson:"surname"`
}

type repo struct {
	collection *mongo.Collection
}

// NewMongoRepository returns a MongoDB implementation of User repository
func NewMongoRepository(c *mongo.Client, db string, collection string) Repository {
	return &repo{
		collection: c.Database(db).Collection(collection),
	}
}

func (r *repo) Delete(ctx context.Context, id ID) error {
	oid, err := objectid.FromHex(id)
	if err != nil {
		return err
	}

	filter := bson.NewDocument(bson.EC.ObjectID("_id", oid))
	_, err = r.collection.DeleteOne(ctx, filter)

	return err
}

func (r *repo) DeleteAll(ctx context.Context) error {
	_, err := r.collection.DeleteMany(ctx, nil)

	return err
}

func (r *repo) Find(ctx context.Context, id ID) (*User, error) {
	result := mongoUser{}
	oid, err := objectid.FromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.NewDocument(bson.EC.ObjectID("_id", oid))
	err = r.collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:      result.ID.Hex(),
		Name:    result.Name,
		Surname: result.Surname,
	}, nil
}

func (r *repo) Save(ctx context.Context, user *User) (ID, error) {
	mongoUser := mongoUserInput{
		Name:    user.Name,
		Surname: user.Surname,
	}

	res, err := r.collection.InsertOne(ctx, mongoUser)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(objectid.ObjectID).Hex(), nil
}

func (r *repo) All(ctx context.Context) ([]*User, error) {
	cur, err := r.collection.Find(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var users []*User

	for cur.Next(ctx) {
		usr := mongoUser{}
		err = cur.Decode(&usr)
		if err != nil {
			return nil, err
		}
		users = append(users, &User{
			ID:      usr.ID.Hex(),
			Name:    usr.Name,
			Surname: usr.Surname,
		})
	}

	if err = cur.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
