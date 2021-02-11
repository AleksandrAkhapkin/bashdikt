package postgres

import (
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"github.com/pkg/errors"
)

func (p *Postgres) AddOnlineDictation(dictation *types.Dictation) error {

	_, err := p.db.Exec("INSERT INTO dictation (text, user_id, status, online) VALUES ($1, $2, $3, true)",
		dictation.Text, dictation.UserID, types.StatusStudentNotChecked)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) CheckOnlineWrite(userID int) (bool, error) {

	var online bool
	err := p.db.QueryRow("SELECT online FROM dictation WHERE user_id = $1", userID).Scan(&online)
	if err != nil {
		return false, err
	}

	return online, err
}

func (p *Postgres) GetOnlineDictationById(userID int) (*types.Dictation, error) {

	dictation := &types.Dictation{UserID: userID}
	err := p.db.QueryRow("SELECT text FROM dictation WHERE user_id = $1", dictation.UserID).Scan(&dictation.Text)
	if err != nil {
		return nil, err
	}

	return dictation, err
}

func (p *Postgres) ChangeRatingAndStatus(dictation *types.ReplyDictation) error {

	_, err := p.db.Exec("UPDATE dictation SET rating = $1, status = $2  WHERE user_id = $3",
		dictation.Rating, types.StatusStudentChecked, dictation.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetMyDictation(userID int) (*types.MyDictation, error) {

	dictation := &types.MyDictation{UserID: userID}
	err := p.db.QueryRow("SELECT text, rating, status, send_cert FROM dictation WHERE user_id = $1", dictation.UserID).
		Scan(&dictation.Text, &dictation.Rating, &dictation.Status, &dictation.SendCert)
	if err != nil {
		return nil, err
	}

	return dictation, err
}

func (p *Postgres) GetRatingAndStatusDictation(userID int) (*types.MyDictation, error) {

	dictation := &types.MyDictation{UserID: userID}
	err := p.db.QueryRow("SELECT rating, status, send_cert FROM dictation WHERE user_id = $1", dictation.UserID).
		Scan(&dictation.Rating, &dictation.Status, &dictation.SendCert)
	if err != nil {
		return nil, err
	}

	return dictation, err
}

func (p *Postgres) AddOfflineDictation(userID int) error {

	_, err := p.db.Exec("INSERT INTO dictation (user_id, status, online) VALUES ($1, $2, false)",
		userID, types.StatusStudentNotChecked)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) AddMarkersForText(markers *types.Markers, userID, teacherID int) error {

	_, err := p.db.Exec("INSERT INTO markers (student_id, teacher_id, text, position) VALUES ($1, $2, $3, $4)",
		userID, teacherID, markers.Text, markers.Position)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetMarkers(userID int) ([]types.Markers, error) {

	markers := make([]types.Markers, 0)
	rows, err := p.db.Query("SELECT text, position FROM markers WHERE student_id = $1", userID)
	if err != nil {
		return nil, errors.Wrap(err, "err with Query")
	}
	defer rows.Close()
	marker := types.Markers{}
	for rows.Next() {
		if err = rows.Scan(&marker.Text, &marker.Position); err != nil {
			return nil, errors.Wrap(err, "err with Scan")
		}
		markers = append(markers, marker)
	}

	return markers, nil
}

func (p *Postgres) ChangeSendCertStatus(userID int) error {

	_, err := p.db.Exec("UPDATE dictation SET send_cert = true WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetStatusDicByUserID(userID int) (string, error) {

	var status string
	err := p.db.QueryRow("SELECT status FROM dictation WHERE user_id = $1", userID).
		Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}
