package service

import (
	"database/sql"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/service/mail"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types/config"
	"github.com/AleksandrAkhapkin/bashdikt/internal/clients/postgres"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/infrastruct"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/logger"
	"github.com/pkg/errors"
	"strings"
)

type Service struct {
	p                *postgres.Postgres
	secretKey        string
	pathForDictation string
	certificate      *config.GenerateCertificate
	email            *mail.Mail
}

func NewService(pq *postgres.Postgres, mail *mail.Mail, cnf *config.Config) (*Service, error) {

	return &Service{
		p:                pq,
		secretKey:        cnf.SecretKeyJWT,
		pathForDictation: cnf.PathForDictation,
		certificate:      cnf.GenerateCertificate,
		email:            mail,
	}, nil
}

func delSpaceRegister(user *types.Register) {
	user.FirstName = strings.ReplaceAll(user.FirstName, " ", "")
	user.LastName = strings.ReplaceAll(user.LastName, " ", "")
	user.MiddleName = strings.ReplaceAll(user.MiddleName, " ", "")
	user.Email = strings.ReplaceAll(user.Email, " ", "")
	user.Pass = strings.ReplaceAll(user.Pass, " ", "")
	user.RepeatPass = strings.ReplaceAll(user.RepeatPass, " ", "")
}

func (s *Service) CheckUserInDBUsers(id int) error {
	if err := s.p.CheckUserInDBUsers(id); err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with CheckUserInDBUsers"))
			return infrastruct.ErrorInternalServerError
		}
		return infrastruct.ErrorPermissionDenied
	}

	return nil
}

func (s *Service) SaveRequestLog(body []byte, route string, ip string) error {
	return s.p.SaveRequestLog(string(body), route, ip)
}
