package config

import (
	"github.com/Pharmeum/pharmeum-users-api/email"
	"github.com/caarlos0/env"
)

type EmailClient struct {
	EmailAddress string `env:"PHARMEUM_EMAIL_ADDRESS,required"`
	Password     string `env:"PHARMEUM_EMAIL_PASSWORD,required"`
	Host         string `env:"PHARMEUM_SMTP_SERVER_HOST" envDefault:"smtp.gmail.com"`
	Port         int    `env:"PHARMEUM_SMTP_SERVER_PORT" envDefault:"465"`
}

func (c *ConfigImpl) EmailClient() *email.ClientImpl {
	if c.email != nil {
		return c.email
	}

	c.Lock()
	defer c.Unlock()

	emailClient := &EmailClient{}
	if err := env.Parse(emailClient); err != nil {
		panic(err)
	}

	c.email = email.New(
		emailClient.EmailAddress,
		emailClient.Password,
		emailClient.Host,
		emailClient.Port,
	)
	return c.email
}
