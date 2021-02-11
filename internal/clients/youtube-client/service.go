package youtubeClient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types/config"
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
	"time"
)

type YouTube struct {
	service                *youtube.Service
	credentionalYouTubeApi string
	pathForClientSecret    string
	channelID              string
	Cache                  *types.Chat
	timeSleep              time.Duration
}

func NewYouTube(cnf *config.Config) (*YouTube, error) {

	mockChat := &types.Chat{}
	a := make([]types.MessageForChat, 1)
	a[0].Message = "Трансляция еще не началась"
	a[0].Name = "Башкирский Диктант"
	timeNow := time.Now()
	a[0].Time = timeNow.Format(time.RFC3339)
	mockChat.Messages = a

	ctx := context.Background()
	b, err := ioutil.ReadFile(cnf.PathSECRETForYouTubeApi)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with ioutil.ReadFile in GetComments"))
		return nil, infrastruct.ErrorInternalServerError
	}
	conf, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with google.ConfigFromJSON in GetComments"))
		return nil, infrastruct.ErrorInternalServerError
	}
	f, err := os.Open(filepath.Join(cnf.PathForClientSecret, url.QueryEscape("youtube-go-quickstart.json")))
	if err != nil {
		logger.LogError(errors.Wrap(err, fmt.Sprintf("err with google.ConfigFromJSON (PWD = %s) in GetComments", cnf.PathForClientSecret)))
		return nil, infrastruct.ErrorInternalServerError
	}

	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	if err != nil {
		logger.LogError(errors.Wrap(err, "NewDecoder"))
		return nil, infrastruct.ErrorInternalServerError
	}
	defer f.Close()
	client := conf.Client(ctx, t)
	serviceYouTube, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with gyoutube.NewService in GetComments"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return &YouTube{
		service:                serviceYouTube,
		credentionalYouTubeApi: cnf.PathSECRETForYouTubeApi,
		pathForClientSecret:    cnf.PathForClientSecret,
		channelID:              "",
		Cache:                  mockChat,
		timeSleep:              time.Second * 5,
	}, nil
}
