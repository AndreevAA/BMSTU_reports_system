package repository

import (
	"reports_system/internal/model/account"
	"reports_system/internal/model/label"
	"reports_system/internal/model/report"
	"reports_system/internal/repository/psql"
	"reports_system/pkg/client/psqlclient"
	"reports_system/pkg/logging"
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

type Label interface {
	Create(userID int, reportID int, t *label.Label) error
	GetAll(userID int) ([]label.Label, error)
	GetAllByReport(userID, reportID int) ([]label.Label, error)
	GetOne(userID, labelID int) (label.Label, error)
	Delete(userID, labelID int) error
	Detach(userID, labelID, reportID int) error
	Assign(labelID, reportID, userID int) error
	Update(userID, labelID int, t label.Label) error
}

type Repository struct {
	Account
	Report
	Label
}

func New(client *psqlclient.Client, logger logging.Logger) *Repository {
	return &Repository{
		Account: psql.NewAuthPostgres(client, logger),
		Report:  psql.NewReportPostgres(client, logger),
		Label:   psql.NewLabelPostgres(client, logger),
	}
}
