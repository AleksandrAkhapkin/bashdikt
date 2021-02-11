package service

import (
	"database/sql"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/infrastruct"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/logger"
	"github.com/pkg/errors"
)

func (s *Service) GetStudentForOrgCabinet(getInfo *types.GetInfoForCabinet) (*types.AllStudentsForCabinet, error) {

	var err error
	allStudents := &types.AllStudentsForCabinet{}
	allStudents.AllStudents, err = s.p.GetAllStudentConfirmEmail(getInfo)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.GetAllStudentConfirmEmail"))
		return nil, infrastruct.ErrorInternalServerError
	}
	allStudents.Total, err = s.p.TotalUserByRole(getInfo.Role)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with TotalUserByRole"))
		return nil, infrastruct.ErrorInternalServerError
	}

	for i, _ := range allStudents.AllStudents {
		status, err := s.p.GetStatusDicByUserID(allStudents.AllStudents[i].ID)
		if err != nil {
			if err != sql.ErrNoRows {
				logger.LogError(errors.Wrap(err, "err with GetStatusDicByUserID in GetStudentForOrgCabinet"))
				return nil, infrastruct.ErrorInternalServerError
			}
			status = ""
		}
		allStudents.AllStudents[i].Status = status
	}

	return allStudents, nil
}

func (s *Service) GetTeacherForOrgCabinet(getInfo *types.GetInfoForCabinet) (*types.AllTeachersForCabinet, error) {

	var err error
	allTeachers := &types.AllTeachersForCabinet{}
	allTeachers.AllTeachers, err = s.p.GetTeachersForOrganizerInCabinet(getInfo)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.GetTeachersForOrganizerInCabinet"))
		return nil, infrastruct.ErrorInternalServerError
	}
	allTeachers.Total, err = s.p.TotalUserByRole(getInfo.Role)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with TotalUserByRole"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return allTeachers, nil
}

func (s *Service) GetStudentForTeacherCabinet(getInfo *types.GetInfoForCabinet) (*types.AllStudentsForCabinet, error) {

	var err error
	allStudents := &types.AllStudentsForCabinet{}
	allStudents.AllStudents, err = s.p.GetAllStudentConfirmEmail(getInfo)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.GetStudentForTeacherCabinet"))
		return nil, infrastruct.ErrorInternalServerError
	}
	allStudents.Total, err = s.p.TotalUserByRole(getInfo.Role)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with TotalUserByRole"))
		return nil, infrastruct.ErrorInternalServerError
	}
	for i, _ := range allStudents.AllStudents {
		status, err := s.p.GetStatusDicByUserID(allStudents.AllStudents[i].ID)
		if err != nil {
			if err != sql.ErrNoRows {
				logger.LogError(errors.Wrap(err, "err with GetStatusDicByUserID in GetStudentForOrgCabinet"))
				return nil, infrastruct.ErrorInternalServerError
			}
			status = ""
		}
		allStudents.AllStudents[i].Status = status
	}

	return allStudents, nil
}

func (s *Service) GetPinStudentForTeacherCabinet(getInfo *types.GetInfoForCabinet, teacherID int) (*types.AllStudentsForCabinet, error) {

	var err error
	allPinStudents := &types.AllStudentsForCabinet{}
	allPinStudents.AllStudents, err = s.p.GetPinStudents(getInfo, teacherID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.GetPinStudentForTeacherCabinet in GetPinStudentForTeacherCabinet"))
		return nil, infrastruct.ErrorInternalServerError
	}
	allPinStudents.Total, err = s.p.TotalPinStudents(teacherID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with TotalPinStudents in GetPinStudentForTeacherCabinet "))
		return nil, infrastruct.ErrorInternalServerError
	}
	for i, _ := range allPinStudents.AllStudents {
		status, err := s.p.GetStatusDicByUserID(allPinStudents.AllStudents[i].ID)
		if err != nil {
			if err != sql.ErrNoRows {
				logger.LogError(errors.Wrap(err, "err with GetStatusDicByUserID in GetPinStudentForTeacherCabinet"))
				return nil, infrastruct.ErrorInternalServerError
			}
			status = ""
		}
		allPinStudents.AllStudents[i].Status = status
	}

	return allPinStudents, nil
}

func (s *Service) GetMyDictation(userID int) (*types.MyDictation, error) {

	namesFile, err := s.GetDictationsNameByStudentID(userID)
	if err != nil {
		if err != infrastruct.ErrorNotHaveDictant {
			return nil, err
		}
		myDictation := &types.MyDictation{UserID: userID, Status: types.StatusStudentNotWrite}
		return myDictation, nil
	}

	if namesFile.Names[0] == "Online" {
		dictation, err := s.p.GetMyDictation(userID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with MyDictation in GetMyDictation"))
			return nil, infrastruct.ErrorInternalServerError
		}
		dictation.Names = namesFile.Names
		markers, err := s.p.GetMarkers(userID)
		if err != nil {
			if err != sql.ErrNoRows {
				logger.LogError(errors.Wrap(err, "err with GetMarkers in GetMyDictation"))
				return nil, infrastruct.ErrorInternalServerError
			}
			markers = []types.Markers{}
		}

		dictation.Markers = markers
		return dictation, nil
	}

	dictation, err := s.p.GetRatingAndStatusDictation(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetRatingAndStatusDictation in GetMyDictation"))
		return nil, infrastruct.ErrorInternalServerError
	}
	dictation.UserID = userID
	dictation.Names = namesFile.Names

	return dictation, nil
}

func (s *Service) GetMyCertificate(userID int, role string, lang string) error {

	user, err := s.p.GetUserByID(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByID in MakeAndSendMyCertificate"))
		return infrastruct.ErrorInternalServerError
	}

	cert := &types.MakeCertificate{}
	cert.FirstName = user.FirstName
	cert.LastName = user.LastName
	cert.MiddleName = user.MiddleName
	cert.EmailTo = user.Email
	cert.PathForCert = s.certificate.PathForStudentTemplateRus //for rus stud
	if lang == "bash" {
		cert.PathForCert = s.certificate.PathForStudentTemplateBash //for bash stud
	}
	if role == types.RoleTeacher {
		cert.PathForCert = s.certificate.PathForTeacherTemplateRus //for rus teach
		if lang == "bash" {
			cert.PathForCert = s.certificate.PathForTeacherTemplateBash //for bash teach
		}
	}
	cert.PathForFonts = s.certificate.PathForFonts
	cert.PathForSave = s.certificate.PathForSaveCert
	cert.UserID = userID

	go func() {
		if err = s.makeAndSendCertificate(cert); err != nil {
			logger.LogError(errors.Wrap(err, "err with GetMyCertificate"))
		}
	}()

	if err := s.p.ChangeSendCertStatus(userID); err != nil {
		logger.LogError(errors.Wrap(err, "err with ChangeSendCertStatus in GetMyDictation"))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}
