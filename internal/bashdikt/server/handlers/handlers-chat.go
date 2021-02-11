package handlers

import (
	"github.com/AleksandrAkhapkin/bashdikt/pkg/infrastruct"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/logger"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

func (h *Handlers) GetChat(w http.ResponseWriter, r *http.Request) {

	apiResponseEncoder(w, h.ytb.Cache)
}

func (h *Handlers) StartAuthYoutube(w http.ResponseWriter, r *http.Request) {

	err := h.ytb.StartAuthYouTube()
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) CallbackYoutube(w http.ResponseWriter, r *http.Request) {

	code := r.FormValue("code")
	code, err := url.QueryUnescape(code)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with url.QueryUnescape in CallbackYoutube"))
		apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		return
	}

	err = h.ytb.CallbackYoutube(code)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

}
