package repository

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
	"time"
	"tourist-alice-skill/internal/api"
)

const (
	userCollection  = "user"
	chatCollection  = "chat_state"
	defaultUser     = "user"
	defaultPassword = "password"

	goLang   = "Golang"
	javaLang = "Java"
)

type UserRepositorySuite struct {
	suite.Suite
	col    *mongo.Collection
	repo   *MongoUserRepository
	closer func()
}

func (s *UserRepositorySuite) SetupSuite() {
	db, cl, err := initTestMongoDb()
	if err != nil {
		s.Error(err)
		return
	}
	s.closer = cl
	s.repo = NewUserRepository(db)
	s.col = db.Collection(userCollection)
}

func (s *UserRepositorySuite) TearDownTest() {

	if err := s.col.Drop(context.Background()); err != nil {
		s.Error(err)
	}
}

func (s *UserRepositorySuite) TearDownSuite() {
	if s.closer != nil {
		s.closer()
	}
}

func TestUserRepositorySuiteUp(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}

func (s *UserRepositorySuite) TestFindUserById_UserExists() {
	//given
	objId := primitive.NewObjectID()

	//when
	_, err := s.col.InsertOne(context.TODO(),
		&api.User{
			ID: objId.Hex(),
		})
	s.NoError(err)

	//then
	u, err := s.repo.FindById(context.TODO(), objId.Hex())
	s.NoError(err)
	s.NotNil(u)
}

func (s *UserRepositorySuite) TestFindUserById_UserNotExists() {
	//when
	u, err := s.repo.FindById(context.TODO(), primitive.NewObjectID().Hex())

	//then
	s.Nil(u)
	s.ErrorIs(err, mongo.ErrNoDocuments)
}

func (s *UserRepositorySuite) TestUpsertUser_UserExists() {
	//given
	objId := primitive.NewObjectID()

	_, err := s.col.InsertOne(context.TODO(),
		&api.User{
			ID:       objId.Hex(),
			UserLang: javaLang,
		})
	s.NoError(err)

	u, err := s.findById(objId)
	s.NoError(err)
	s.Equal(javaLang, u.UserLang)

	//when
	u, err = s.repo.UpsertUser(context.TODO(), api.User{
		ID:       objId.Hex(),
		UserLang: goLang,
	})
	s.NoError(err)
	s.Equal(u.UserLang, goLang)

	//then
	u, err = s.findById(objId)
	s.NoError(err)
	s.Equal(u.UserLang, goLang)
}

func (s *UserRepositorySuite) TestUpsertUser_UserNotExists() {
	//given
	objId := primitive.NewObjectID()

	u, err := s.findById(objId)
	s.ErrorIs(err, mongo.ErrNoDocuments)
	s.Nil(u)

	//when
	u, err = s.repo.UpsertUser(context.TODO(), api.User{
		ID:       objId.Hex(),
		UserLang: goLang,
	})
	s.NoError(err)
	s.Equal(goLang, u.UserLang)

	//then
	u, err = s.findById(objId)
	s.NoError(err)
	s.Equal(goLang, u.UserLang)
}

func (s *UserRepositorySuite) TestUpsertLangUser_UserExists() {

	//given
	objId := primitive.NewObjectID()

	_, err := s.col.InsertOne(context.TODO(),
		&api.User{
			ID:           objId.Hex(),
			SelectedLang: javaLang,
		})
	s.NoError(err)

	u, err := s.findById(objId)
	s.NoError(err)
	s.Equal(javaLang, u.SelectedLang)

	//when
	err = s.repo.UpsertLangUser(context.TODO(), objId.Hex(), javaLang)
	s.NoError(err)

	//then
	u, err = s.findById(objId)
	s.NoError(err)
	s.Equal(goLang, u.SelectedLang)
}

func (s *UserRepositorySuite) TestUpsertLangUser_UserNotExists() {
	//given
	objId := primitive.NewObjectID()

	u, err := s.findById(objId)
	s.ErrorIs(err, mongo.ErrNoDocuments)

	//when
	err = s.repo.UpsertLangUser(context.TODO(), objId.Hex(), goLang)
	s.NoError(err)

	//then
	u, err = s.findById(objId)
	s.NoError(err)
	s.Equal(goLang, u.SelectedLang)
}

