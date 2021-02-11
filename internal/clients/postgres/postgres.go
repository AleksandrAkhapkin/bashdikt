package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(dsn string) (*Postgres, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "err with Open DB")
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "err with ping DB")
	}

	return &Postgres{db}, nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func (p *Postgres) SaveRequestLog(body, route, ip string) error {

	_, err := p.db.Exec("INSERT INTO request_log (dump_request, ip, route) VALUES ($1, $2, $3)", body, ip, route)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) AddWhiteEmail(email string) error {

	var err error
	_, err = p.db.Exec("INSERT INTO white_emails (email) VALUES ($1)", email)
	if err != nil {
		return err
	}

	return nil
}
