package handlers

import (
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/infrastruct"
	"net/http"
	"strconv"
)

func (h *Handlers) GetOrgCabinetInfo(w http.ResponseWriter, r *http.Request) {

	var err error
	getInfo := &types.GetInfoForCabinet{}
	getInfo.Role = r.FormValue("type")
	getInfo.Offset, err = strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}
	getInfo.Limit, err = strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	switch getInfo.Role {
	case "student":
		{
			info, err := h.srv.GetStudentForOrgCabinet(getInfo)
			if err != nil {
				apiErrorEncode(w, err)
				return
			}

			apiResponseEncoder(w, info)
			return
		}
	case "teacher":
		{
			info, err := h.srv.GetTeacherForOrgCabinet(getInfo)
			if err != nil {
				apiErrorEncode(w, err)
				return
			}

			apiResponseEncoder(w, info)
			return
		}
	default:
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
	}
}

func (h *Handlers) GetTeacherCabinetInfo(w http.ResponseWriter, r *http.Request) {

	var err error
	getInfo := &types.GetInfoForCabinet{}
	getInfo.Role = r.FormValue("type")
	getInfo.Offset, err = strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}
	getInfo.Limit, err = strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	switch getInfo.Role {
	case "student":
		{
			info, err := h.srv.GetStudentForTeacherCabinet(getInfo)
			if err != nil {
				apiErrorEncode(w, err)
				return
			}

			apiResponseEncoder(w, info)
			return
		}
	case "pinstudent":
		{
			info, err := h.srv.GetPinStudentForTeacherCabinet(getInfo, claims.UserID)
			if err != nil {
				apiErrorEncode(w, err)
				return
			}

			apiResponseEncoder(w, info)
			return
		}
	default:
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
	}
}

func (h *Handlers) GetMyDictation(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	myDictation, err := h.srv.GetMyDictation(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, myDictation)
}

func (h *Handlers) MakeAndSendMyCertificate(w http.ResponseWriter, r *http.Request) {

	lang := r.FormValue("lang") //bash or rus
	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	if err := h.srv.GetMyCertificate(claims.UserID, claims.Role, lang); err != nil {
		apiErrorEncode(w, err)
		return
	}
}
