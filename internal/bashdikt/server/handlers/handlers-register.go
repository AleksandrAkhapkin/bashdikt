package handlers

import (
	"encoding/json"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/infrastruct"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/logger"
	"github.com/pkg/errors"
	"net/http"
)

func (h *Handlers) RegisterNew(w http.ResponseWriter, r *http.Request) {

	userRole := r.FormValue("role")
	var err error

	switch userRole {
	case "student":
		{
			user := types.StudentRegister{}
			if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
				apiErrorEncode(w, infrastruct.ErrorBadRequest)
				return
			}
			if user.Pass != user.RepeatPass {
				apiErrorEncode(w, infrastruct.ErrorPasswordsDoNotMatch)
				return
			}
			err = h.srv.RegisterStudent(&user)
			if err != nil {
				apiErrorEncode(w, err)
				return
			}
		}
	case "teacher":
		{
			user := types.TeacherRegister{}
			if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
				apiErrorEncode(w, infrastruct.ErrorBadRequest)
				return
			}
			if user.Pass != user.RepeatPass {
				apiErrorEncode(w, infrastruct.ErrorPasswordsDoNotMatch)
				return
			}
			if err := h.srv.RegisterTeacher(&user); err != nil {
				apiErrorEncode(w, err)
				return
			}
		}
	case "organizer":
		{
			user := types.OrganizerRegister{}
			if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
				apiErrorEncode(w, infrastruct.ErrorBadRequest)
				return
			}
			if user.Pass != user.RepeatPass {
				apiErrorEncode(w, infrastruct.ErrorPasswordsDoNotMatch)
				return
			}
			err = h.srv.OrganizerRegister(&user)
			if err != nil {
				apiErrorEncode(w, err)
				return
			}
		}
	default:
		err = infrastruct.ErrorBadRequest
		apiErrorEncode(w, err)
		return
	}
}

//func (h *Handlers) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
//
//	user := types.ConfirmationEmail{}
//	var err error
//
//	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
//		apiErrorEncode(w, infrastruct.ErrorBadRequest)
//		return
//	}
//
//	token, err := h.srv.ConfirmationEmail(&user)
//	if err != nil {
//		apiErrorEncode(w, err)
//		return
//	}
//
//	apiResponseEncoder(w, token)
//}

func (h *Handlers) ConfirmByLink(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	code := r.FormValue("key")
	var err error

	if err = h.srv.AuthByLink(email, code); err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) Auth(w http.ResponseWriter, r *http.Request) {

	userAuth := &types.Auth{}

	if err := json.NewDecoder(r.Body).Decode(&userAuth); err != nil {
		logger.LogError(errors.Wrap(err, "err with json.Decode"))
		err = infrastruct.ErrorInternalServerError
		apiErrorEncode(w, err)
		return
	}
	token, err := h.srv.AuthUsers(userAuth)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, token)
}

func (h *Handlers) RecoverPassword(w http.ResponseWriter, r *http.Request) {

	rec := types.RecoverPass{}
	var err error

	if err = json.NewDecoder(r.Body).Decode(&rec); err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	if rec.Email == "" {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	err = h.srv.RecoverPassword(&rec)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}
