package email

import (
	"github.com/go-gomail/gomail"
	"github.com/stellar/go/support/errors"
)

type Client interface {
	Signup(to, link string) error
	Forgot(to, link string) error
	NewPassword(to string) error
}

type ClientImpl struct {
	emailAddress string
	password     string
	host         string
	port         int
}

func New(emailAddress, password, host string, port int) *ClientImpl {
	return &ClientImpl{
		emailAddress: emailAddress,
		password:     password,
		host:         host,
		port:         port,
	}
}

func (c ClientImpl) Signup(to, link string) error {
	dialer := gomail.NewPlainDialer(c.host, c.port, c.emailAddress, c.password)
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", c.emailAddress, "Pharmeum")
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Account Verification Pharmeum")
	msg.SetBody("text/html", "To verify your account, please click on the link: "+link)
	if err := dialer.DialAndSend(msg); err != nil {
		return errors.Wrap(err, "failed to send sign up verification email")
	}

	return nil
}

func (c ClientImpl) Forgot(to, link string) error {
	dialer := gomail.NewPlainDialer(c.host, c.port, c.emailAddress, c.password)
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", c.emailAddress, "Pharmeum")
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Account Verification Pharmeum")
	msg.SetBody("text/html", "To change your password, please click on the link: <a href=\""+link+
		"\">"+link+"</a><br><br>Best Regards,<br>Pharmeum")

	//TODO wait until Zain provide valid HTML template
	if err := dialer.DialAndSend(msg); err != nil {
		return errors.Wrap(err, "failed to send forgot password email")
	}

	return nil
}

func (c ClientImpl) NewPassword(to string) error {
	dialer := gomail.NewPlainDialer(c.host, c.port, c.emailAddress, c.password)
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", c.emailAddress, "Pharmeum")
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Pharmeum password changed")
	msg.SetBody("text/html", "Your password was successfully changed")

	//TODO wait until Zain provide valid HTML template
	if err := dialer.DialAndSend(msg); err != nil {
		return errors.Wrap(err, "failed to send new password request")
	}

	return nil
}
