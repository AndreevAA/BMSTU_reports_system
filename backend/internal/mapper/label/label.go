package label

import (
	"reports_system/internal/model/label"
	"reports_system/pkg/logging"
)

type mapper struct {
	logger logging.Logger
}

func New(logger logging.Logger) *mapper {
	return &mapper{logger: logger}
}

func (m *mapper) MapCreateLabelDTO(dto label.CreateLabelDTO) label.Label {
	return label.Label{
		ID:         0,
		Name:       dto.Name,
		Department: dto.Department,
	}
}

func (m *mapper) MapUpdateLabelDTO(dto label.UpdateLabelDTO) label.Label {
	return label.Label{
		ID:         0,
		Name:       dto.Name,
		Department: dto.Department,
	}
}

func (m *mapper) MapGetAllLabelsDTO(labels []label.Label) label.GetAllLabelsDTO {
	return label.GetAllLabelsDTO{
		Labels: labels,
	}
}
