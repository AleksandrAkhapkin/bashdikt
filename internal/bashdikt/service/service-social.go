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
)

func (s *Service) RegisterOrAuthStudentBySocial(auth *types.RegisterOrAuthSocial) (*types.Token, error) {

	if auth.Email == "" {
		//если соцсетка не вернула емейл - говорим об ошибке
		return nil, infrastruct.ErrorNotEmailByAuth
	}

	auth.Email = strings.ToLower(auth.Email)
	auth.Email = strings.ReplaceAll(auth.Email, " ", "")

	userInDB, err := s.p.GetUserByEmail(auth.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with CheckConfirmEmail in RegisterOrAuthStudentBySocial"))
			return nil, infrastruct.ErrorInternalServerError
		}
	}

	if userInDB != nil {
		//авторизуем если нашли аккаунт
		if !userInDB.ConfirmEmail && userInDB.Role != types.RoleStudent {
			//аккаунт не студента ожидает подтверждения
			return nil, infrastruct.ErrorAccountWaitConfirm
		}

		//авторизуем аккаунт
		if !userInDB.ConfirmEmail {
			//подтверждаем почту
			if err := s.p.UpdateStatusEmail(userInDB.Email); err != nil {
				logger.LogError(errors.Wrap(err, "err with UpdateStatusEmail in RegisterOrAuthStudentBySocial "))
				return nil, infrastruct.ErrorInternalServerError
			}
		}

		token, err := infrastruct.GenerateJWT(userInDB.ID, userInDB.Role, s.secretKey)
		if err != nil {
			logger.LogError(err)
			return nil, infrastruct.ErrorInternalServerError
		}

		return &types.Token{Token: token}, nil
	}

	//если пользователь небыл найден - регистрируем
	user := &types.StudentRegister{}

	newPass, err := uuid.GenerateUUID()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GenerateUUID in RegisterOrAuthStudentBySocial"))
		return nil, infrastruct.ErrorInternalServerError
	}
	user.FirstName = auth.Firstname
	user.LastName = auth.Lastname
	user.Email = auth.Email
	user.Pass = newPass[:8]
	user.Role = types.RoleStudent

	user.ID, err = s.p.CreateUser(&user.Register)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateUser in RegisterOrAuthStudentBySocial"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if err := s.p.CreateStudentInfo(user); err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateStudentInfo in RegisterOrAuthStudentBySocial "))
		return nil, infrastruct.ErrorInternalServerError
	}

	err = s.email.SendMail(
		mail.RegistrationTmpl,
		fmt.Sprintf(mail.BodyRegisterBySocialText, user.Pass),
		fmt.Sprintf(mail.BodyRegisterBySocialHTML, user.Pass),
		user.Email)
	if err != nil {
		logger.LogError(errors.Wrap(err, fmt.Sprintf("err with send Email in RegisterOrAuthStudentBySocial FOR USER_ID = %d", user.ID)))
		return nil, infrastruct.ErrorInternalServerError
	}

	token, err := infrastruct.GenerateJWT(user.ID, user.Role, s.secretKey)
	if err != nil {
		logger.LogError(err)
		return nil, infrastruct.ErrorInternalServerError
	}

	if err := s.p.UpdateStatusEmail(user.Email); err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.UpdateStatusEmail in RegisterOrAuthStudentBySocial"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return &types.Token{Token: token}, nil

}

