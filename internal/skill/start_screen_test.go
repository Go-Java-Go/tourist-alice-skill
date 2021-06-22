// +build test

package skill

import (
	"context"
	"fmt"
	"github.com/azzzak/alice"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"tourist-alice-skill/internal/api"
	"tourist-alice-skill/test/mocks"
)

type StartScreenSuite struct {
	css *mocks.ChatStateService
	ss  *StartScreen
	suite.Suite
}

func TestStartScreenSuiteUp(t *testing.T) {
	suite.Run(t, new(StartScreenSuite))
}

func (s *StartScreenSuite) SetupSuite() {
	s.css = new(mocks.ChatStateService)
	s.ss = NewStartScreen(s.css)
}

func (s *StartScreenSuite) TestStartScreen_HasReact() {

	type args struct {
		u api.Update
	}
	tests := []struct {
		name string
		args args
		res  bool
	}{
		{
			name: "Test when bot react",
			args: args{
				api.Update{
					Request: &alice.Request{
						Request: struct {
							Command           string `json:"command"`
							OriginalUtterance string `json:"original_utterance"`
							Type              string `json:"type"`
							Markup            struct {
								DangerousContext *bool `json:"dangerous_context,omitempty"`
							} `json:"markup,omitempty"`
							Payload interface{} `json:"payload,omitempty"`
							NLU     struct {
								Tokens   []string       `json:"tokens"`
								Entities []alice.Entity `json:"entities,omitempty"`
							} `json:"nlu"`
						}{"Привет",
							"",
							"",
							struct {
								DangerousContext *bool `json:"dangerous_context,omitempty"`
							}{DangerousContext: nil},
							nil,
							struct {
								Tokens   []string       `json:"tokens"`
								Entities []alice.Entity `json:"entities,omitempty"`
							}{Tokens: nil, Entities: nil}},
						Session: struct {
							New       bool   `json:"new"`
							MessageID int    `json:"message_id"`
							SessionID string `json:"session_id"`
							SkillID   string `json:"skill_id"`
							UserID    string `json:"user_id"`
						}{
							true,
							0,
							"",
							"",
							"",
						},
					},
				}},
			res: true,
		},
		{
			name: "Test when bot react with correct command",
			args: args{
				api.Update{
					Request: &alice.Request{
						Request: struct {
							Command           string `json:"command"`
							OriginalUtterance string `json:"original_utterance"`
							Type              string `json:"type"`
							Markup            struct {
								DangerousContext *bool `json:"dangerous_context,omitempty"`
							} `json:"markup,omitempty"`
							Payload interface{} `json:"payload,omitempty"`
							NLU     struct {
								Tokens   []string       `json:"tokens"`
								Entities []alice.Entity `json:"entities,omitempty"`
							} `json:"nlu"`
						}{"Привет",
							"",
							"",
							struct {
								DangerousContext *bool `json:"dangerous_context,omitempty"`
							}{DangerousContext: nil},
							nil,
							struct {
								Tokens   []string       `json:"tokens"`
								Entities []alice.Entity `json:"entities,omitempty"`
							}{Tokens: nil, Entities: nil}},
						Session: struct {
							New       bool   `json:"new"`
							MessageID int    `json:"message_id"`
							SessionID string `json:"session_id"`
							SkillID   string `json:"skill_id"`
							UserID    string `json:"user_id"`
						}{
							false,
							0,
							"",
							"",
							"",
						},
					},
				}},
			res: true,
		},
		{
			name: "Test when bot react with new session",
			args: args{
				api.Update{
					Request: &alice.Request{
						Request: struct {
							Command           string `json:"command"`
							OriginalUtterance string `json:"original_utterance"`
							Type              string `json:"type"`
							Markup            struct {
								DangerousContext *bool `json:"dangerous_context,omitempty"`
							} `json:"markup,omitempty"`
							Payload interface{} `json:"payload,omitempty"`
							NLU     struct {
								Tokens   []string       `json:"tokens"`
								Entities []alice.Entity `json:"entities,omitempty"`
							} `json:"nlu"`
						}{"Test command",
							"",
							"",
							struct {
								DangerousContext *bool `json:"dangerous_context,omitempty"`
							}{DangerousContext: nil},
							nil,
							struct {
								Tokens   []string       `json:"tokens"`
								Entities []alice.Entity `json:"entities,omitempty"`
							}{Tokens: nil, Entities: nil}},
						Session: struct {
							New       bool   `json:"new"`
							MessageID int    `json:"message_id"`
							SessionID string `json:"session_id"`
							SkillID   string `json:"skill_id"`
							UserID    string `json:"user_id"`
						}{
							true,
							0,
							"",
							"",
							"",
						},
					},
				}},
			res: true,
		},
		{
			name: "Test when bot not react",
			args: args{
				api.Update{
					Request: &alice.Request{
						Request: struct {
							Command           string `json:"command"`
							OriginalUtterance string `json:"original_utterance"`
							Type              string `json:"type"`
							Markup            struct {
								DangerousContext *bool `json:"dangerous_context,omitempty"`
							} `json:"markup,omitempty"`
							Payload interface{} `json:"payload,omitempty"`
							NLU     struct {
								Tokens   []string       `json:"tokens"`
								Entities []alice.Entity `json:"entities,omitempty"`
							} `json:"nlu"`
						}{"Test command",
							"",
							"",
							struct {
								DangerousContext *bool `json:"dangerous_context,omitempty"`
							}{DangerousContext: nil},
							nil,
							struct {
								Tokens   []string       `json:"tokens"`
								Entities []alice.Entity `json:"entities,omitempty"`
							}{Tokens: nil, Entities: nil}},
						Session: struct {
							New       bool   `json:"new"`
							MessageID int    `json:"message_id"`
							SessionID string `json:"session_id"`
							SkillID   string `json:"skill_id"`
							UserID    string `json:"user_id"`
						}{
							false,
							0,
							"",
							"",
							"",
						},
					},
				}},
			res: false,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.Equal(tt.res, s.ss.HasReact(tt.args.u))
		})
	}
}

