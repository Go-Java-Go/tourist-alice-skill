package skill

import (
	"context"
	travel_client "github.com/Go-Java-Go/izi-travel-client"
	"github.com/azzzak/alice"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tourist-alice-skill/internal/api"
)

// OperationScreen send /room, after click on the button 'Присоединиться'
type OperationScreen struct {
	css ChatStateService
	tc  *travel_client.Client
}

func NewOperationScreen(css ChatStateService, tc *travel_client.Client) *OperationScreen {
	return &OperationScreen{
		css: css,
		tc:  tc,
	}
}

func (s OperationScreen) HasReact(u api.Update) bool {
	return u.ChatState != nil && u.ChatState.Action == wantSelectedCity
}

func (s *OperationScreen) OnMessage(ctx context.Context, u api.Update) (d *alice.Response, err error) {
	defer func(css ChatStateService, ctx context.Context, id primitive.ObjectID) {
		err := css.DeleteById(ctx, id)
		if err != nil {
			log.Error().Err(err).Msg("")
		}
	}(s.css, ctx, u.ChatState.ID)

	cs := &api.ChatState{UserId: u.User.ID, Action: selectedCity, CallbackData: &api.CallbackData{SelectedCity: u.Request.Command()}}
	err = s.css.Save(ctx, cs)
	if err != nil {
		log.Error().Err(err).Msg("create chat state failed")
		return
	}
	text := "Выбранный город " + u.Request.Command() + "\n"
	text += u.Request.Command()
	text += ` город на юго-западе России, расположенный на берегах Волги и Казанки. 
В столице полуавтономной Республики Татарстан находится древний кремль – крепость, 
известная своими музеями и святыми местами. Башня Сююмбике, синие и золотые купола Благовещенского собора и яркая джума-мечеть 
Кул-Шариф – одни из самых интересных достопримечательностей кремля.
	`
	u.Response.Text(text)
	u.Response.Buttons(alice.NewButton("Расскажи", "", false))
	return u.Response, nil
}
