package handlers

import (
	"encoding/json"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"net/http"
)

func (h *Handlers) AddWhiteEmail(w http.ResponseWriter, r *http.Request) {

	addEmail := types.AddEmail{}
	if err := json.NewDecoder(r.Body).Decode(&addEmail); err != nil {
		apiErrorEncode(w, err)
		return
	}

	err := h.srv.AddWhiteEmail(addEmail.Email)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}
