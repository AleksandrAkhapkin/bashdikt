package service

import (
	"database/sql"
	"fmt"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/service/mail"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/infrastruct"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/logger"
	"github.com/hashicorp/go-uuid"

	"github.com/pkg/errors"
	"strings"
	"time"
)

func (s *Service) RegisterStudent(user *types.StudentRegister) error {

	user.Role = types.RoleStudent
	user.Email = strings.ToLower(user.Email)
	delSpaceRegister(&user.Register)
	var err error

	if err = s.checkDoubleUserEmail(user.Email); err != nil {
		return err
	}
	if err = s.confirmURL(user.Email); err != nil {
		return err
	}

	user.ID, err = s.p.CreateUser(&user.Register)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateUser in RegisterStudent"))
		return infrastruct.ErrorInternalServerError
	}

	if err := s.p.CreateStudentInfo(user); err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateStudentInfo in RegisterStudent "))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}

func (s *Service) RegisterTeacher(user *types.TeacherRegister) error {

	user.Role = types.RoleTeacher
	user.Email = strings.ToLower(user.Email)
	delSpaceRegister(&user.Register)
	var err error

	if err = s.checkDoubleUserEmail(user.Email); err != nil {
		return err
	}

	user.ID, err = s.p.CreateUser(&user.Register)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateUser in RegisterTeacher"))
		return infrastruct.ErrorInternalServerError
	}

	if err := s.p.CreateTeacherInfo(user); err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateStudentInfo in RegisterTeacher "))
		return infrastruct.ErrorInternalServerError
	}

	whiteEmails, err := s.p.GetWhiteEmails()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetWhiteEmail in RegisterTeacher"))
		return infrastruct.ErrorInternalServerError
	}

	for i := range whiteEmails {
		if whiteEmails[i] == user.Email {
			if err = s.confirmURL(user.Email); err != nil {
				return err
			}
			return nil
		}
	}

	err = s.email.SendMail(
		mail.RegistrationTmpl,
		fmt.Sprintf(mail.BodyWaitTeacherText, user.Email),
		fmt.Sprintf(mail.BodyWaitTeacherHTML, user.Email),
		user.Email)
	if err != nil {
		logger.LogError(errors.Wrap(err, fmt.Sprintf("err with send Email in RegisterTeacher FOR email = %s", user.Email)))
		return infrastruct.ErrorInternalServerError
	}

	logger.LogInfo(fmt.Sprintf("\nНе найден в белом списке:\nФио: %s %s %s\nПочта: %s\nГород: %s\nДолжность/место работы: %s\n",
		user.LastName, user.FirstName, user.MiddleName, user.Email, user.Address, user.Info))

	return infrastruct.ErrorNotWhiteEmail
}

func (s *Service) OrganizerRegister(user *types.OrganizerRegister) error {

	user.Role = types.RoleOrganizer
	user.Email = strings.ToLower(user.Email)
	delSpaceRegister(&user.Register)
	var err error

	if err = s.checkDoubleUserEmail(user.Email); err != nil {
		return err
	}
	if err = s.checkDoublePhoneOrganizer(user.Phone); err != nil {
		return err
	}

	//todo подумать
	//if err = s.confirmURL(user.Email); err != nil {
	//	return err
	//}

	user.ID, err = s.p.CreateUser(&user.Register)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateUser in RegisterStudent"))
		return infrastruct.ErrorInternalServerError
	}

	if err := s.p.CreateOrganizerInfo(user); err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateStudentInfo in RegisterStudent "))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}

