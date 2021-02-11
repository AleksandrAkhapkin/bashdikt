package mail

import (
	"fmt"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types/config"
	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type Mail struct {
	email      *config.ForSendEmail
	serverPort string
}

func NewMail(cnf *config.ForSendEmail, port string) (*Mail, error) {
	return &Mail{
		email:      cnf,
		serverPort: port,
	}, nil
}

func (r *Mail) SendMail(template string, data, dataHTML string, to string) error {

	m := gomail.NewMessage()
	m.SetAddressHeader("From", r.email.EmailLogin, r.email.NameSender)
	m.SetAddressHeader("To", to, to)
	m.SetHeader("From", fmt.Sprintf("%s <%s>", r.email.NameSender, r.email.EmailLogin))
	m.SetHeader("To", to)
	m.SetHeader("Subject", template)
	m.SetHeader("MIME-Version:", "1.0")
	m.SetHeader("Reply-To", r.email.EmailLogin)
	//m.SetHeader("List-Unsubscribe", fmt.Sprintf("<mailto: %s>, <https://supreme-cheese.ru%s/mail/unsubscribe/%s>", r.serverPort, r.email.EmailUnsubscribe, to)) //r.to[0])) //, <https://supreme-cheese.ru:8080/mail/unsubscribe/%s>
	//m.SetHeader("List-Unsubscribe-Post", "List-Unsubscribe=One-Click")

	m.SetBody("text/plain", data)
	port, err := strconv.Atoi(r.email.EmailPort)
	if err != nil {
		return errors.Wrap(err, "err while Atoi ")
	}
	d := gomail.NewDialer(r.email.EmailHost, port, r.email.EmailLogin, r.email.EmailPass)

	stopMail := 0
	for stopMail < 5 {
		stopMail++

		if err := d.DialAndSend(m); err != nil {

			time.Sleep(time.Second * 30)
			if stopMail < 5 {
				continue
			}
			return errors.Wrap(err, fmt.Sprintf("err with DialAndSend in FOR USER_email = %s", to))
		}
		break
	}

	return nil
}

func (r *Mail) SendCert(pathForCert string, to string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", r.email.EmailLogin, r.email.NameSender)
	m.SetAddressHeader("To", to, to)
	m.SetHeader("From", r.email.EmailLogin)
	m.SetHeader("To", to)
	m.SetHeader("Subject", DictantTmplt)
	m.SetBody("text/plain", BodyCertificateInAttachText)
	m.Attach(pathForCert)
	port, err := strconv.Atoi(r.email.EmailPort)
	if err != nil {
		return err
	}
	d := gomail.NewDialer(r.email.EmailHost, port, r.email.EmailLogin, r.email.EmailPass)

	stopMail := 0
	for stopMail < 5 {
		stopMail++

		if err := d.DialAndSend(m); err != nil {

			time.Sleep(time.Second * 30)
			if stopMail < 5 {
				continue
			}
			return errors.Wrap(err, fmt.Sprintf("err with DialAndSend in SendCert USER_email = %s", to))
		}
		break
	}

	return nil
}
