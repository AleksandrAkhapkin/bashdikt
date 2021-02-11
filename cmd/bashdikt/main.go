package main

import (
	"flag"
	"fmt"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/server"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/server/handlers"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/service"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/service/mail"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types/config"
	"github.com/AleksandrAkhapkin/bashdikt/internal/clients/postgres"
	youtubeClient "github.com/AleksandrAkhapkin/bashdikt/internal/clients/youtube-client"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/infrastruct"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/logger"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
)

func main() {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		err := fmt.Errorf("PANIC:'%v'\nRecovered in: %s", r, infrastruct.IdentifyPanic())
		logger.LogError(err)
	}()

	configPath := new(string)
	debug := new(bool)
	flag.StringVar(configPath, "config-path", "config/config-local.yaml", "path to yaml config")
	flag.BoolVar(debug, "debug-mod", true, "debug usage")
	flag.Parse()

	cnfFile, err := os.Open(*configPath)
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "err with os.Open"))
	}

	logger.Debug = *debug
	logger.CheckDebug()
	cnf := config.Config{}
	if err := yaml.NewDecoder(cnfFile).Decode(&cnf); err != nil {
		logger.LogFatal(errors.Wrap(err, "err with Decode config"))
	}

	err = logger.NewLogger(cnf.Telegram)
	if err != nil {
		logger.LogFatal(err)
	}

	pq, err := postgres.NewPostgres(cnf.PostgresDsn)
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "err with NewPostgres"))
	}
	email, err := mail.NewMail(cnf.Email, cnf.ServerPort)
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "err with NewService"))
	}

	srv, err := service.NewService(pq, email, &cnf)
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "err with NewService"))
	}

	ytb, err := youtubeClient.NewYouTube(&cnf)
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "with NewYouTube"))
	}

	handls, err := handlers.NewHandlers(srv, ytb, &cnf)
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "with NewHandlers"))
	}

	//go srv.PingYouTubeComment() //проверка работоспособности ютуба

	ch := make(chan struct{})
	go ytb.GetChatCycle(ch) //запускаем параллельный процесс вытягивания комментов
	go func() {
		for {
			// в паралелль запускаем бесконечный цикл вычитывания канала
			_ = <-ch
			go ytb.GetChatCycle(ch)
			//если мы вычитали канал - значит была ошибка и надо запустить процесс вытягивания комментов заново
		}
	}()

	server.StartServer(handls, cnf.ServerPort)
}
