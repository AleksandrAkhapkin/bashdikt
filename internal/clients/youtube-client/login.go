package youtubeClient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/infrastruct"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/logger"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
)

func (y *YouTube) StartAuthYouTube() error {

	b, err := ioutil.ReadFile(y.credentionalYouTubeApi)
	if err != nil {
		logger.LogError(errors.Wrap(err, fmt.Sprintf("err with ReadFile(s.credentionalYouTubeApi) in StartAuthYouTube PWD = %s", y.credentionalYouTubeApi)))
		return infrastruct.ErrorInternalServerError
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with google.ConfigFromJSON in StartAuthYouTube"))
		return infrastruct.ErrorInternalServerError
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	logger.LogInfo(fmt.Sprintf("Ссылка для авторизации на ЮТУБ: %s&%s\n", authURL, "prompt=consent"))

	return nil
}

func (y *YouTube) CallbackYoutube(code string) error {

	b, err := ioutil.ReadFile(y.credentionalYouTubeApi)
	if err != nil {
		logger.LogError(errors.Wrap(err, fmt.Sprintf("err with ReadFile(s.credentionalYouTubeApi) in CallbackYoutube PWD = %s", y.credentionalYouTubeApi)))
		return infrastruct.ErrorInternalServerError
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with google.ConfigFromJSON in CallbackYoutube"))
		return infrastruct.ErrorInternalServerError
	}

	os.Mkdir(y.pathForClientSecret, 0700)
	cacheFile := filepath.Join(y.pathForClientSecret, url.QueryEscape("youtube-go-quickstart.json"))

	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with config.Exchange in CallbackYoutube"))
		return infrastruct.ErrorInternalServerError
	}

	f, err := os.Create(cacheFile)
	if err != nil {
		logger.LogError(errors.Wrap(err, fmt.Sprintf("err with os.Create in CallbackYoutube for %s", cacheFile)))
		return infrastruct.ErrorInternalServerError
	}
	logger.LogInfo(fmt.Sprintf("Saving credential file to: %s\n", cacheFile))

	defer f.Close()
	json.NewEncoder(f).Encode(tok)

	ctx := context.Background()

	file, err := os.Open(cacheFile)
	if err != nil {
		logger.LogError(errors.Wrap(err, fmt.Sprintf("err with google.ConfigFromJSON (PWD = %s) in GetComments", y.pathForClientSecret)))
		return infrastruct.ErrorInternalServerError
	}

	t := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(t)
	defer f.Close()
	client := config.Client(ctx, t)
	y.service, err = youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with gyoutube.NewService in GetComments"))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}