func (s *StartScreenSuite) TestStartScreen_OnMessage_WithoutError() {

	//given
	upd := api.Update{
		User: &api.User{
			ID:           "testId",
			UserLang:     "",
			SelectedLang: "",
		},
		Response: &alice.Response{Response: &struct {
			Text       string         `json:"text"`
			TTS        string         `json:"tts,omitempty"`
			Card       *alice.Card    `json:"card,omitempty"`
			Buttons    []alice.Button `json:"buttons,omitempty"`
			EndSession bool           `json:"end_session"`
		}{Text: "", TTS: "", Card: nil, Buttons: make([]alice.Button, 0, 0), EndSession: false}},
	}

	cs := &api.ChatState{UserId: "testId", Action: wantSelectedCity}

	s.css.On("Save", mock.Anything, cs).Return(nil)

	//when
	r, err := s.ss.OnMessage(context.TODO(), upd)

	//then
	s.css.AssertCalled(s.T(), "Save", mock.Anything, cs)
	s.css.AssertNumberOfCalls(s.T(), "Save", 1)
	s.NoError(err)
	s.Equal(r.Response.Text, "Приветствую тебя в навыке\nВ каком городе хочешь посмотреть достопримечательности?")
	s.Len(r.Response.Buttons, 2)
}

func (s *StartScreenSuite) TestStartScreen_OnMessage_WithError() {
	//given
	upd := api.Update{
		User: &api.User{
			ID:           "testId",
			UserLang:     "",
			SelectedLang: "",
		},
		Response: &alice.Response{Response: &struct {
			Text       string         `json:"text"`
			TTS        string         `json:"tts,omitempty"`
			Card       *alice.Card    `json:"card,omitempty"`
			Buttons    []alice.Button `json:"buttons,omitempty"`
			EndSession bool           `json:"end_session"`
		}{Text: "", TTS: "", Card: nil, Buttons: make([]alice.Button, 0, 0), EndSession: false}},
	}

	testErr := fmt.Errorf("test error")
	s.css.On("Save", mock.Anything, mock.Anything).Return(testErr)

	//when
	_, err := s.ss.OnMessage(context.TODO(), upd)

	//then
	s.ErrorIs(err, testErr)
	s.css.AssertNumberOfCalls(s.T(), "Save", 1)

}
