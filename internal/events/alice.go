package events

import (
	"context"
	"github.com/azzzak/alice"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"tourist-alice-skill/internal/api"
	"tourist-alice-skill/internal/skill"
)

type UserService interface {
	UpsertUser(ctx context.Context, u api.User) (*api.User, error)
}

type ChatStateService interface {
	FindByUserId(ctx context.Context, userId string) (*api.ChatState, error)
}

// AliceListener listens to alice update, forward to skill bots and send back responses
type AliceListener struct {
	AliceAPI         aliceAPI
	Bots             skill.Interface
	UserService      UserService
	ChatStateService ChatStateService
}

type aliceAPI interface {
	ListenForWebhook(addr string, opts ...func(*alice.Options)) alice.Stream
}

// Do process all events
func (l *AliceListener) Do(ctx context.Context) (err error) {
	updates := alice.ListenForWebhook("/hook", alice.Timeout(250000))

	updates.Loop(func(kit alice.Kit) *alice.Response {
		upd := &api.Update{Request: kit.Req, Response: kit.Resp}
		upd.User, err = l.UserService.UpsertUser(ctx, api.User{ID: kit.Req.Session.UserID})
		if err != nil {
			log.Error().Err(err).Stack().Msgf("failed to upsert user, %v", err)
			return nil
		}
		if err := l.populateChatState(ctx, upd); err != nil {
			log.Error().Err(err).Stack().Msgf("failed to populateChatState")
			return nil
		}

		res, err := l.Bots.OnMessage(ctx, *upd)
		if err != nil {
			//TODO: сделать канал, в который будем псиать ошибки, а в другом месте вычитывать из него и писать в логи
			return nil
		}
		return res
	})

	return err
}

func (l *AliceListener) populateChatState(ctx context.Context, upd *api.Update) error {
	cs, err := l.ChatStateService.FindByUserId(ctx, upd.Request.Session.UserID)
	if err != nil {
		return errors.Wrapf(err, "failed to find ChatState by id %q", err)
	}
	upd.ChatState = cs
	return nil
}
