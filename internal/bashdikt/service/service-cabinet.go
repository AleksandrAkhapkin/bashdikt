package service

import (
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/infrastruct"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/logger"
	"github.com/pkg/errors"
)

func (s *Service) GetCabinetStudent(id int) (*types.StudentProfile, error) {

	cabinet := &types.StudentProfile{}

	cabinet, err := s.p.GetCabinetStudent(id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetStudentCabinet"))
		return nil, infrastruct.ErrorInternalServerError
	}
	cabinet.ID = id
	cabinet.Role = types.RoleStudent

	return cabinet, nil
}

func (s *Service) GetCabinetTeacher(id int) (*types.TeacherProfile, error) {

	cabinet := &types.TeacherProfile{}

	cabinet, err := s.p.GetCabinetTeacher(id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetCabinetTeacher"))
		return nil, infrastruct.ErrorInternalServerError
	}
	cabinet.ID = id
	cabinet.Role = types.RoleTeacher

	return cabinet, nil
}

func (s *Service) GetCabinetOrganizer(id int) (*types.OrganizerProfile, error) {

	cabinet := &types.OrganizerProfile{}

	cabinet, err := s.p.GetCabinetOrganizer(id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetCabinetOrganizer"))
		return nil, infrastruct.ErrorInternalServerError
	}
	cabinet.ID = id
	cabinet.Role = types.RoleOrganizer

	//что бы в выдаче джисона не было null
	if len(cabinet.AddEmail) == 0 {
		cabinet.AddEmail = make([]string, 0)
	}
	if len(cabinet.AddPhone) == 0 {
		cabinet.AddPhone = make([]string, 0)
	}

	return cabinet, nil
}

func (s *Service) PutCabinetStudent(newInfoForCabinet *types.StudentProfileForPUT) (*types.StudentProfile, error) {

	if newInfoForCabinet.Pass != "" {
		oldUserPass, err := s.p.GetPassByUserID(newInfoForCabinet.ID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetCabinetStudent in PutCabinetStudent"))
			return nil, infrastruct.ErrorInternalServerError
		}
		if oldUserPass == newInfoForCabinet.OldPass {
			if err := s.p.UpdatePassUser(newInfoForCabinet.Pass, newInfoForCabinet.ID); err != nil {
				logger.LogError(errors.Wrap(err, "err with UpdatePassUser in PutCabinetStudent"))
				return nil, infrastruct.ErrorInternalServerError
			}
		} else {
			return nil, infrastruct.ErrorOldPassDontMatch
		}
	}

	err := s.p.PutCabinetStudent(newInfoForCabinet)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with PutStudentCabinet in PutCabinetStudent"))
		return nil, infrastruct.ErrorInternalServerError
	}

	newUserInfo, err := s.GetCabinetStudent(newInfoForCabinet.ID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetCabinetStudent in PutCabinetStudent"))
		return nil, err
	}

	return newUserInfo, nil
}

func (s *Service) PutCabinetTeacher(newInfoForCabinet *types.TeacherProfileForPUT) (*types.TeacherProfile, error) {

	if newInfoForCabinet.Pass != "" {
		oldUserPass, err := s.p.GetPassByUserID(newInfoForCabinet.ID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetCabinetTeacher in PutCabinetTeacher"))
			return nil, infrastruct.ErrorInternalServerError
		}
		if oldUserPass == newInfoForCabinet.OldPass {
			if err := s.p.UpdatePassUser(newInfoForCabinet.Pass, newInfoForCabinet.ID); err != nil {
				logger.LogError(errors.Wrap(err, "err with UpdatePassUser in PutCabinetTeacher"))
				return nil, infrastruct.ErrorInternalServerError
			}
		} else {
			return nil, infrastruct.ErrorOldPassDontMatch
		}
	}

	err := s.p.PutCabinetTeacher(newInfoForCabinet)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with PutCabinetTeacher in GetCabinetTeacher"))
		return nil, infrastruct.ErrorInternalServerError
	}

	newUserInfo, err := s.GetCabinetTeacher(newInfoForCabinet.ID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetCabinetTeacher in GetCabinetTeacher"))
		return nil, err
	}

	return newUserInfo, nil
}

func (s *Service) PutCabinetOrganizer(newInfoForCabinet *types.OrganizerProfileForPUT) (*types.OrganizerProfile, error) {

	if newInfoForCabinet.Pass != "" {
		oldUserPass, err := s.p.GetPassByUserID(newInfoForCabinet.ID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetCabinetOrganizer in PutCabinetOrganizer"))
			return nil, infrastruct.ErrorInternalServerError
		}
		if oldUserPass == newInfoForCabinet.OldPass {
			if err := s.p.UpdatePassUser(newInfoForCabinet.Pass, newInfoForCabinet.ID); err != nil {
				logger.LogError(errors.Wrap(err, "err with UpdatePassUser in PutCabinetOrganizer"))
				return nil, infrastruct.ErrorInternalServerError
			}
		} else {
			return nil, infrastruct.ErrorOldPassDontMatch
		}
	}

	err := s.p.PutCabinetOrganizer(newInfoForCabinet)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with PutCabinetOrganizer in PutCabinetOrganizer"))
		return nil, infrastruct.ErrorInternalServerError
	}

	newUserInfo, err := s.GetCabinetOrganizer(newInfoForCabinet.ID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetCabinetOrganizer in PutCabinetOrganizer"))
		return nil, err
	}

	return newUserInfo, nil
}
