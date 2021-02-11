package handlers

import (
	"encoding/json"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/infrastruct"
	"net/http"
)

func (h *Handlers) Cabinet(w http.ResponseWriter, r *http.Request) {

	var err error
	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	switch claims.Role {
	case types.RoleStudent:
		{
			cabinet, err := h.srv.GetCabinetStudent(claims.UserID)
			if err != nil {
				apiErrorEncode(w, err)
				return
			}
			apiResponseEncoder(w, cabinet)
			return
		}

	case types.RoleTeacher:
		{
			cabinet, err := h.srv.GetCabinetTeacher(claims.UserID)
			if err != nil {
				apiErrorEncode(w, err)
				return
			}
			apiResponseEncoder(w, cabinet)
			return
		}

	case types.RoleOrganizer:
		{
			cabinet, err := h.srv.GetCabinetOrganizer(claims.UserID)
			if err != nil {
				apiErrorEncode(w, err)
				return
			}
			apiResponseEncoder(w, cabinet)
			return
		}
	default:
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
	}
}

func (h *Handlers) PutCabinet(w http.ResponseWriter, r *http.Request) {

	var err error
	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	switch claims.Role {
	case types.RoleStudent:
		{
			user := types.StudentProfileForPUT{StandardProfileForPUT: types.StandardProfileForPUT{ID: claims.UserID, Role: claims.Role}}
			if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
				apiErrorEncode(w, infrastruct.ErrorBadRequest)
				return
			}
			cabinet, err := h.srv.PutCabinetStudent(&user)
			if err != nil {
				apiErrorEncode(w, err)
				return
			}
			apiResponseEncoder(w, cabinet)

		}

	case types.RoleTeacher:
		{
			user := types.TeacherProfileForPUT{StandardProfileForPUT: types.StandardProfileForPUT{ID: claims.UserID, Role: claims.Role}}
			if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
				apiErrorEncode(w, infrastruct.ErrorBadRequest)
				return
			}
			cabinet, err := h.srv.PutCabinetTeacher(&user)
			if err != nil {
				apiErrorEncode(w, err)
				return
			}
			apiResponseEncoder(w, cabinet)

		}

	case types.RoleOrganizer:
		{
			user := types.OrganizerProfileForPUT{StandardProfileForPUT: types.StandardProfileForPUT{ID: claims.UserID, Role: claims.Role}}
			if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
				apiErrorEncode(w, infrastruct.ErrorBadRequest)
				return
			}
			cabinet, err := h.srv.PutCabinetOrganizer(&user)
			if err != nil {
				apiErrorEncode(w, err)
				return
			}
			apiResponseEncoder(w, cabinet)
		}
	default:
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
	}

}
