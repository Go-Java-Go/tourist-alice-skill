package api

import (
	"github.com/azzzak/alice"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/language"
)

type Update struct {
	Request   *alice.Request
	Response  *alice.Response
	User      *User
	ChatState *ChatState
}

// ChatState stores user state
type ChatState struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId       string             `json:"userId" bson:"user_id"`
	Action       Action             `json:"action" bson:"action"`
	CallbackData *CallbackData      `json:"callbackData" bson:"callback_data"`
}

type Action string

type CallbackData struct {
	SelectedCity string `bson:"selected_city,omitempty"`
	Page         int    `bson:"page,omitempty"`
}

// User defines user info of the Message
type User struct {
	ID           string `json:"id" bson:"_id"`
	UserLang     string `json:"userLang" bson:"user_lang"`
	SelectedLang string `json:"selectedLang" bson:"selected_lang"`
}

func DefineLang(u *User) string {
	if u.SelectedLang != "" {
		return u.SelectedLang
	} else {
		if u.UserLang == language.English.String() || u.UserLang == language.Russian.String() {
			return u.UserLang
		} else {
			return language.English.String()
		}
	}
}
