package skill

import (
	"context"
	"fmt"
	"github.com/azzzak/alice"
	"github.com/rs/zerolog/log"
	"runtime/debug"
	"sync"
	"tourist-alice-skill/internal/api"
)

//actions
const (
	selectedCity     api.Action = "selected_city"
	wantSelectedCity api.Action = "want_selected_city"
)

// Interface is a skill reactive spec. response will be sent if "send" result is true
type Interface interface {
	OnMessage(ctx context.Context, update api.Update) (f *alice.Response, err error)
	HasReact(update api.Update) bool
}

// MultiSkill combines many bots to one virtual
type MultiSkill []Interface

func (b MultiSkill) OnMessage(ctx context.Context, update api.Update) (ha *alice.Response, err error) {

	resps := make(chan *alice.Response)

	var wg sync.WaitGroup
	for _, bot := range b {
		wg.Add(1)
		fmt.Printf("%v", bot)
		bot := bot
		go func(ctx context.Context, wg *sync.WaitGroup) {
			defer wg.Done()
			defer handlePanic(bot)
			if bot.HasReact(update) {
				resp, _ := bot.OnMessage(ctx, update)
				resps <- resp
			}
		}(ctx, &wg)
	}

	go func() {
		wg.Wait()
		close(resps)
	}()

	var handler *alice.Response
	//TODO: что то сделать с возаратом нескольких респонсов
	for r := range resps {
		log.Debug().Msgf("collect %v", r)
		handler = r
	}

	return handler, nil
}
func handlePanic(bot Interface) {
	if err := recover(); err != nil {
		switch e := err.(type) {
		case error:
			log.Error().Err(e).Stack().Msgf("panic! skill: %T, stack: %s", bot, string(debug.Stack()))
		default:
			log.Error().Stack().Msgf("panic! skill: %t, err: %v, stack: %s", bot, err, string(debug.Stack()))
		}
	}
}

func (b MultiSkill) HasReact(u api.Update) bool {
	var hasReact bool
	for _, bot := range b {
		hasReact = hasReact && bot.HasReact(u)
	}
	fmt.Printf("hasReact MultiSkill %v", hasReact)
	return hasReact
}