func (s *Service) AuthUsers(auth *types.Auth) (*types.Token, error) {

	auth.Email = strings.ToLower(auth.Email)
	auth.Email = strings.ReplaceAll(auth.Email, " ", "")

	user, err := s.p.GetUserByEmail(auth.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with CheckConfirmEmail in AuthUsers"))
			return nil, infrastruct.ErrorInternalServerError
		}
		return nil, infrastruct.ErrorPasswordOrEmailIsIncorrect
	}

	if !user.ConfirmEmail {
		if user.Role != types.RoleStudent {
			return nil, infrastruct.ErrorAccountWaitConfirm
		}
		return nil, infrastruct.ErrorNOTConfirmEmail
	}

	oldUserPass, err := s.p.GetPassByUserID(user.ID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetCabinetOrganizer in AuthUsers"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if oldUserPass != auth.Pass {
		return nil, infrastruct.ErrorPasswordOrEmailIsIncorrect
	}

	token, err := infrastruct.GenerateJWT(user.ID, user.Role, s.secretKey)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with  infrastruct.GenerateJWT"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return &types.Token{Token: token}, nil
}

func (s *Service) deleteUserByID(userID int, role string) error {

	if err := s.p.DeleteUserByID(userID); err != nil {
		logger.LogError(errors.Wrap(err, "err with DeleteUserByID in RegisterUser"))
		return infrastruct.ErrorInternalServerError
	}

	switch role {
	case types.RoleStudent:
		{
			if err := s.p.DeleteDopInfoStudent(userID); err != nil {
				logger.LogError(errors.Wrap(err, "err with DeleteDopInfoStudent in RegisterUser"))
				return infrastruct.ErrorInternalServerError
			}
			return nil
		}
	case types.RoleTeacher:
		{
			if err := s.p.DeleteDopInfoTeacher(userID); err != nil {
				logger.LogError(errors.Wrap(err, "err with DeleteDopInfoTeacher in RegisterUser"))
				return infrastruct.ErrorInternalServerError
			}
			return nil
		}
	case types.RoleOrganizer:
		{
			if err := s.p.DeleteDopInfoOrganizer(userID); err != nil {
				logger.LogError(errors.Wrap(err, "err with DeleteDopInfoOrganizer in RegisterUser"))
				return infrastruct.ErrorInternalServerError
			}
			return nil
		}
	}
	return nil
}

func (s *Service) RecoverPassword(rec *types.RecoverPass) error {

	rec.Email = strings.ToLower(rec.Email)
	rec.Email = strings.ReplaceAll(rec.Email, " ", "")

	//проверяем наличие в базе по емейлу
	userInDB, err := s.p.GetUserByEmail(rec.Email)
	if err != nil {
		//если нету в базе - пользователю не говорим, но сворачиваем движуху
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with GetUserByEmail in RegisterUser"))
			return infrastruct.ErrorInternalServerError
		}
		return nil
	}

	//генерируем пароль
	tmpPass, err := uuid.GenerateUUID()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GenerateUUID in sendPassToEmail"))
		return infrastruct.ErrorInternalServerError
	}
	tmpPass = tmpPass[:8]
	tmpPass = strings.ToUpper(tmpPass)
	rec.GeneratePass = tmpPass

	//изменяем юзера
	err = s.p.UpdatePassUser(rec.GeneratePass, userInDB.ID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateUser in RegisterUser"))
		return infrastruct.ErrorInternalServerError
	}

	//todo PROD DEL
	logger.LogInfo(rec.GeneratePass)
	//отправляем на почту

	err = s.email.SendMail(mail.RecoveryPasswordTmpl,
		fmt.Sprintf(mail.RecoveryPasswordText, rec.GeneratePass),
		fmt.Sprintf(mail.RecoveryPasswordHTML, rec.GeneratePass), rec.Email)
	if err != nil {
		logger.LogError(errors.Wrap(err, fmt.Sprintf("err with sendEmail  in RecoverPassword FOR USER_ID = %d", userInDB.ID)))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}

func (s *Service) checkDoubleUserEmail(email string) error {

	//ищем пользователя с емейлом который хотят зарегистрировать
	userInDB, err := s.p.GetUserByEmail(email)
	if err != nil && err != sql.ErrNoRows {
		//если ошибка при обращении к базе - логируем
		logger.LogError(errors.Wrap(err, "err with GetUserByEmail in RegisterUser"))
		return infrastruct.ErrorInternalServerError
	}

	if userInDB == nil {
		//если никого не нашли - ок
		return nil
	}

	if !userInDB.ConfirmEmail {
		if userInDB.Role == types.RoleStudent {
			//если профиль студента не подтвержден, удаялем строчку
			if err = s.deleteUserByID(userInDB.ID, userInDB.Role); err != nil {
				return infrastruct.ErrorInternalServerError
			}
			return nil
		}
		//если это не профиль студента - говорим что надо получить подтверждение
		return infrastruct.ErrorAccountWaitConfirm
	}

	return infrastruct.ErrorEmailIsExist
}

