package main

import (
	"context"
	"fmt"
	itclient "github.com/Go-Java-Go/izi-travel-client"
	"github.com/etherlabsio/healthcheck"
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
	bot "tourist-alice-skill/internal/skill"

	"github.com/xlab/closer"
)

func main() {
	defer closer.Close()
	ctx, ctxCl := context.WithCancel(context.Background())
	closer.Bind(ctxCl)

	cfg, err := initConfig()
	if err != nil {
		log.Err(err).Msg("Can not init config")
	}

	if err := initLogger(cfg); err != nil {
		log.Fatal().Err(err).Msg("Can not init logger")
	}
	initI18n(cfg)

	conn, cl, err := initMongoConnection(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Can not init mongo connection")
	}
	closer.Bind(cl)

	go initHealthCheck(conn)

	rand.Seed(int64(time.Now().Nanosecond()))

	app, err := initApp(ctx, conn, cfg)
	if err != nil {
		log.Error().Err(err).Msg("Can not init application")
		return
	}
	if err = app.Do(ctx); err != nil {
		log.Error().Err(err).Msg("telegram listener failed")
		return
	}
}

func initSkillConfig(bots []bot.Interface, us events.UserService, css events.ChatStateService) (*events.AliceListener, error) {
	multiSkill := bot.MultiSkill(bots)
	tgListener := &events.AliceListener{
		Bots:             multiSkill,
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

func initMongoDatabase(cli *mongo.Client, cfg *config) *mongo.Database {
	return cli.Database(cfg.DbName)
}

func initMongoConnection(ctx context.Context, cfg *config) (*mongo.Client, func(), error) {
	client, err := mongo.NewClient(
		options.Client().
			ApplyURI(cfg.DbAddr).
			SetAuth(options.Credential{
				Username: cfg.DbUser,
				Password: cfg.DbPassword,
			}))
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
	return client, func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal().Err(err).Msg("error while connect to mongo")
		}
	}, nil
}

func initIziTravelClient(c *config) (*itclient.Client, error) {
	client, err := itclient.NewClient(itclient.Config{APIKey: "d", Host: c.IziTravelHost})
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

func initHealthCheck(cli *mongo.Client) {
	http.HandleFunc("/health", healthcheck.HandlerFunc(

		healthcheck.WithTimeout(5*time.Second),

		healthcheck.WithChecker("mongo", healthcheck.CheckerFunc(
			func(ctx context.Context) error {
				return cli.Ping(ctx, nil)
			})),

		healthcheck.WithChecker("heartbeat", healthcheck.CheckerFunc(func(ctx context.Context) error {
			log.Debug().Msg("health check called")
			return nil
		})),
	))
	log.Info().Str("startup port", "3000").Msg("Server started")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Error().Err(err).Msg("Can not start server")
	}
}
