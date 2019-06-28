package config

import (
	"net/url"
	"sync"

	"github.com/Pharmeum/pharmeum-users-api/db"
	"github.com/Pharmeum/pharmeum-users-api/email"

	"github.com/go-chi/jwtauth"

	"github.com/sirupsen/logrus"
)

type Config interface {
	HTTP() *HTTP
	Log() *logrus.Entry
	EmailClient() *email.ClientImpl
	WebsiteURL() *url.URL
	DB() *db.DB
	JWT() *jwtauth.JWTAuth
}

type ConfigImpl struct {
	sync.Mutex

	//internal objects
	http   *HTTP
	log    *logrus.Entry
	email  *email.ClientImpl
	webApp *url.URL
	db     *db.DB
	jwt    *jwtauth.JWTAuth
}

func New() Config {
	return &ConfigImpl{
		Mutex: sync.Mutex{},
	}
}
