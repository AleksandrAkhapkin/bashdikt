package handlers

import (
	"errors"
	"fmt"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/infrastruct"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/logger"
	"net/http"
	"net/http/httputil"
	"strings"
)

func (h *Handlers) RecoveryPanic(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}

			err := errors.New(fmt.Sprintf("PANIC:'%v'\nRecovered in: %s", r, infrastruct.IdentifyPanic()))
			logger.LogError(err)
			apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		}()

		handler.ServeHTTP(w, r)
	})
}

func (h *Handlers) CheckRoleStudent(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
		if err != nil {
			apiErrorEncode(w, err)
			return
		}

		if claims.Role != types.RoleStudent {
			apiErrorEncode(w, infrastruct.ErrorPermissionDenied)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (h *Handlers) CheckRoleTeacher(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
		if err != nil {
			apiErrorEncode(w, err)
			return
		}

		if claims.Role != types.RoleTeacher {
			apiErrorEncode(w, infrastruct.ErrorPermissionDenied)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (h *Handlers) CheckRoleOrganizer(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
		if err != nil {
			apiErrorEncode(w, err)
			return
		}

		if claims.Role != types.RoleOrganizer {
			apiErrorEncode(w, infrastruct.ErrorPermissionDenied)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (h *Handlers) CheckRoleAdmin(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
		if err != nil {
			apiErrorEncode(w, err)
			return
		}

		if claims.Role != types.RoleAdmin {
			apiErrorEncode(w, infrastruct.ErrorPermissionDenied)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (h *Handlers) CheckUserInDBUsers(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
		if err != nil {
			apiErrorEncode(w, err)
			return
		}
		if err = h.srv.CheckUserInDBUsers(claims.UserID); err != nil {
			apiErrorEncode(w, err)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (h *Handlers) RequestLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		func() {
			req, err := httputil.DumpRequest(r, !strings.Contains(r.RequestURI, "/dictation/upload"))
			if err != nil {
				logger.LogError(err)
				return
			}

			ip := r.Header.Get("X-Real-IP")

			if err = h.srv.SaveRequestLog(req, r.RequestURI, ip); err != nil {
				logger.LogError(err)
				return
			}
		}()

		handler.ServeHTTP(w, r)
	})
}
