package postgres

import "github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"

func (p *Postgres) GetUserIDByEmail(email string) (int, error) {

	var id int
	if err := p.db.QueryRow("SELECT user_id FROM users WHERE email = $1", email).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (p *Postgres) GetUserByEmail(email string) (*types.User, error) {

	user := &types.User{}
	if err := p.db.QueryRow("SELECT user_id, user_role, email, confirm_email FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Role, &user.Email, &user.ConfirmEmail); err != nil {
		return nil, err
	}

	return user, nil
}

func (p *Postgres) GetOrganizerIDByPhone(phone string) (int, error) {

	var id int
	if err := p.db.QueryRow("SELECT user_id FROM organizer_info WHERE phone = $1", phone).
		Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (p *Postgres) CheckUserInDBUsers(userID int) error {

	if err := p.db.QueryRow("SELECT FROM users WHERE user_id = $1", userID).Scan(); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) TotalUserByRole(role string) (int, error) {

	var count int
	if err := p.db.QueryRow("SELECT COUNT (user_id) FROM users WHERE user_role = $1 AND confirm_email = true", role).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (p *Postgres) GetPassByUserID(userID int) (string, error) {

	var passInDB string
	err := p.db.QueryRow("SELECT pass FROM users WHERE user_id = $1", userID).Scan(&passInDB)
	if err != nil {
		return "", err
	}

	return passInDB, nil
}

func (p *Postgres) GetUserByID(userID int) (*types.User, error) {

	user := &types.User{ID: userID}
	if err := p.db.QueryRow("SELECT last_name, first_name, middle_name, user_role, email, confirm_email FROM users WHERE user_id = $1", userID).
		Scan(&user.LastName, &user.FirstName, &user.MiddleName, &user.Role, &user.Email, &user.ConfirmEmail); err != nil {
		return nil, err
	}

	return user, nil
}

func (p *Postgres) TotalPinStudents(teacherID int) (int, error) {

	var count int
	if err := p.db.QueryRow("SELECT COUNT (student_id) FROM pin_student WHERE teacher_id = $1", teacherID).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}
