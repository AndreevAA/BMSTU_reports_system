package report

import (
	"reports_system/internal/model/report"
	"reports_system/pkg/logging"
)

type mapper struct {
	logger logging.Logger
}

func New(logger logging.Logger) *mapper {
	return &mapper{logger: logger}
}

func (m *mapper) MapCreateReportDTO(dto report.CreateReportDTO) report.Report {
	n := report.Report{
		ID:        0,
		Header:    dto.Header,
		Body:      dto.Body,
		ShortBody: "",
		Labels:    nil,
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

	n := report.Report{
		ID:        dto.ID,
		Header:    dto.Header,
		Body:      dto.Body,
		ShortBody: "",
	}

	n.GenerateShortBody()
	m.logger.Infof("Generated short body with length of %v symbols", len(n.ShortBody))

	return n
}
