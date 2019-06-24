package email

import (
	"github.com/alexcesaro/mail/gomail"
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
	from := c.emailAddress
	password := c.password
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", from, "Pharmeum")
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Account Verification Pharmeum")
	msg.SetBody("text/html", "To verify your account, please click on the link: "+link)

	m := gomail.NewMailer(c.host, from, password, c.port)
	if err := m.Send(msg); err != nil {
		return err
	}

	return nil
}

func (c ClientImpl) Forgot(to, link string) error {
	from := c.emailAddress
	password := c.password
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", from, "Pharmeum")
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Account Verification Pharmeum")
	msg.SetBody("text/html", "To change your password, please click on the link: <a href=\""+link+
		"\">"+link+"</a><br><br>Best Regards,<br>Pharmeum")

	m := gomail.NewMailer(c.host, from, password, c.port)
	if err := m.Send(msg); err != nil {
		return err
	}

	return nil
}

func (c ClientImpl) NewPassword(to string) error {
	from := c.emailAddress
	password := c.password
	msg := gomail.NewMessage()

	msg.SetAddressHeader("From", from, "Pharmeum")
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Pharmeum password changed")
	msg.SetBody("text/html", "Your password was successfully changed")

	m := gomail.NewMailer(c.host, from, password, c.port)
	if err := m.Send(msg); err != nil {
		return err
	}

	return nil
}
