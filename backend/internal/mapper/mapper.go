package mapper

import (
	authMapper "neatly/internal/mapper/account"
	labelMapper "neatly/internal/mapper/label"
	reportMapper "neatly/internal/mapper/report"
	"neatly/internal/model/account"
	"neatly/internal/model/label"
	"neatly/internal/model/report"
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

type Label interface {
	MapCreateLabelDTO(dto label.CreateLabelDTO) label.Label
	MapUpdateLabelDTO(dto label.UpdateLabelDTO) label.Label
	MapGetAllLabelsDTO(labels []label.Label) label.GetAllLabelsDTO
}

type Mapper struct {
	Account
	Report
	Label
}

func New(l logging.Logger) *Mapper {
	return &Mapper{
		Account: authMapper.New(l),
		Report:  reportMapper.New(l),
		Label:   labelMapper.New(l),
	}
}
