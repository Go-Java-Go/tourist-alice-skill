//+build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"tourist-alice-skill/internal/events"
	"tourist-alice-skill/internal/handler"
	"tourist-alice-skill/internal/repository"
	"tourist-alice-skill/internal/service"
	"tourist-alice-skill/internal/skill"
)

func initApp(ctx context.Context, cfg *config) (tg *events.AliceListener, closer func(), err error) {
	wire.Build(initMongoConnection, initSkillConfig, initIziTravelClient,
		service.NewUserService, service.NewChatStateService, handler.NewErrorHandler,
		wire.Bind(new(events.ChatStateService), new(*service.ChatStateService)),
		wire.Bind(new(skill.ChatStateService), new(*service.ChatStateService)),
		wire.Bind(new(events.UserService), new(*service.UserService)),
		ProvideBotList, bots,
		repository.NewUserRepository, wire.Bind(new(repository.UserRepository), new(*repository.MongoUserRepository)),
		repository.NewChatStateRepository, wire.Bind(new(repository.ChatStateRepository), new(*repository.MongoChatStateRepository)))
	return nil, nil, nil
}

var bots = wire.NewSet(skill.NewStartScreen, skill.NewOperationScreen)

func ProvideBotList(s1 *skill.StartScreen, s2 *skill.OperationScreen) []skill.Interface {
	return []skill.Interface{s1, s2}
}
