package repository

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"tourist-alice-skill/internal/api"
)

type UserRepository interface {
	UpsertUser(ctx context.Context, u api.User) (*api.User, error)
	FindById(ctx context.Context, id string) (*api.User, error)
}

type ChatStateRepository interface {
	Save(ctx context.Context, u *api.ChatState) error
	FindByUserId(ctx context.Context, userId string) (*api.ChatState, error)
	DeleteById(ctx context.Context, id primitive.ObjectID) error
}

type MongoUserRepository struct {
	col *mongo.Collection
}

type MongoChatStateRepository struct {
	col *mongo.Collection
}

func NewUserRepository(col *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{col: col.Collection("user")}
}

func NewChatStateRepository(col *mongo.Database) *MongoChatStateRepository {
	return &MongoChatStateRepository{col: col.Collection("chat_state")}
}

func (r MongoUserRepository) FindById(ctx context.Context, id string) (*api.User, error) {
	res := r.col.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: eqFilter(id)}})
	if res.Err() != nil {
		return nil, res.Err()
	}
	cs := &api.User{}
	if err := res.Decode(cs); err != nil {
		return nil, err
	}
	return cs, nil
}

func (r MongoUserRepository) UpsertUser(ctx context.Context, u api.User) (*api.User, error) {
	opts := options.Update().SetUpsert(true)
	f := bson.D{primitive.E{Key: "_id", Value: eqFilter(u.ID)}}
	update := bson.D{primitive.E{Key: "$set", Value: u}}
	_, err := r.col.UpdateOne(ctx, f, update, opts)
	if err != nil {
		return nil, err
	}
	return r.FindById(ctx, u.ID)
}

func (r MongoUserRepository) UpsertLangUser(ctx context.Context, userId string, lang string) error {
	opts := options.Update().SetUpsert(true)
	f := bson.D{primitive.E{Key: "_id", Value: eqFilter(userId)}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.M{"selected_lang": lang}}}
	_, err := r.col.UpdateOne(ctx, f, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func (csr MongoChatStateRepository) Save(ctx context.Context, cs *api.ChatState) error {
	res, err := csr.col.InsertOne(ctx, cs)
	if err != nil {
		log.Error().Err(err).Msg("insert failed")
	}
	if res != nil && res.InsertedID == nil {
		return errors.New("insert failed")
	}
	return err
}

func (csr MongoChatStateRepository) FindByUserId(ctx context.Context, userId string) (*api.ChatState, error) {
	res := csr.col.FindOne(ctx, bson.D{primitive.E{Key: "user_id", Value: eqFilter(userId)}})
	if res.Err() == mongo.ErrNoDocuments {
		log.Debug().Err(res.Err()).Msgf("chat_state not found by user_id %v", userId)
		return nil, nil
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	cs := &api.ChatState{}
	if err := res.Decode(cs); err != nil {
		return nil, err
	}
	return cs, nil
}

func (csr MongoChatStateRepository) DeleteById(ctx context.Context, id primitive.ObjectID) error {
	_, err := csr.col.DeleteOne(ctx, bson.D{primitive.E{Key: "_id", Value: eqFilter(id)}})
	if err != nil {
		log.Error().Err(err).Msg("delete failed")
		return err
	}
	return nil
}

func eqFilter(id interface{}) bson.D {
	return bson.D{primitive.E{Key: "$eq", Value: id}}
}
