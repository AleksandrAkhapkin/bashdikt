package postgres

import (
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"github.com/pkg/errors"
)

func (p *Postgres) GetAllStudentConfirmEmail(getInfo *types.GetInfoForCabinet) ([]types.AllStudents, error) {

	students := make([]types.AllStudents, 0)
	rows, err := p.db.Query("SELECT users.user_id, user_role, email, last_name, first_name, middle_name, address, "+
		"level "+
		"FROM users RIGHT JOIN student_info ON users.user_id = student_info.user_id "+
		"WHERE user_role = 'student' AND confirm_email = true ORDER BY user_id ASC LIMIT $1 OFFSET $2",
		getInfo.Limit, getInfo.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "err with Query")
	}
	defer rows.Close()
	student := types.AllStudents{}
	for rows.Next() {
		if err = rows.Scan(&student.ID, &student.Role, &student.Email, &student.LastName, &student.FirstName, &student.MiddleName, &student.Address,
			&student.Level); err != nil {
			return nil, errors.Wrap(err, "err with Scan")
		}
		students = append(students, student)
	}

	return students, nil
}

func (p *Postgres) GetTeachersForOrganizerInCabinet(getInfo *types.GetInfoForCabinet) ([]types.AllTeachers, error) {

	teachers := make([]types.AllTeachers, 0)
	rows, err := p.db.Query("SELECT users.user_id, user_role, email, last_name, first_name, middle_name, "+
		"address, info FROM users RIGHT JOIN teacher_info ON users.user_id = teacher_info.user_id "+
		"WHERE user_role = 'teacher' AND confirm_email = true ORDER BY user_id ASC LIMIT $1 OFFSET $2", getInfo.Limit, getInfo.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "err with Query")
	}
	defer rows.Close()
	teacher := types.AllTeachers{}
	for rows.Next() {
		if err = rows.Scan(&teacher.ID, &teacher.Role, &teacher.Email, &teacher.LastName, &teacher.FirstName,
			&teacher.MiddleName, &teacher.Address, &teacher.InfoAboutWork); err != nil {
			return nil, errors.Wrap(err, "err with Scan")
		}
		teachers = append(teachers, teacher)
	}

	return teachers, nil
}

func (p *Postgres) GetPinStudents(getInfo *types.GetInfoForCabinet, teacherID int) ([]types.AllStudents, error) {

	students := make([]types.AllStudents, 0)
	rows, err := p.db.Query("SELECT users.user_id, user_role, email, last_name, first_name, middle_name, address, "+
		"level "+
		"FROM users RIGHT JOIN student_info ON users.user_id = student_info.user_id "+
		"RIGHT JOIN pin_student ON pin_student.student_id = student_info.user_id "+
		"WHERE pin_student.teacher_id = $3 ORDER BY user_id ASC LIMIT $1 OFFSET $2",
		getInfo.Limit, getInfo.Offset, teacherID)
	if err != nil {
		return nil, errors.Wrap(err, "err with Query")
	}
	defer rows.Close()
	student := types.AllStudents{}
	for rows.Next() {
		if err = rows.Scan(&student.ID, &student.Role, &student.Email, &student.LastName, &student.FirstName, &student.MiddleName, &student.Address,
			&student.Level); err != nil {
			return nil, errors.Wrap(err, "err with Scan")
		}
		students = append(students, student)
	}

	return students, nil
}
