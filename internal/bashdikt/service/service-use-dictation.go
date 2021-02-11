package service

import (
	"database/sql"
	"fmt"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/service/mail"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/infrastruct"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/logger"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func (s *Service) UploadFile(file *types.UploadFile) error {

	_, err := s.p.GetMyDictation(file.UserID)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with GetMyDictation in UploadFiles"))
			return infrastruct.ErrorInternalServerError
		}
		if err := s.p.AddOfflineDictation(file.UserID); err != nil {
			logger.LogError(errors.Wrap(err, "err with AddOfflineDictation in UploadFiles"))
			return infrastruct.ErrorInternalServerError
		}
		////прикрепляем ученика
		//находим последнего прикрепленного
		lastPinTeacher, err := s.p.GetLastPinTeacher()
		if err != nil {
			if err != sql.ErrNoRows {
				logger.LogError(errors.Wrap(err, "err with GetLastPin in UploadFile"))
				return infrastruct.ErrorInternalServerError
			}
			//если ранее небыло прикрепленных - находим любого учителя
			randTeacher, err := s.p.GetRandTeacher()
			if err != nil {
				if err != sql.ErrNoRows {
					logger.LogError(errors.Wrap(err, "err with GetRandTeacher in UploadFile"))
					return infrastruct.ErrorInternalServerError
				}
				logger.LogError(errors.Wrap(err, "err with GetRandTeacher in UploadFile (НЕВОЗМОЖНО ПРИКРЕПИИТЬ УЧЕНИКА - НЕ ЗАРЕГИСТРИРОВАНО НИ ОДНОГО УЧИТЕЛЯ!!!!!"))
				return infrastruct.ErrorInternalServerError
			}
			//прикрепляем к рандомному
			if err = s.p.PinnedStudent(file.UserID, randTeacher, types.FormatDicFile); err != nil {
				logger.LogError(errors.Wrap(err, "err with PinnedStudent(randTeacher) in UploadFile"))
				return infrastruct.ErrorInternalServerError
			}
		} else {
			//прикрепляем к следующему после последнего
			nextAfterLastTeacherPin, err := s.p.NextAfterLastTeacherPin(lastPinTeacher)
			if err = s.p.PinnedStudent(file.UserID, nextAfterLastTeacherPin, types.FormatDicFile); err != nil {
				logger.LogError(errors.Wrap(err, "err with PinnedStudent (lastTeacherPin) in UploadFile"))
				return infrastruct.ErrorInternalServerError
			}
		}
	}

	if err := os.Mkdir(filepath.Join(s.pathForDictation, fmt.Sprintf("%d", file.UserID)), 0777); err != nil {
		if !os.IsExist(err) {
			logger.LogError(errors.Wrap(err, "err with os.Mkdir in UploadFiles"))
			return infrastruct.ErrorInternalServerError
		}
	}

	dst, err := os.Create(filepath.Join(s.pathForDictation, fmt.Sprintf("%d/%s", file.UserID, file.Head)))
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with UploadFiles in create file"))
		return infrastruct.ErrorInternalServerError
	}

	_, err = io.Copy(dst, file.Body)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with UploadFiles in io.Copy"))
		return infrastruct.ErrorInternalServerError
	}

	user, err := s.p.GetUserByID(file.UserID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByID in WriteDictation"))
		return infrastruct.ErrorInternalServerError
	}

	//отправка емейла о приеме работы
	err = s.email.SendMail(
		mail.DictantTmplt,
		mail.BodyYouWriteDictantText,
		mail.BodyYouWriteDictantHTML,
		user.Email)
	if err != nil {
		logger.LogError(errors.Wrap(err, fmt.Sprintf("err with sendEmail in WriteDictation FOR USER_ID = %d", user.ID)))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}