func (s *Service) confirmURL(email string) error {

	//проверяем наличие пользователя в базе подтверждений
	have, err := s.p.CheckCodeInConfEmail(email)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with CheckDBConfPhone in confirmURL"))
			return infrastruct.ErrorInternalServerError
		}
	}
	if have {
		//если код есть - удаляем
		if err := s.p.DeleteConfLink(email); err != nil {
			logger.LogError(errors.Wrap(err, "err with DeleteConfPhone in confirmURL"))
			return infrastruct.ErrorInternalServerError
		}
	}

	//генерируем код
	now := time.Now()
	nanos := now.UnixNano()
	if err := s.p.AddCodeForConfEmail(email, nanos); err != nil {
		logger.LogError(errors.Wrap(err, "err with AddCodeForConfPhone in confirmURL"))
		return infrastruct.ErrorInternalServerError
	}

	//TODO PROD удалить код
	link := fmt.Sprintf("https://lk.bashdiktant.ru/register/link?email=%s&key=%d", email, nanos)
	logger.LogInfo(link)

	err = s.email.SendMail(
		mail.RegistrationTmpl,
		fmt.Sprintf(mail.BodyRegisterSendConfirmURLText, link),
		fmt.Sprintf(mail.BodyRegisterSendConfirmURLHTML, link),
		email)
	if err != nil {
		logger.LogError(errors.Wrap(err, fmt.Sprintf("err with SendMail in confirmURL FOR email = %s", email)))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}

func (s *Service) checkDoublePhoneOrganizer(phone string) error {

	//проверяем наличие в базе по фону
	userIDByPhone, err := s.p.GetUserIDByPhone(phone)
	if err != nil {
		//если ошибка не ноуровс
		if err != sql.ErrNoRows {
			//логируем и завершаем
			logger.LogError(errors.Wrap(err, "err with GetUserIDByPhone in checkDoublePhoneOrganizer"))
			return infrastruct.ErrorInternalServerError
		}
		return nil
	}

	//если пользователь с таким телефоном зарегистрирован проверяем подтверждение почты
	user, err := s.p.GetUserByID(userIDByPhone)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByID in checkDoublePhoneOrganizer"))
		return infrastruct.ErrorInternalServerError
	}

	if user.ConfirmEmail {
		//если подтверждена возвращаем ошибку
		return infrastruct.ErrorPhoneIsExist
	} else {
		//если не подтвержден, удаялем строчку Л
		if err = s.p.DeleteUserByID(userIDByPhone); err != nil {
			logger.LogError(errors.Wrap(err, "err with DeleteUserByID in checkDoublePhoneOrganizer"))
			return infrastruct.ErrorInternalServerError
		}
	}
	return nil
}

func (s *Service) AuthByLink(email, urlKey string) error {

	//получаем емейл пользователя по токену
	emailInDB, err := s.p.GetUserEmailByTokenInConfLink(urlKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return infrastruct.ErrorCodeIsIncorrect
		}
		logger.LogError(errors.Wrap(err, "err with CheckCodeInConfEmail in AuthByLink"))
		return infrastruct.ErrorInternalServerError
	}

	if emailInDB != email {
		return infrastruct.ErrorCodeIsIncorrect
	}

	//если токен валидный - меняем статус емейла
	if err := s.p.UpdateStatusEmail(email); err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.UpdateStatusPhone in AuthByLink"))
		return infrastruct.ErrorInternalServerError
	}

	//удаляем запись с токеном
	if err := s.p.DeleteConfLink(email); err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with s.p.DeleteConfEmail in AuthByLink"))
			return infrastruct.ErrorInternalServerError
		}
	}

	return nil
}
