package twitter

import (
	"github.com/DipandaAser/tweetwatcher/config"
	"go.mongodb.org/mongo-driver/bson"
)

type Tweet struct {
	ID        string `bson:"_id"`
	User      string `bson:"User"`
	Text      string `bson:"Text"`
	Url       string `bson:"Url"`
	CreatedAt string `bson:"CreatedAt"`
}

var collectionName string = "Tweets"

func Save(id string, text string, user string, url string, pubDate string) (*Tweet, error) {

	theTweet := Tweet{
		ID:        id,
		User:      user,
		Text:      text,
		Url:       url,
		CreatedAt: pubDate,
	}

	_, err := config.DB.Collection(collectionName).InsertOne(*config.MongoCtx, theTweet)
	if err != nil {
		return nil, err
	}

	return &theTweet, nil
}

// IsExist check if a tweet are already been scrapped
func IsExist(tweetId string) bool {

	tweet := &Tweet{}
	filter := bson.M{"_id": tweetId}

	err := config.DB.Collection(collectionName).FindOne(*config.MongoCtx, filter).Decode(tweet)
	if err != nil {
		return false
	}

	return true
}