func (s *Service) RegisterOrAuthTeacherBySocial(auth *types.RegisterOrAuthSocial) (*types.Token, error) {

	if auth.Email == "" {
		//если соцсетка не вернула емейл - говорим об ошибке
		return nil, infrastruct.ErrorNotEmailByAuth
	}

	auth.Email = strings.ToLower(auth.Email)
	auth.Email = strings.ReplaceAll(auth.Email, " ", "")

	userInDB, err := s.p.GetUserByEmail(auth.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with CheckConfirmEmail in RegisterOrAuthTeacherBySocial"))
			return nil, infrastruct.ErrorInternalServerError
		}
	}

	if userInDB != nil {
		//авторизуем если нашли аккаунт
		if !userInDB.ConfirmEmail && userInDB.Role != types.RoleStudent {
			//аккаунт не студента ожидает подтверждения
			return nil, infrastruct.ErrorAccountWaitConfirm
		}

		//авторизуем аккаунт
		if !userInDB.ConfirmEmail {
			//подтверждаем почту
			if err := s.p.UpdateStatusEmail(userInDB.Email); err != nil {
				logger.LogError(errors.Wrap(err, "err with UpdateStatusEmail in RegisterOrAuthTeacherBySocial "))
				return nil, infrastruct.ErrorInternalServerError
			}
		}

		token, err := infrastruct.GenerateJWT(userInDB.ID, userInDB.Role, s.secretKey)
		if err != nil {
			logger.LogError(err)
			return nil, infrastruct.ErrorInternalServerError
		}

		return &types.Token{Token: token}, nil
	}

	//если пользователь небыл найден - регистрируем
	user := &types.TeacherRegister{}

	newPass, err := uuid.GenerateUUID()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GenerateUUID in RegisterOrAuthTeacherBySocial"))
		return nil, infrastruct.ErrorInternalServerError
	}
	user.FirstName = auth.Firstname
	user.LastName = auth.Lastname
	user.Email = auth.Email
	user.Pass = newPass[:8]
	user.Role = types.RoleTeacher

	user.ID, err = s.p.CreateUser(&user.Register)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateUser in RegisterOrAuthTeacherBySocial"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if err := s.p.CreateTeacherInfo(user); err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateStudentInfo in RegisterOrAuthTeacherBySocial "))
		return nil, infrastruct.ErrorInternalServerError
	}

	whiteEmails, err := s.p.GetWhiteEmails()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetWhiteEmail in RegisterOrAuthTeacherBySocial"))
		return nil, infrastruct.ErrorInternalServerError
	}

	for i := range whiteEmails {
		if whiteEmails[i] == user.Email {
			if err := s.p.UpdateStatusEmail(user.Email); err != nil {
				logger.LogError(errors.Wrap(err, "err with UpdateStatusEmail in RegisterOrAuthTeacherBySocial"))
			}

			token, err := infrastruct.GenerateJWT(user.ID, user.Role, s.secretKey)
			if err != nil {
				logger.LogError(errors.Wrap(err, "err with  infrastruct.GenerateJWT in RegisterOrAuthTeacherBySocial"))
				return nil, infrastruct.ErrorInternalServerError
			}

			return &types.Token{Token: token}, nil
		}
	}
	err = s.email.SendMail(
		mail.RegistrationTmpl,
		fmt.Sprintf(mail.BodyWaitTeacherForSocialText, user.Email, user.Pass),
		fmt.Sprintf(mail.BodyWaitTeacherForSocialHTML, user.Email, user.Pass),
		user.Email)
	if err != nil {
		logger.LogError(errors.Wrap(err, fmt.Sprintf("err with send Email in RegisterTeacher FOR email = %d", user.Email)))
		return nil, infrastruct.ErrorInternalServerError
	}

	logger.LogInfo(fmt.Sprintf("\nНе найден в белом списке:\nФио: %s %s %s\nПочта: %s\n%s\n%s\n",
		user.LastName, user.FirstName, user.MiddleName, user.Email, user.Address, user.Info))

	return nil, infrastruct.ErrorNotWhiteEmail
}

func (s *Service) AuthSocial(auth *types.RegisterOrAuthSocial) (*types.Token, error) {

	auth.Email = strings.ToLower(auth.Email)
	auth.Email = strings.ReplaceAll(auth.Email, " ", "")

	user, err := s.p.GetUserByEmail(auth.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with CheckConfirmEmail in AuthUsers"))
			return nil, infrastruct.ErrorInternalServerError
		}
		return nil, infrastruct.ErrorNotHaveUser
	}

	if !user.ConfirmEmail {
		if user.Role != types.RoleStudent {
			return nil, infrastruct.ErrorAccountWaitConfirm
		}
	}

	if err := s.p.UpdateStatusEmail(auth.Email); err != nil {
		logger.LogError(errors.Wrap(err, "err with UpdateStatusEmail in AuthSocial"))
		return nil, infrastruct.ErrorInternalServerError
	}

	token, err := infrastruct.GenerateJWT(user.ID, user.Role, s.secretKey)
	if err != nil {
		logger.LogError(err)
		return nil, infrastruct.ErrorInternalServerError
	}

	return &types.Token{Token: token}, nil
}
