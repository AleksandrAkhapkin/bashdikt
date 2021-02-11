package postgres

import (
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func (p *Postgres) GetCabinetStudent(id int) (*types.StudentProfile, error) {

	profile := &types.StudentProfile{}
	err := p.db.QueryRow("SELECT email, last_name, first_name, middle_name, address, level "+
		"FROM users RIGHT JOIN student_info ON users.user_id = student_info.user_id WHERE users.user_id = $1", id).
		Scan(&profile.Email, &profile.LastName,
			&profile.FirstName, &profile.MiddleName, &profile.Address, &profile.Level)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (p *Postgres) GetCabinetTeacher(id int) (*types.TeacherProfile, error) {

	profile := &types.TeacherProfile{}
	err := p.db.QueryRow("SELECT email, last_name, first_name, middle_name, address, info "+
		"FROM users RIGHT JOIN teacher_info ON users.user_id=teacher_info.user_id WHERE users.user_id = $1", id).
		Scan(&profile.Email, &profile.LastName, &profile.FirstName, &profile.MiddleName, &profile.Address,
			&profile.Info)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (p *Postgres) GetCabinetOrganizer(id int) (*types.OrganizerProfile, error) {

	profile := &types.OrganizerProfile{}
	err := p.db.QueryRow("SELECT email, last_name, first_name, middle_name, address, phone, soc_url, "+
		"count_student, format_dictation, add_phone, add_email FROM users RIGHT JOIN organizer_info "+
		"ON users.user_id=organizer_info.user_id WHERE users.user_id = $1", id).Scan(
		&profile.Email, &profile.LastName, &profile.FirstName, &profile.MiddleName,
		&profile.Address, &profile.Phone, &profile.SocURL, &profile.CountStudent, &profile.FormatDictation,
		pq.Array(&profile.AddPhone), pq.Array(&profile.AddEmail))
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (p *Postgres) PutCabinetStudent(cabinet *types.StudentProfileForPUT) error {

	tx, err := p.db.Begin()
	if err != nil {
		return errors.Wrap(err, "err with Begin")
	}

	_, err = tx.Exec("UPDATE users SET email = $1, last_name = $2, first_name = $3, "+
		"middle_name = $4, address = $5, updated_at = NOW() WHERE user_id = $6",
		cabinet.Email, cabinet.LastName, cabinet.FirstName, cabinet.MiddleName,
		cabinet.Address, cabinet.ID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "err with Exec users")
	}

	_, err = tx.Exec("UPDATE student_info SET level = $1 WHERE user_id = $2",
		cabinet.Level,
		cabinet.ID,
	)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "err with Exec organizer_info")
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "err with Commit")
	}

	return nil
}

func (p *Postgres) PutCabinetTeacher(cabinet *types.TeacherProfileForPUT) error {

	tx, err := p.db.Begin()
	if err != nil {
		return errors.Wrap(err, "err with Begin")
	}

	_, err = tx.Exec("UPDATE users SET email = $1, last_name = $2, first_name = $3, middle_name = $4, "+
		"address = $5, updated_at = NOW() WHERE user_id = $6",
		cabinet.Email, cabinet.LastName, cabinet.FirstName, cabinet.MiddleName, cabinet.Address, cabinet.ID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "err with Exec users")
	}

	_, err = tx.Exec("UPDATE teacher_info SET info = $1 WHERE user_id = $2", cabinet.Info, cabinet.ID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "err with Exec teacher_info")
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "err with Commit")
	}

	return nil
}

func (p *Postgres) PutCabinetOrganizer(cabinet *types.OrganizerProfileForPUT) error {

	tx, err := p.db.Begin()
	if err != nil {
		return errors.Wrap(err, "err with Begin")
	}

	_, err = tx.Exec("UPDATE users SET email = $1, last_name = $2, first_name = $3, "+
		"middle_name = $4, address = $5, updated_at = NOW() WHERE user_id = $6",
		cabinet.Email, cabinet.LastName, cabinet.FirstName, cabinet.MiddleName,
		cabinet.Address, cabinet.ID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "err with Exec users")
	}

	_, err = tx.Exec("UPDATE organizer_info SET phone = $1, soc_url = $2, count_student = $3, "+
		"format_dictation = $4, add_phone = $5, add_email = $6 WHERE user_id = $7", cabinet.Phone, cabinet.SocURL,
		cabinet.CountStudent, cabinet.FormatDictation, pq.Array(cabinet.AddPhone), pq.Array(cabinet.AddEmail), cabinet.ID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "err with Exec organizer_info")
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "err with Commit")
	}

	return nil
}

func (p *Postgres) UpdatePassUser(newPass string, userID int) error {

	_, err := p.db.Exec("UPDATE users SET pass = $1 WHERE user_id = $2", newPass, userID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) DeleteUserByID(id int) error {

	_, err := p.db.Exec("DELETE FROM users WHERE user_id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) DeleteDopInfoStudent(id int) error {

	_, err := p.db.Exec("DELETE FROM student_info WHERE user_id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
func (p *Postgres) DeleteDopInfoTeacher(id int) error {

	_, err := p.db.Exec("DELETE FROM teacher_info WHERE user_id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
func (p *Postgres) DeleteDopInfoOrganizer(id int) error {

	_, err := p.db.Exec("DELETE FROM organizer_info WHERE user_id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
