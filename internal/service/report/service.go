package report

import (
	"neatly/internal/model/report"
	"neatly/internal/repository"
	"neatly/pkg/logging"
)

type Service struct {
	reportsRepository repository.Report
	labelsRepository  repository.Label
	logger            logging.Logger
}

func NewService(reportsRepository repository.Report, labelsRepository repository.Label, logger logging.Logger) *Service {
	return &Service{reportsRepository: reportsRepository, labelsRepository: labelsRepository, logger: logger}
}

func (s *Service) Create(userID int, n *report.Report) error {
	err := s.reportsRepository.Create(userID, n)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetAll(userID int) ([]report.Report, error) {
	reports, err := s.reportsRepository.GetAll(userID)
	if err != nil {
		return []report.Report{}, err
	}

	for i := 0; i < len(reports); i++ {
		reportID := reports[i].ID
		labels, err := s.labelsRepository.GetAllByReport(userID, reportID)
		if err != nil {
			return []report.Report{}, err
		}
		reports[i].Labels = labels
	}

	return reports, nil
}

func (s *Service) GetOne(userID, reportID int) (report.Report, error) {
	n, err := s.reportsRepository.GetOne(userID, reportID)
	if err != nil {
		return n, err
	}

	labels, err := s.labelsRepository.GetAllByReport(userID, n.ID)
	if err != nil {
		return n, err
	}
	n.Labels = labels

	return n, nil
}

func (s *Service) Delete(userID, reportID int) error {
	return s.reportsRepository.Delete(userID, reportID)
}

func (s *Service) Update(userID int, n report.Report, needBodyUpdate bool) error {
	prev, err := s.reportsRepository.GetOne(userID, n.ID)
	if err != nil {
		return err
	}
	if n.Header == "" {
		n.Header = prev.Header
	}

	if !needBodyUpdate {
		n.Body = prev.Body
		n.ShortBody = prev.ShortBody
	}

	return s.reportsRepository.Update(userID, n)
}

func (s *Service) FindByLabels(userID int, labelNames []string) ([]report.Report, error) {
	ns, err := s.reportsRepository.GetAll(userID)
	if err != nil {
		return ns, err
	}

	var (
		reportsWithAllLabels []report.Report
	)

	for _, n := range ns {
		n.Labels, err = s.labelsRepository.GetAllByReport(userID, n.ID)
		if err != nil {
			return ns, err
		}

		s.logger.Infof("Found labels from report %v: %v", n.ID, n.Labels)
		if n.HasEveryLabel(labelNames) {
			reportsWithAllLabels = append(reportsWithAllLabels, n)
		}
	}

	return reportsWithAllLabels, nil
}
