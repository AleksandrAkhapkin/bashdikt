package postgres

func (p *Postgres) AddCodeForConfEmail(email string, code int64) error {

	if _, err := p.db.Exec("INSERT INTO confirm_url (email, url_key) VALUES ($1, $2)", email, code); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) CheckCodeInConfEmail(email string) (bool, error) {

	if err := p.db.QueryRow("SELECT email FROM confirm_url WHERE email = $1", email).Scan(&email); err != nil {
		return false, err
	}

	return true, nil
}

func (p *Postgres) DeleteConfLink(email string) error {

	var err error
	_, err = p.db.Exec("DELETE FROM confirm_url WHERE email = $1", email)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) UpdateStatusEmail(email string) error {

	_, err := p.db.Exec("UPDATE users SET confirm_email = true WHERE email = $1", email)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetUserIDByPhone(phone string) (int, error) {

	var userID int
	err := p.db.QueryRow("SELECT user_id FROM organizer_info WHERE phone = $1", phone).
		Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (p *Postgres) GetUserEmailByTokenInConfLink(urlKey string) (string, error) {

	var email string

	err := p.db.QueryRow("SELECT email FROM confirm_url WHERE url_key = $1", urlKey).Scan(&email)
	if err != nil {
		return "", err
	}

	return email, nil
}
