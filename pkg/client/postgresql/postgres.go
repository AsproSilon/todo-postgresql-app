package postgresql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	UsersTable        = "users"
	PersonTable       = "person"
	DeviceTable       = "device"
	UsersPersonTable  = "users_person"
	PersonDeviceTable = "person_device"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

func NewConfig(cfg *Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
