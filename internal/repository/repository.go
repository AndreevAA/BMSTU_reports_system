package repository

import (
	"neatly/internal/model/account"
	"neatly/internal/model/report"
	"neatly/internal/model/tag"
	"neatly/internal/repository/psql"
	"neatly/pkg/client/psqlclient"
	"neatly/pkg/logging"
)

type Account interface {
	CreateAccount(a *account.Account) error
	AuthorizeAccount(a *account.Account) error
	GetOne(userID int) (account.Account, error)
}

type Report interface {
	Create(userID int, report *report.Report) error
	GetAll(userID int) ([]report.Report, error)
	GetOne(userID, reportID int) (report.Report, error)
	Delete(userID, reportID int) error
	Update(userID int, n report.Report) error
}

type Tag interface {
	Create(userID int, reportID int, t *tag.Tag) error
	GetAll(userID int) ([]tag.Tag, error)
	GetAllByReport(userID, reportID int) ([]tag.Tag, error)
	GetOne(userID, tagID int) (tag.Tag, error)
	Delete(userID, tagID int) error
	Detach(userID, tagID, reportID int) error
	Assign(tagID, reportID, userID int) error
	Update(userID, tagID int, t tag.Tag) error
}

type Repository struct {
	Account
	Report
	Tag
}

func New(client *psqlclient.Client, logger logging.Logger) *Repository {
	return &Repository{
		Account: psql.NewAuthPostgres(client, logger),
		Report:    psql.NewReportPostgres(client, logger),
		Tag:     psql.NewTagPostgres(client, logger),
	}
}
