package service

import (
	"neatly/internal/model/account"
	"neatly/internal/model/report"
	"neatly/internal/model/tag"
	"neatly/internal/repository"
	authService "neatly/internal/service/account"
	reportService "neatly/internal/service/report"
	tagService "neatly/internal/service/tag"
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
	FindByTags(userID int, tagNames []string) ([]report.Report, error)
}

type Tag interface {
	Create(userID, reportID int, tag *tag.Tag) error
	GetAll(userID int) ([]tag.Tag, error)
	GetAllByReport(userID, reportID int) ([]tag.Tag, error)
	GetOne(userID, tagID int) (tag.Tag, error)
	Delete(userID, tagID int) error
	Update(userID, tagID int, t tag.Tag) error
	Detach(userID, tagID, reportID int) error
}

type Service struct {
	Account
	Report
	Tag
}

func New(repo *repository.Repository, logger logging.Logger) *Service {
	return &Service{
		Account: authService.NewService(repo.Account),
		Report:    reportService.NewService(repo.Report, repo.Tag, logger),
		Tag:     tagService.NewService(repo.Tag, repo.Report, logger),
	}
}
