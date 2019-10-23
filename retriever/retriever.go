// Package retriever handles getting information from Twitter.
package retriever

import (
	"errors"

	"github.com/dghubble/go-twitter/twitter"
)

// Retriever represents a handler for retrieving information from Twitter.
type Retriever struct {
	client  *twitter.Client
	targets []Target
}

// Target represents a target to be saved into the archive, this is a Twitter
// user ID, screen name and Last tweet we know of.
type Target struct {
	ID       int64
	Username string
	Last     int64
}

var (
	ErrProtected = errors.New("user timeline is protected")
)

// New returns a Retriever pointer ready to use.
func New(client *twitter.Client) *Retriever {
	return &Retriever{
		client: client,
	}
}

// Latest returns an array of tweets for the given user ID. It will return an
// error if something goes wrong during the request.
func (r *Retriever) Latest(id int64) ([]twitter.Tweet, error) {
	t, _, e := r.client.Timelines.UserTimeline(&twitter.UserTimelineParams{UserID: id})
	return t, e
}

// Since returns an array of tweets for the given user ID since the tweet on
// the second argument. If there's an error it will return it in the second
// parameter.
func (r *Retriever) Since(user, tweet int64) ([]twitter.Tweet, error) {
	t, _, e := r.client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		UserID:  user,
		SinceID: tweet,
	})
	return t, e
}

// From returns an array of tweets for the given user ID and uses the ID passed
// for the tweet to mark the moment in the past that sets the initial point to
// read into the past. This works the other way around than Since().
func (r *Retriever) From(user, tweet int64) ([]twitter.Tweet, error) {
	t, _, e := r.client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		UserID: user,
		MaxID:  tweet,
	})
	return t, e
}

// UserID returns the numerical and unmodifiable ID of a given screen name. If
// the error is ErrProtected it will anyway return the user numerical ID but
// you need to follow that user in order to be able to retrieve tweets.
func (r *Retriever) UserID(username string) (int64, error) {
	u, _, e := r.client.Users.Show(&twitter.UserShowParams{ScreenName: username})
	if u.Protected {
		return u.ID, ErrProtected
	}
	return u.ID, e
}
