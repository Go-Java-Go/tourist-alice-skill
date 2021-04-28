package skill

import (
	"context"
	"github.com/gookit/i18n"
	"tourist-alice-skill/internal/api"
)

type UserService interface {
	FindById(ctx context.Context, id int) (*api.User, error)
}

type Config struct {
	SkillName string
}

// I18n define text by user lang
func I18n(u *api.User, text string, args ...interface{}) string {
	return i18n.Tr(api.DefineLang(u), text, args...)
}
