package postgres

func (p *Postgres) GetLastPinTeacher() (int, error) {

	var lastTeacher int
	err := p.db.QueryRow("SELECT teacher_id FROM pin_student ORDER by time_pin DESC LIMIT 1").Scan(&lastTeacher)
	if err != nil {
		return 0, err
	}

	return lastTeacher, nil
}

func (p *Postgres) NextAfterLastTeacherPin(lastTeacher int) (int, error) {

	var nextAfterLastTeacher int
	err := p.db.QueryRow("SELECT user_id FROM users WHERE user_id > $1 AND user_role = 'teacher' ORDER BY user_id ASC LIMIT 1", lastTeacher).Scan(&nextAfterLastTeacher)
	if err != nil {
		err := p.db.QueryRow("SELECT user_id FROM users WHERE user_id > 0 AND user_role = 'teacher' ORDER BY user_id ASC LIMIT 1").Scan(&nextAfterLastTeacher)
		if err != nil {
			return 0, err
		}
	}
	return nextAfterLastTeacher, nil
}

func (p *Postgres) GetRandTeacher() (int, error) {

	var randTeacher int
	err := p.db.QueryRow("SELECT user_id FROM users WHERE user_role = 'teacher' LIMIT 1 ").Scan(&randTeacher)
	if err != nil {
		return 0, err
	}

	return randTeacher, nil
}

func (p *Postgres) PinnedStudent(stdntID, tchrID int, format string) error {

	_, err := p.db.Exec("INSERT INTO pin_student (student_id, teacher_id, format_dictation) VALUES ($1, $2, $3)",
		stdntID, tchrID, format)
	if err != nil {
		return err
	}

	return nil
}
