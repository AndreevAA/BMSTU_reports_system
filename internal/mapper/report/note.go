package report

import (
	"neatly/internal/model/report"
	"neatly/pkg/logging"
)

type mapper struct {
	logger logging.Logger
}

func New(logger logging.Logger) *mapper {
	return &mapper{logger: logger}
}

func (m *mapper) MapCreateReportDTO(dto report.CreateReportDTO) report.Report {
	n := report.Report{
		ID:         0,
		Header:     dto.Header,
		Body:       dto.Body,
		ShortBody:  "",
		Labels:     nil,
		Department: dto.Department,
	}

	n.GenerateShortBody()
	m.logger.Infof("Generated short body with length of %v symbols", len(n.ShortBody))

	return n
}

func (m *mapper) MapGetAllReportsDTO(ns []report.Report) report.GetAllReportsDTO {
	return report.GetAllReportsDTO{
		Reports: ns,
	}
}

func (m *mapper) MapUpdateReportDTO(dto report.UpdateReportDTO) report.Report {
	if dto.Department == "" {
		dto.Department = report.DefaultReportDepartment
	}

	n := report.Report{
		ID:         dto.ID,
		Header:     dto.Header,
		Body:       dto.Body,
		ShortBody:  "",
		Department: dto.Department,
	}

	n.GenerateShortBody()
	m.logger.Infof("Generated short body with length of %v symbols", len(n.ShortBody))

	return n
}