func (s *UserRepositorySuite) findById(id primitive.ObjectID) (*api.User, error) {
	res := s.col.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: bson.D{primitive.E{Key: "$eq", Value: id.Hex()}}}})
	if res.Err() != nil {
		return nil, res.Err()
	}
	cs := &api.User{}
	if err := res.Decode(cs); err != nil {
		return nil, err
	}
	return cs, nil
}

type ChatRepositorySuite struct {
	suite.Suite
	col    *mongo.Collection
	repo   *MongoChatStateRepository
	closer func()
}

func (s *ChatRepositorySuite) SetupSuite() {
	db, cl, err := initTestMongoDb()
	if err != nil {
		s.Error(err)
		return
	}
	s.closer = cl
	s.repo = NewChatStateRepository(db)
	s.col = db.Collection(chatCollection)

}

func (s *ChatRepositorySuite) TearDownTest() {
	if err := s.col.Drop(context.Background()); err != nil {
		s.Error(err)
	}
}

func (s *ChatRepositorySuite) TearDownSuite() {
	if s.closer != nil {
		s.closer()
	}
}

func (s *ChatRepositorySuite) TestSave_UserChatStateNotExists_AutoGenerateId() {

	//given
	var (
		action = api.Action("testAction")
		cd     = &api.CallbackData{SelectedCity: "Kazan", Page: 1}
		uId    = "testUser"
	)

	//when
	err := s.repo.Save(context.TODO(), &api.ChatState{
		UserId:       uId,
		Action:       action,
		CallbackData: cd,
	})
	s.NoError(err)

	//then
	res := s.col.FindOne(context.TODO(), bson.D{primitive.E{Key: "user_id", Value: bson.D{primitive.E{Key: "$eq", Value: uId}}}})
	s.NoError(res.Err())
	cs := &api.ChatState{}
	err = res.Decode(cs)
	s.NoError(err)

	s.Equal(uId, cs.UserId)
	s.Equal(action, cs.Action)
	s.Equal(cd, cs.CallbackData)

}

func (s *ChatRepositorySuite) TestSave_UserChatStateNotExists_ManualGenerateId() {
	//given
	var (
		uId    = "testUser"
		action = api.Action("testAction")
		cd     = &api.CallbackData{SelectedCity: "Kazan", Page: 1}
		objId  = primitive.NewObjectID()
	)

	//when
	err := s.repo.Save(context.TODO(), &api.ChatState{
		ID:           objId,
		UserId:       uId,
		Action:       action,
		CallbackData: cd,
	})
	s.NoError(err)

	//then
	res := s.col.FindOne(context.TODO(), bson.D{primitive.E{Key: "user_id", Value: bson.D{primitive.E{Key: "$eq", Value: uId}}}})
	s.NoError(res.Err())
	cs := &api.ChatState{}
	err = res.Decode(cs)
	s.NoError(err)

	s.Equal(objId, cs.ID)
	s.Equal(uId, cs.UserId)
	s.Equal(action, cs.Action)
	s.Equal(cd, cs.CallbackData)
}

func (s *ChatRepositorySuite) TestSave_UserChatStateExists() {
	//given
	objId := primitive.NewObjectID()
	uId := "testsUser"
	testCs := &api.ChatState{
		ID:           objId,
		UserId:       uId,
		Action:       "testAction",
		CallbackData: &api.CallbackData{SelectedCity: "Kazan", Page: 1},
	}
	_, err := s.col.InsertOne(context.TODO(), testCs)
	s.NoError(err)

	//when
	err = s.repo.Save(context.TODO(), testCs)

	//then
	s.NoError(err)
}

