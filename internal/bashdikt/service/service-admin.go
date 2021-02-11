package service

import (
	"database/sql"
	"fmt"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/service/mail"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/infrastruct"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/logger"
	"github.com/pkg/errors"
	"strings"
)

//добавляет новый емейл в белый лист
func (s *Service) AddWhiteEmail(email string) error {

	email = strings.ReplaceAll(email, " ", "")
	email = strings.ToLower(email)

	if err := s.p.AddWhiteEmail(email); err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.AddWhiteEmail"))
		return infrastruct.ErrorInternalServerError
	}

	if err := s.p.UpdateStatusEmail(email); err != nil {
		logger.LogError(errors.Wrap(err, "err with UpdateStatusEmail in AddWhiteEmail"))
		return infrastruct.ErrorInternalServerError
	}

	userInDB, err := s.p.GetUserByEmail(email)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with GetUserByEmail in AddWhiteEmail"))
			return infrastruct.ErrorInternalServerError
		}
		//если не нашли в базе
		err = s.email.SendMail(
			mail.DictantTmplt,
			fmt.Sprintf(mail.BodyInvateTeacherText, fmt.Sprintf(mail.LinkForRegisterTeacher, email)),
			fmt.Sprintf(mail.BodyInvateTeacherHTML, fmt.Sprintf(mail.LinkForRegisterTeacher, email)),
			email)
		if err != nil {
			logger.LogError(errors.Wrap(err, fmt.Sprintf("err with SendMail in confirmURL FOR email = %s", email)))
			return infrastruct.ErrorInternalServerError
		}
		return nil
	}

	if !userInDB.ConfirmEmail {
		//если нашли в базе
		err = s.email.SendMail(
			mail.DictantTmplt,
			mail.BodyActivateTeacherText,
			mail.BodyActivateTeacherHTML,
			email)
		if err != nil {
			logger.LogError(errors.Wrap(err, fmt.Sprintf("err with SendMail in confirmURL FOR email = %s", email)))
			return infrastruct.ErrorInternalServerError
		}
	}

	return nil
}
