package youtubeClient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
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

func (y *YouTube) GetChatCycle(ch chan struct{}) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		err := fmt.Errorf("PANIC:'%v'\nRecovered in: %s", r, infrastruct.IdentifyPanic())
		logger.LogError(err)
		ch <- struct{}{}
	}()
	for {
		time.Sleep(y.timeSleep)
		if err := y.GetComments(); err != nil {
			ch <- struct{}{}
			return
		}
	}
}

func (y *YouTube) GetComments() error {
	if y.service == nil {
		ctx := context.Background()

		b, err := ioutil.ReadFile(y.credentionalYouTubeApi)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with ioutil.ReadFile in GetComments"))
			return infrastruct.ErrorInternalServerError
		}
		config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with google.ConfigFromJSON in GetComments"))
			return infrastruct.ErrorInternalServerError
		}
		f, err := os.Open(filepath.Join(y.pathForClientSecret, url.QueryEscape("youtube-go-quickstart.json")))
		if err != nil {
			logger.LogError(errors.Wrap(err, fmt.Sprintf("err with google.ConfigFromJSON (PWD = %s) in GetComments", y.pathForClientSecret)))
			return infrastruct.ErrorInternalServerError
		}

		t := &oauth2.Token{}
		err = json.NewDecoder(f).Decode(t)
		if err != nil {
			logger.LogError(errors.Wrap(err, "NewDecoder"))
			return infrastruct.ErrorInternalServerError
		}
		defer f.Close()
		client := config.Client(ctx, t)
		y.service, err = youtube.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with gyoutube.NewService in GetComments"))
			return infrastruct.ErrorInternalServerError
		}
	}

	if y.channelID == "" {
		res, err := y.service.LiveBroadcasts.List([]string{"id,snippet"}).BroadcastStatus("active").BroadcastType("all").Do()
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with service.LiveBroadcasts.List in GetComments"))
			return infrastruct.ErrorInternalServerError
		}

		if len(res.Items) == 0 {
			a := make([]types.MessageForChat, 1)
			a[0].Message = "Трансляция еще не началась"
			a[0].Name = "Башкирский Диктант"
			timeNow := time.Now()
			a[0].Time = timeNow.Format(time.RFC3339)

			y.Cache = &types.Chat{Messages: a}
			y.timeSleep = time.Minute * 5

			return nil
		}
		y.channelID = res.Items[0].Snippet.LiveChatId
	}
	res2, err := y.service.LiveChatMessages.List(y.channelID, []string{"id,snippet,authorDetails"}).Do()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with service.LiveChatMessages.List in GetComments"))
		y.channelID = ""
		return infrastruct.ErrorInternalServerError
	}
	y.timeSleep = time.Second * 3

	if len(res2.Items) == 0 {
		a := make([]types.MessageForChat, 1)
		a[0].Message = "Сообщений в чате еще нету"
		a[0].Name = "Башкирский Диктант"
		timeNow := time.Now()
		a[0].Time = timeNow.Format(time.RFC3339)

		y.Cache = &types.Chat{Messages: a}

		return nil
	}

	j := 0
	if len(res2.Items) > 30 {
		j = 30
	} else {
		j = len(res2.Items)
	}

	chatMessages := make([]types.MessageForChat, j)

	l := len(res2.Items) - 1

	for ; j > 0; j-- {
		chatMessage := types.MessageForChat{
			Message: res2.Items[l].Snippet.TextMessageDetails.MessageText,
			Name:    res2.Items[l].AuthorDetails.DisplayName,
			Time:    res2.Items[l].Snippet.PublishedAt,
		}
		chatMessages[j-1] = chatMessage
		l--
	}

	y.Cache = &types.Chat{Messages: chatMessages}
	return nil
}
