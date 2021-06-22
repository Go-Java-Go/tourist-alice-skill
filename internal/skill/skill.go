package skill

import (
	"context"
	"fmt"
	"github.com/azzzak/alice"
	"github.com/rs/zerolog/log"
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

	childCtx, cancelFunc := context.WithCancel(ctx)

	resps := make(chan *alice.Response)
	errors := make(chan error)

	var wg sync.WaitGroup
	for _, bot := range b {
		wg.Add(1)
		fmt.Printf("%v", bot)
		bot := bot
		go func(ctx context.Context, wg *sync.WaitGroup) {
			defer wg.Done()
			if bot.HasReact(update) {
				resp, er := bot.OnMessage(ctx, update)
				if er != nil {
					log.Error().Err(er).Msg("1")
					errors <- er
				} else {
					resps <- resp
				}

			}
		}(childCtx, &wg)
	}

	go func() {
		wg.Wait()
		close(resps)
		close(errors)
		cancelFunc()
	}()

	var handler *alice.Response
	var eror error

tobreake:
	for {
		select {
		case resp, ok := <-resps:
			if !ok {
				break tobreake
			}
			handler = resp

		case err, ok := <-errors:
			if !ok {
				break tobreake
			}
			eror = err

		default:

		}
	}

	return handler, eror
}

func (b MultiSkill) HasReact(u api.Update) bool {
	var hasReact bool
	for _, bot := range b {
		hasReact = hasReact && bot.HasReact(u)
	}
	fmt.Printf("hasReact MultiSkill %v", hasReact)
	return hasReact
}
