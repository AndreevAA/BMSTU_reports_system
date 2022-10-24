package mapper

import (
	authMapper "reports_system/internal/mapper/account"
	labelMapper "reports_system/internal/mapper/label"
	reportMapper "reports_system/internal/mapper/report"
	"reports_system/internal/model/account"
	"reports_system/internal/model/label"
	"reports_system/internal/model/report"
	"reports_system/pkg/logging"
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
