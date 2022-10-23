package service

import (
	"neatly/internal/model/account"
	"neatly/internal/model/label"
	"neatly/internal/model/report"
	"neatly/internal/repository"
	authService "neatly/internal/service/account"
	labelService "neatly/internal/service/label"
	reportService "neatly/internal/service/report"
	"neatly/pkg/logging"
)

type Account interface {
	CreateAccount(u *account.Account) error
	GenerateJWT(u *account.Account) (string, error)
	GetOne(userID int) (account.Account, error)
}

type Report interface {
	Create(userID int, n *report.Report) error
	GetAll(userID int) ([]report.Report, error)
	GetOne(userID, reportID int) (report.Report, error)
	Delete(userID, reportID int) error
	Update(userID int, n report.Report, needBodyUpdate bool) error
	FindByLabels(userID int, labelNames []string) ([]report.Report, error)
}

type Label interface {
	Create(userID, reportID int, label *label.Label) error
	GetAll(userID int) ([]label.Label, error)
	GetAllByReport(userID, reportID int) ([]label.Label, error)
	GetOne(userID, labelID int) (label.Label, error)
	Delete(userID, labelID int) error
	Update(userID, labelID int, t label.Label) error
	Detach(userID, labelID, reportID int) error
}

type Service struct {
	Account
	Report
	Label
}

func New(repo *repository.Repository, logger logging.Logger) *Service {
	return &Service{
		Account: authService.NewService(repo.Account),
		Report:  reportService.NewService(repo.Report, repo.Label, logger),
		Label:   labelService.NewService(repo.Label, repo.Report, logger),
	}
}
