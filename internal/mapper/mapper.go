package mapper

import (
	authMapper "neatly/internal/mapper/account"
	reportMapper "neatly/internal/mapper/report"
	tagMapper "neatly/internal/mapper/tag"
	"neatly/internal/model/account"
	"neatly/internal/model/report"
	"neatly/internal/model/tag"
	"neatly/pkg/logging"
)

type Account interface {
	MapRegisterAccountDTO(dto account.RegisterAccountDTO) (account.Account, error)
	MapLogInAccountDTO(dto account.LoginAccountDTO) account.Account
	MapAccountWithTokenDTO(token string, a account.Account) account.WithTokenDTO
	MapAccountDTO(a account.Account) account.GetAccountDTO
}

type Report interface {
	MapCreateReportDTO(dto report.CreateReportDTO) report.Report
	MapUpdateReportDTO(dto report.UpdateReportDTO) report.Report
	MapGetAllReportsDTO(ns []report.Report) report.GetAllReportsDTO
}

type Tag interface {
	MapCreateTagDTO(dto tag.CreateTagDTO) tag.Tag
	MapUpdateTagDTO(dto tag.UpdateTagDTO) tag.Tag
	MapGetAllTagsDTO(tags []tag.Tag) tag.GetAllTagsDTO
}

type Mapper struct {
	Account
	Report
	Tag
}

func New(l logging.Logger) *Mapper {
	return &Mapper{
		Account: authMapper.New(l),
		Report:    reportMapper.New(l),
		Tag:     tagMapper.New(l),
	}
}
