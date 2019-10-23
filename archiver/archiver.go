package archiver

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/dghubble/go-twitter/twitter"
)

// Archiver represents a handler for the storage provider that archives and
// retrieves the archive.
type Archiver struct {
	mcli   *mongo.Client
	tweets *mongo.Collection
	users  *mongo.Collection
}

func New(mcli *mongo.Client) *Archiver {
	a := &Archiver{mcli: mcli}
	a.tweets = a.mcli.Database("archiver").Collection("tweets")
	a.users = a.mcli.Database("archiver").Collection("users")
	return a
}

// Add returns an error if something goes wrong while adding an user to the
// mongo database.
func (a *Archiver) AddUser(user twitter.User) error {
	_, err := a.users.InsertOne(context.TODO(), user)
	return err
}

// GetUsers will return a slice of users saved in the database.
func (a *Archiver) GetUsers() ([]twitter.User, error) {
	res := []twitter.User{}
	cur, err := a.users.Find(context.TODO(), bson.D{}, options.Find())
	if err != nil {
		return res, err
	}

	for cur.Next(context.TODO()) {
		var item twitter.User
		if err := cur.Decode(&item); err != nil {
			return res, err
		}

		res = append(res, item)
	}

	if err := cur.Err(); err != nil {
		return res, err
	}

	cur.Close(context.TODO())

	return res, nil
}

// Get returns the archived tweets for a given user, it will return an error if
// something goes wrong during the query to the data storage.
func (a *Archiver) Get(id int64) ([]twitter.Tweet, error) {
	res := []twitter.Tweet{}
	cur, err := a.tweets.Find(context.TODO(), bson.M{
		"User": bson.M{
			"ID": id,
		},
	}, options.Find())
	if err != nil {
		return res, err
	}

	for cur.Next(context.TODO()) {
		var item twitter.Tweet
		if err := cur.Decode(&item); err != nil {
			return res, err
		}

		res = append(res, item)
	}

	if err := cur.Err(); err != nil {
		return res, err
	}

	cur.Close(context.TODO())

	return res, nil
}

// Save will return an error if during the save procedure for the tweets passed
// something goes wrong.
func (a *Archiver) Save(tweets []twitter.Tweet) error {
	x := make([]interface{}, len(tweets))
	for k, v := range tweets {
		x[k] = v
	}
	_, err := a.tweets.InsertMany(context.TODO(), x)
	return err
}
