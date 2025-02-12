package config

import (
	"fmt"

	"github.com/Pharmeum/pharmeum-users-api/db"
	"github.com/caarlos0/env"
)

type Database struct {
	Name     string `env:"PHARMEUM_DATABASE_NAME,required"`
	Host     string `env:"PHARMEUM_DATABASE_HOST,required"`
	Port     int    `env:"PHARMEUM_DATABASE_PORT,required"`
	User     string `env:"PHARMEUM_DATABASE_USER,required"`
	Password string `env:"PHARMEUM_DATABASE_PASSWORD,required"`
	SSL      string `env:"PHARMEUM_DATABASE_SSL,required"`
}

func (d Database) URL() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s", d.Host, d.Port, d.User, d.Password, d.Name, d.SSL)
}

func (c *ConfigImpl) DB() *db.DB {
	if c.db != nil {
		return c.db
	}

	c.Lock()
	defer c.Unlock()

	var database Database
	if err := env.Parse(&database); err != nil {
		panic(err)
	}

	dbInstance, err := db.New(database.URL())
	if err != nil {
		panic(err)
	}

	c.db = dbInstance

	return c.db
}
