package skill

import (
	"context"
	"github.com/azzzak/alice"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tourist-alice-skill/internal/api"
)

type ChatStateService interface {
	Save(ctx context.Context, u *api.ChatState) error
	DeleteById(ctx context.Context, id primitive.ObjectID) error
}

type StartScreen struct {
	css ChatStateService
}

func NewStartScreen(css ChatStateService) *StartScreen {
	return &StartScreen{
		css: css,
	}
}

func (s StartScreen) HasReact(u api.Update) bool {
	return u.Request.IsNewSession() || u.Request.Command() == "Привет"
}

func (s *StartScreen) OnMessage(ctx context.Context, u api.Update) (d *alice.Response, err error) {
	text := "Приветствую тебя в навыке\n"
	text += "В каком городе хочешь посмотреть достопримечательности?"

	cs := &api.ChatState{UserId: u.User.ID, Action: wantSelectedCity}
	err = s.css.Save(ctx, cs)
	if err != nil {
		return nil, err
	}
	u.Response.Text(text)
	u.Response.Buttons(
		alice.NewButton("Москва", "", false),
		alice.NewButton("Казань", "", false))
	return u.Response, nil
}
