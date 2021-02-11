package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/service"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types/config"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/infrastruct"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/logger"

	youtubeClient "github.com/AleksandrAkhapkin/bashdikt/internal/clients/youtube-client"

	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type Handlers struct {
	srv       *service.Service
	secretKey string
	soc       *config.SocAuth
	ytb       *youtubeClient.YouTube
}

func (h *Handlers) Ping(w http.ResponseWriter, _ *http.Request) {

	_, _ = w.Write([]byte("pong 1102 1216 THIS IN PROD BRANCH"))
}

func (h *Handlers) TimeDictation(w http.ResponseWriter, r *http.Request) {
	td := struct {
		Time string `json:"time"`
		URL  string `json:"url"`
	}{}
	level := r.FormValue("level")

	//todo UPDATE PROD
	switch level {
	case "start":
		td.Time = "2020-12-30T12:00:00"
		td.URL = "https://www.youtube.com/?gl=UA"
	case "advanced":
		td.Time = "2020-12-30T15:00:00"
		td.URL = "https://www.youtube.com/?gl=UA"
	case "dialect":
		td.Time = "2020-12-30T18:00:00"
		td.URL = "https://www.youtube.com/?gl=UA"
	}

	apiResponseEncoder(w, td)
}

func NewHandlers(srv *service.Service, ytb *youtubeClient.YouTube, cnf *config.Config) (*Handlers, error) {
	return &Handlers{
		srv:       srv,
		secretKey: cnf.SecretKeyJWT,
		soc:       cnf.Soc,
		ytb:       ytb,
	}, nil
}

func apiErrorEncode(w http.ResponseWriter, err error) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if customError, ok := err.(*infrastruct.CustomError); ok {
		w.WriteHeader(customError.Code)
	}

	result := struct {
		Err string `json:"error"`
	}{
		Err: err.Error(),
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		logger.LogError(err)
	}
}

func apiResponseEncoder(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.LogError(err)
	}
}

func (h *Handlers) SendFrontError(w http.ResponseWriter, r *http.Request) {

	frontError := types.FrontError{}

	respLogerBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with ReadAll in SendFrontError"))
		return
	}
	if err := json.Unmarshal(respLogerBytes, &frontError); err != nil {
		logger.LogError(errors.Wrap(err, "err with Unmarshal in SendFrontError"))
		return
	}

	if frontError.Device.IsMobile {
		logger.LogError(fmt.Errorf("\n\nRoute: %s\nError:\n  Line: %d\n  Column: %d\n  SourceURL: %s\nBROWSER: \n   Name: %s\n   Version: %s\nDevice:\n   MOBILE\n   Vendor: %s\n   Model: %s\n   OS: %s\n   osVersion: %s\n   ua: %s\n\nJSON: %s",
			frontError.Route, frontError.Error.Line, frontError.Error.Column, frontError.Error.SourceURL, frontError.Device.Vendor, frontError.Device.Model, frontError.Device.OS, frontError.Device.OsVersion, frontError.Device.UA, frontError.Browser.Name, frontError.Browser.Version, string(respLogerBytes)))
		return
	}
	logger.LogError(fmt.Errorf("\n\nRoute: %s\nError:\n  Line: %d\n  Column: %d\n  SourceURL: %s\nDevice:\n   Is browser: %t\n   Browser major version: %s\n   Browser full version: %s\n   Browser name: %s\n   Engine name: %s\n   Engine version: %s\n   Os name: %s\n   Os version: %s\n   User agent: %s\n\nJSON: %s",
		frontError.Route, frontError.Error.Line, frontError.Error.Column, frontError.Error.SourceURL, frontError.Device.IsBrowser, frontError.Device.BrowserMajorVersion, frontError.Device.BrowserFullVersion, frontError.Device.BrowserName, frontError.Device.EngineName, frontError.Device.EngineVersion, frontError.Device.OsName, frontError.Device.OsVersion, frontError.Device.UserAgent, string(respLogerBytes)))
}