func (s *Service) WriteDictation(dictation *types.Dictation) error {

	_, err := s.p.GetMyDictation(dictation.UserID)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with GetMyDictation in UploadFiles"))
			return infrastruct.ErrorInternalServerError
		}
		if err := s.p.AddOnlineDictation(dictation); err != nil {
			logger.LogError(errors.Wrap(err, "err with AddOnlineDictation in WriteDictation"))
			return infrastruct.ErrorInternalServerError
		}
		////прикрепляем ученика
		//находим последнего прикрепленного
		lastPinTeacher, err := s.p.GetLastPinTeacher()
		if err != nil {
			if err != sql.ErrNoRows {
				logger.LogError(errors.Wrap(err, "err with GetLastPin in WriteDictation"))
				return infrastruct.ErrorInternalServerError
			}
			//если ранее небыло прикрепленных - находим любого учителя
			randTeacher, err := s.p.GetRandTeacher()
			if err != nil {
				if err != sql.ErrNoRows {
					logger.LogError(errors.Wrap(err, "err with GetRandTeacher in WriteDictation"))
					return infrastruct.ErrorInternalServerError
				}
				logger.LogError(errors.Wrap(err, "err with GetRandTeacher in WriteDictation (НЕВОЗМОЖНО ПРИКРЕПИИТЬ УЧЕНИКА - НЕ ЗАРЕГИСТРИРОВАНО НИ ОДНОГО УЧИТЕЛЯ!!!!!"))
				return infrastruct.ErrorInternalServerError
			}
			//прикрепляем к рандомному
			if err = s.p.PinnedStudent(dictation.UserID, randTeacher, types.FormatDicOnline); err != nil {
				logger.LogError(errors.Wrap(err, "err with PinnedStudent(randTeacher) in WriteDictation"))
				return infrastruct.ErrorInternalServerError
			}
		} else {
			//прикрепляем к следующему после последнего
			nextAfterLastTeacherPin, err := s.p.NextAfterLastTeacherPin(lastPinTeacher)
			if err = s.p.PinnedStudent(dictation.UserID, nextAfterLastTeacherPin, types.FormatDicOnline); err != nil {
				logger.LogError(errors.Wrap(err, "err with PinnedStudent (lastTeacherPin) in WriteDictation"))
				return infrastruct.ErrorInternalServerError
			}
		}
	}

	user, err := s.p.GetUserByID(dictation.UserID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByID in WriteDictation"))
		return infrastruct.ErrorInternalServerError
	}

	//отправка емейла о приеме работы
	err = s.email.SendMail(
		mail.DictantTmplt,
		mail.BodyYouWriteDictantText,
		mail.BodyYouWriteDictantHTML,
		user.Email)
	if err != nil {
		logger.LogError(errors.Wrap(err, fmt.Sprintf("err with sendEmail  in WriteDictation FOR USER_ID = %d", dictation.UserID)))
		return infrastruct.ErrorInternalServerError
	}
	return nil
}

func (s *Service) GetDictationsNameByStudentID(studentID int) (*types.NameFiles, error) {

	names := &types.NameFiles{}

	//смотрим был ли написан диктант на сайте
	online, err := s.p.CheckOnlineWrite(studentID)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with CheckOnlineWrite in GetDictationsNameByStudentID"))
			return nil, infrastruct.ErrorInternalServerError
		}
		return names, infrastruct.ErrorNotHaveDictant
	}

	if online {
		names.Names = make([]string, 1)
		names.Names[0] = "Online"
		return names, nil
	}

	info, err := ioutil.ReadDir(filepath.Join(s.pathForDictation, fmt.Sprintf("%d", studentID)))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, infrastruct.ErrorBadRequest
		}
		logger.LogError(errors.Wrap(err, "err with ioutil.ReadDir in GetDictationsNameByStudentID"))
		return nil, infrastruct.ErrorInternalServerError
	}

	namesStr := make([]string, 0)
	for i, _ := range info {
		namesStr = append(namesStr, info[i].Name())
	}
	names.Names = namesStr

	return names, nil
}

func (s *Service) GetURLFile(studentID, fileName string) string {
	return filepath.Join(s.pathForDictation, fmt.Sprintf("%s/%s", studentID, fileName))
}

func (s *Service) GetOnlineDictant(studentID int) (*types.Dictation, error) {

	dictation, err := s.p.GetOnlineDictationById(studentID)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with GetOnlineDictationById in GetOnlineDictant"))
			return nil, infrastruct.ErrorInternalServerError
		}
		return nil, infrastruct.ErrorBadRequest
	}

	return dictation, nil
}

func (s *Service) ReplyDictation(replyDictation *types.ReplyDictation) error {

	if err := s.p.ChangeRatingAndStatus(replyDictation); err != nil {
		logger.LogError(errors.Wrap(err, "err with ChangeRatingAndStatus in ReplyDictation"))
		return infrastruct.ErrorInternalServerError
	}

	for i := range replyDictation.Markers {
		if err := s.p.AddMarkersForText(&replyDictation.Markers[i], replyDictation.UserID, replyDictation.TeacherID); err != nil {
			logger.LogError(errors.Wrap(err, "err with AddMarkersForText in ReplyDictation"))
			return infrastruct.ErrorInternalServerError
		}
	}

	user, err := s.p.GetUserByID(replyDictation.UserID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByID in ReplyDictation"))
		return infrastruct.ErrorInternalServerError
	}
	err = s.email.SendMail(
		mail.DictantTmplt,
		mail.BodyYourDictationIsCheckedText,
		mail.BodyYourDictationIsCheckedHTML,
		user.Email)
	if err != nil {
		logger.LogError(errors.Wrap(err, fmt.Sprintf("err with sendEmail in ReplyDictation FOR USER_ID = %d", replyDictation.UserID)))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}

func (s *Service) DeleteDictationFile(userID int, filename string) error {

	if err := os.Remove(filepath.Join(s.pathForDictation, strconv.Itoa(userID), filename)); err != nil {
		return infrastruct.ErrorBadRequest
	}

	files, err := ioutil.ReadDir(filepath.Join(s.pathForDictation, fmt.Sprintf("%d", userID)))
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with UploadFiles in ioutil.ReadDir"))
	}

	for i := range files {
		oldName := files[i].Name()
		err := os.Rename(filepath.Join(s.pathForDictation, strconv.Itoa(userID), oldName),
			filepath.Join(s.pathForDictation, strconv.Itoa(userID), fmt.Sprintf("file_%d", i+1)))
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with Rename in DeleteDictationFile"))
			return infrastruct.ErrorInternalServerError
		}
	}

	return nil
}