func (s *ChatRepositorySuite) TestSave_UserChatStateNotExists_SameUserId() {
	//given
	objId := primitive.NewObjectID()
	uId := "testsUser"

	oldAction := api.Action("oldAction")
	newAction := api.Action("newAction")

	oldChatState := &api.ChatState{
		ID:           objId,
		UserId:       uId,
		Action:       oldAction,
		CallbackData: &api.CallbackData{SelectedCity: "Kazan", Page: 1},
	}
	_, err := s.col.InsertOne(context.TODO(), oldChatState)
	s.NoError(err)

	res := s.col.FindOne(context.TODO(), bson.D{primitive.E{Key: "user_id", Value: bson.D{primitive.E{Key: "$eq", Value: uId}}}})
	s.NoError(res.Err())
	cs := &api.ChatState{}
	err = res.Decode(cs)
	s.NoError(err)
	log.Info().Interface("user", cs).Msg("user")
	s.Equal(oldChatState, cs)

	testCs := &api.ChatState{
		UserId:       uId,
		Action:       newAction,
		CallbackData: &api.CallbackData{SelectedCity: "Kazan", Page: 1},
	}

	//when
	err = s.repo.Save(context.TODO(), testCs)
	s.NoError(err)

	//then
	res = s.col.FindOne(context.TODO(), bson.D{primitive.E{Key: "user_id", Value: bson.D{primitive.E{Key: "$eq", Value: uId}}}})
	s.NoError(res.Err())
	cs = &api.ChatState{}
	err = res.Decode(cs)
	s.NoError(err)

	log.Info().Interface("user", cs).Msg("user")

	s.Equal(newAction, cs.Action)
	s.Equal(objId, cs.ID)
}

func (s *ChatRepositorySuite) TestFindByUserId_UserChatStateExists() {

	//given
	objId := primitive.NewObjectID()
	uId := "testsUser"
	testCs := &api.ChatState{
		ID:           objId,
		UserId:       uId,
		Action:       "testAction",
		CallbackData: &api.CallbackData{SelectedCity: "Kazan", Page: 1},
	}

	_, err := s.col.InsertOne(context.TODO(), testCs)
	s.NoError(err)

	//when
	cs, err := s.repo.FindByUserId(context.TODO(), uId)
	s.NoError(err)

	//then
	s.Equal(testCs, cs)
}

func TestChatRepositorySuiteUp(t *testing.T) {
	suite.Run(t, new(ChatRepositorySuite))
}

func initTestMongoDb() (*mongo.Database, func(), error) {
	ctx := context.TODO()
	mgdb := fmt.Sprintf("repo-test-%d", os.Getpid())
	log.Info().Str("database", mgdb).Msg("Mongo Test Database")
	path, err := os.Getwd()
	log.Info().Str("curr directory", path).Msg("Mongo test directory")
	if err != nil {
		return nil, nil, err
	}
	req := tc.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(2 * time.Minute),
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": defaultUser,
			"MONGO_INITDB_ROOT_PASSWORD": defaultPassword,
			"MONGO_INITDB_DATABASE":      mgdb,
		},
		BindMounts: map[string]string{
			path + "/scripts/mongo-init.js": "/docker-entrypoint-initdb.d/mongo-init.js",
		},
	}

	mgContainer, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	h, err := mgContainer.Host(ctx)
	if err != nil {
		return nil, nil, err
	}

	port, err := mgContainer.MappedPort(ctx, "27017/tcp")

	if err != nil {
		return nil, nil, err
	}

	mongoHost := fmt.Sprintf("mongodb://%s:%s@%s:%d", defaultUser, defaultPassword, h, port.Int())
	log.Info().Str("mongo host", mongoHost).Msg("Staring mongoDB container")
	db, cl, err := initMongoConnection(ctx, mongoHost, mgdb)
	if err != nil {
		return nil, nil, err
	}
	return db, cl, nil
}

func initMongoConnection(ctx context.Context, host string, db string) (*mongo.Database, func(), error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(host))
	if err != nil {
		return nil, nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	return client.Database(db), func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal().Err(err).Msg("error while connect to mongo")
		}
	}, nil
}

func TestMongoChatStateRepository_UniqueIndex(t *testing.T) {

	mg, cl, err := initTestMongoDb()
	assert.NoError(t, err)
	defer cl()
	col := mg.Collection(chatCollection)

	_, err = col.InsertOne(context.TODO(), &api.ChatState{
		ID:           primitive.NewObjectID(),
		UserId:       "testsUser",
		Action:       "testAction",
		CallbackData: &api.CallbackData{SelectedCity: "Kazan", Page: 1},
	})
	assert.NoError(t, err)

	_, err = col.InsertOne(context.TODO(), &api.ChatState{
		ID:           primitive.NewObjectID(),
		UserId:       "testsUser",
		Action:       "testAction",
		CallbackData: &api.CallbackData{SelectedCity: "Kazan", Page: 1},
	})
	log.Error().Err(err).Msg("Constraint")
	assert.Contains(t, err.Error(), "E11000 duplicate key error collection")

}
