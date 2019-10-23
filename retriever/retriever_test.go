package retriever

import (
	"os"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func setup() *Retriever {
	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_KEY"), os.Getenv("TWITTER_ACCESS_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	return New(client)
}

func TestNew(t *testing.T) {
	n := New(nil)
	if n == nil {
		t.Error("result of New shouldn't be nil")
	}
}

func TestLatest(t *testing.T) {
	r := setup()
	tw, err := r.Latest(879000394223022080)
	if err != nil {
		t.Errorf("retrieving tweets: %q", err)
	}
	if len(tw) == 0 {
		t.Error("length of the result shouldn't be zero")
	}
}

func TestUserID(t *testing.T) {
	r := setup()
	_, err := r.UserID("packethost")
	if err != nil {
		t.Errorf("retrieving user id: %q", err)
	}
	_, err = r.UserID("ghostbar")
	if err == nil || err != ErrProtected {
		t.Errorf("error should return ErrProtected")
	}
}
