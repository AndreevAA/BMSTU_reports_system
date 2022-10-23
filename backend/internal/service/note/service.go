package report

import (
	"neatly/internal/model/report"
	"neatly/internal/repository"
	"neatly/pkg/logging"
)

type Service struct {
	reportsRepository repository.Report
	tagsRepository  repository.Tag
	logger          logging.Logger
}

func NewService(reportsRepository repository.Report, tagsRepository repository.Tag, logger logging.Logger) *Service {
	return &Service{reportsRepository: reportsRepository, tagsRepository: tagsRepository, logger: logger}
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
		tags, err := s.tagsRepository.GetAllByReport(userID, reportID)
		if err != nil {
			return []report.Report{}, err
		}
		reports[i].Tags = tags
	}

	return reports, nil
}

func (s *Service) GetOne(userID, reportID int) (report.Report, error) {
	n, err := s.reportsRepository.GetOne(userID, reportID)
	if err != nil {
		return n, err
	}

	tags, err := s.tagsRepository.GetAllByReport(userID, n.ID)
	if err != nil {
		return n, err
	}
	n.Tags = tags

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

func (s *Service) FindByTags(userID int, tagNames []string) ([]report.Report, error) {
	ns, err := s.reportsRepository.GetAll(userID)
	if err != nil {
		return ns, err
	}

	var (
		reportsWithAllTags []report.Report
	)

	for _, n := range ns {
		n.Tags, err = s.tagsRepository.GetAllByReport(userID, n.ID)
		if err != nil {
			return ns, err
		}

		s.logger.Infof("Found tags from report %v: %v", n.ID, n.Tags)
		if n.HasEveryTag(tagNames) {
			reportsWithAllTags = append(reportsWithAllTags, n)
		}
	}

	return reportsWithAllTags, nil
}
