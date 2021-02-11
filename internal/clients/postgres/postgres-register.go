package postgres

import (
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
)

func (p *Postgres) CreateUser(user *types.Register) (int, error) {

	var id int
	err := p.db.QueryRow("INSERT INTO users (user_role, email, pass, last_name, first_name, middle_name, address) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING user_id",
		user.Role, user.Email, user.Pass, user.LastName, user.FirstName, user.MiddleName, user.Address).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *Postgres) CreateStudentInfo(user *types.StudentRegister) error {

	_, err := p.db.Exec("INSERT INTO student_info (user_id, level) VALUES ($1, $2)",
		user.ID, user.Level)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) CreateTeacherInfo(user *types.TeacherRegister) error {

	_, err := p.db.Exec("INSERT INTO teacher_info (user_id, info) VALUES ($1, $2)", user.ID, user.Info)
	if err != nil {
		return err
	}

	return nil
}
func (p *Postgres) CreateOrganizerInfo(user *types.OrganizerRegister) error {

	_, err := p.db.Exec("INSERT INTO organizer_info (user_id, phone, soc_url, count_student, format_dictation) "+
		"VALUES ($1, $2, $3, $4, $5)", user.ID, user.Phone, user.SocURL, user.CountStudent, user.FormatDictation)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetWhiteEmails() ([]string, error) {

	var emails []string
	var email string

	rows, err := p.db.Query("SELECT email FROM white_emails")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&email); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	return emails, nil
}
