package psqlclient

import (
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"reports_system/internal/session"
	"reports_system/pkg/logging"
)

type Client struct {
	DB *sqlx.DB
}

func NewClient(cfg session.DB) (*Client, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
		cfg.SSLMode,
	))
	if err != nil {
		logging.GetLogger().Info("Error while connecting to db")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Client{
		DB: db,
	}, nil
}
