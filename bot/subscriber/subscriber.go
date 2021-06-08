package subscriber

import (
	"github.com/DipandaAser/tweetwatcher/config"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"time"
)

type Chat struct {
	ID        string `bson:"_id"`
	Name      string `bson:"Name"`
	Activate  bool   `bson:"Activate"`
	CreatedAt string `bson:"CreatedAt"`
}

func (c Chat) Recipient() int64 {

	atoi, _ := strconv.Atoi(c.ID)
	return int64(atoi)
}

var collectionName = "Subscriber"

func New(name string, chatId string) (*Chat, error) {

	t, _ := time.Now().UTC().MarshalText()
	theChat := Chat{
		ID:        chatId,
		Name:      name,
		Activate:  true,
		CreatedAt: string(t),
	}

	_, err := config.DB.Collection(collectionName).InsertOne(*config.MongoCtx, theChat)
	if err != nil {
		return nil, err
	}

	return &theChat, nil
}

func IsExist(chatId string) bool {

	chat := &Chat{}
	filter := bson.M{"_id": chatId}
	err := config.DB.Collection(collectionName).FindOne(*config.MongoCtx, filter).Decode(chat)
	if err != nil {
		return false
	}

	return true
}

func IsActivate(chatId string) bool {

	chat := &Chat{}
	filter := bson.M{"_id": chatId}
	err := config.DB.Collection(collectionName).FindOne(*config.MongoCtx, filter).Decode(chat)
	if err != nil {
		return false
	}

	return chat.Activate
}

func DesactivateSubscriber(chatId string) error {

	filter := bson.M{"_id": chatId}
	updates := bson.M{"$set": bson.M{"Activate": false}}
	result := config.DB.Collection(collectionName).FindOneAndUpdate(*config.MongoCtx, filter, updates)
	if err := result.Err(); err != nil {
		return err
	}
	return nil
}

func ActivateSubscriber(chatId string) error {

	filter := bson.M{"_id": chatId}
	updates := bson.M{"$set": bson.M{"Activate": true}}
	result := config.DB.Collection(collectionName).FindOneAndUpdate(*config.MongoCtx, filter, updates)
	if err := result.Err(); err != nil {
		return err
	}
	return nil
}

func GetActivatedSubscribers() ([]Chat, error) {

	filter := bson.M{"Activate": true}
	cur, err := config.DB.Collection(collectionName).Find(*config.MongoCtx, filter)
	if err != nil {
		return nil, err
	}

	chats := []Chat{}
	err = cur.All(*config.MongoCtx, &chats)
	if err != nil {
		return nil, err
	}

	return chats, nil
}
