package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/infrastruct"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"strconv"
)

func (h *Handlers) UploadFiles(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	fileType := types.UploadFile{UserID: claims.UserID}
	err = r.ParseMultipartForm(209715200)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with ParseMultipartForm in UploadFiles"))
		apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		return
	}

	files := r.MultipartForm.File["file"]
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with r.MultipartForm.File in UploadFiles"))
			apiErrorEncode(w, infrastruct.ErrorInternalServerError)
			return
		}
		head := file.Filename
		fileType.Body = f
		fileType.Head = head
		if err = h.srv.UploadFile(&fileType); err != nil {
			apiErrorEncode(w, infrastruct.ErrorInternalServerError)
			return
		}
	}

}

func (h *Handlers) GetDictationsNameByStudentID(w http.ResponseWriter, r *http.Request) {

	studentID, err := strconv.Atoi(mux.Vars(r)["idStudent"])
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	name, err := h.srv.GetDictationsNameByStudentID(studentID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, name)
}

func (h *Handlers) WriteDictation(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	dictation := types.Dictation{}
	if err := json.NewDecoder(r.Body).Decode(&dictation); err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}
	dictation.UserID = claims.UserID

	if err := h.srv.WriteDictation(&dictation); err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) GetDictationFile(w http.ResponseWriter, r *http.Request) {

	studentID := mux.Vars(r)["idStudent"]
	fileName := r.FormValue("name")

	if fileName != "Online" {
		fileURL, err := os.Open(h.srv.GetURLFile(studentID, fileName))
		if err != nil {
			if os.IsNotExist(err) {
				apiErrorEncode(w, infrastruct.ErrorBadRequest)
				return
			}
			logger.LogError(errors.Wrap(err, "err with os.Open in GetDictationsFile "))
			apiErrorEncode(w, infrastruct.ErrorInternalServerError)
			return
		}

		fileinfo, err := fileURL.Stat()
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with fileIRL.Stat"))
			apiErrorEncode(w, infrastruct.ErrorBadRequest)
			return
		}

		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileinfo.Size()))
		if _, err = io.Copy(w, fileURL); err != nil {
			logger.LogError(errors.Wrap(err, "err with io.Copy in GetDictationFile"))
			apiErrorEncode(w, infrastruct.ErrorInternalServerError)
			return
		}
		return
	}

	stdntID, err := strconv.Atoi(studentID)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	text, err := h.srv.GetOnlineDictant(stdntID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, text)
}

func (h *Handlers) GetDictationFileForStudent(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	studentID := strconv.Itoa(claims.UserID)
	fileName := r.FormValue("name")

	if fileName != "Online" {
		fileURL, err := os.Open(h.srv.GetURLFile(studentID, fileName))
		if err != nil {
			if os.IsNotExist(err) {
				apiErrorEncode(w, infrastruct.ErrorBadRequest)
				return
			}
			logger.LogError(errors.Wrap(err, "err with os.Open in GetDictationsFile "))
			apiErrorEncode(w, infrastruct.ErrorInternalServerError)
			return
		}

		if _, err = io.Copy(w, fileURL); err != nil {
			logger.LogError(errors.Wrap(err, "err with io.Copy in GetDictationFile"))
			apiErrorEncode(w, infrastruct.ErrorInternalServerError)
			return
		}
	}

	text, err := h.srv.GetOnlineDictant(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, text)
}

func (h *Handlers) DeleteFileForStudent(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	fileName := r.FormValue("name")

	if err := h.srv.DeleteDictationFile(claims.UserID, fileName); err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) ReplyDictation(w http.ResponseWriter, r *http.Request) {

	stdntID := mux.Vars(r)["idStudent"]
	studentID, err := strconv.Atoi(stdntID)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	replyDictation := types.ReplyDictation{}
	if err := json.NewDecoder(r.Body).Decode(&replyDictation); err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}
	replyDictation.UserID = studentID
	replyDictation.TeacherID = claims.UserID

	if err := h.srv.ReplyDictation(&replyDictation); err != nil {
		apiErrorEncode(w, err)
		return
	}
}
