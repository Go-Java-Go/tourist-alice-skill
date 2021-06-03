package main

import (
	"context"
	"fmt"
	it_client "github.com/Go-Java-Go/izi-travel-client"
	"github.com/gookit/i18n"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/text/language"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"tourist-alice-skill/internal/events"
	"tourist-alice-skill/internal/handler"
	bot "tourist-alice-skill/internal/skill"

	"github.com/xlab/closer"
)

func main() {
	defer closer.Close()
	ctx := context.Background()

	cfg, err := initConfig()
	if err != nil {
		log.Err(err).Msg("Can not init config")
	}

	if err := initLogger(cfg); err != nil {
		log.Fatal().Err(err).Msg("Can not init logger")
	}
	initI18n(cfg)

	go func() {
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			log.Error().Err(err).Msg("")
		}
	}()

	rand.Seed(int64(time.Now().Nanosecond()))

	app, cl, err := initApp(ctx, cfg)
	if err != nil {
		log.Error().Err(err).Msg("Can not init application")
		return
	}
	if err = app.Do(ctx); err != nil {
		log.Error().Err(err).Msg("telegram listener failed")
		return
	}
	closer.Bind(cl)
}

func initSkillConfig(bots []bot.Interface, us events.UserService, css events.ChatStateService, eh *handler.ErrorHandler) (*events.AliceListener, error) {
	multiSkill := bot.MultiSkill(bots)
	tgListener := &events.AliceListener{
		Bots:             multiSkill,
		ErrorHandler:     eh,
		UserService:      us,
		ChatStateService: css,
	}
	return tgListener, nil
}

func initLogger(c *config) error {
	log.Debug().Msg("initialize logger")
	logLvl, err := zerolog.ParseLevel(strings.ToLower(c.LogLevel))
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(logLvl)
	switch c.LogFmt {
	case "console":
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	case "json":
	default:
		return fmt.Errorf("unknown output format %service", c.LogFmt)
	}
	return nil
}

func initMongoConnection(ctx context.Context, cfg *config) (*mongo.Database, func(), error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.DbAddr))
	if err != nil {
		return nil, nil, err
	}

	// Create connect
	err = client.Connect(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	return client.Database(cfg.DbName), func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal().Err(err).Msg("error while connect to mongo")
		}
	}, nil
}

func initIziTravelClient(c *config) (*it_client.Client, error) {
	client, err := it_client.NewClient(it_client.Config{APIKey: "d", Host: c.IziTravelHost})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func initI18n(c *config) {
	languages := map[string]string{
		language.English.String(): "English",
		language.Russian.String(): "Русский",
	}
	i18n.Init("conf/lang", c.DefaultLanguage, languages)
}
