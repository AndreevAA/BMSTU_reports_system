package label

import (
	"errors"
	"neatly/internal/model/label"
	"neatly/internal/repository"
	"neatly/pkg/logging"
	"strings"
)

type Service struct {
	labelsRepository  repository.Label
	reportsRepository repository.Report
	logger            logging.Logger
}

func NewService(tr repository.Label, nr repository.Report, l logging.Logger) *Service {
	return &Service{labelsRepository: tr, reportsRepository: nr, logger: l}
}

func (s *Service) Create(userID, reportID int, t *label.Label) error {
	_, err := s.reportsRepository.GetOne(userID, reportID)
	if err != nil {
		return errors.New("report does not exists or does not belong to accounts")
	}

	labels, err := s.labelsRepository.GetAll(userID)
	if err != nil {
		return err
	}

	unique, tuID := s.checkIfUnique(labels, *t)
	if !unique {
		s.logger.Infof("Label with ID %v is not unique", tuID)
		assigned, err := s.checkIfAssigned(tuID, reportID, userID)
		if err != nil {
			return err
		}
		if !assigned {
			s.logger.Infof("Label with ID %v is not assigned to report %v", tuID, reportID)
			t.ID = tuID
			err := s.labelsRepository.Assign(tuID, reportID, userID)
			return err
		}
		t.ID = tuID
		return nil
	}

	s.logger.Infof("Label with ID %v is inuque and will be assigned to report with ID %v", t.ID, reportID)
	err = s.labelsRepository.Create(userID, reportID, t)
	if err != nil {
		return err
	}
	err = s.labelsRepository.Assign(t.ID, reportID, userID)
	return err
}

func (s *Service) GetAll(userID int) ([]label.Label, error) {
	return s.labelsRepository.GetAll(userID)
}

func (s *Service) GetAllByReport(userID, reportID int) ([]label.Label, error) {
	return s.labelsRepository.GetAllByReport(userID, reportID)
}

func (s *Service) GetOne(userID, labelID int) (label.Label, error) {
	return s.labelsRepository.GetOne(userID, labelID)
}

func (s *Service) Delete(userID, labelID int) error {
	return s.labelsRepository.Delete(userID, labelID)
}

func (s *Service) Update(userID, labelID int, t label.Label) error {
	tp, err := s.labelsRepository.GetOne(userID, labelID)
	if err != nil {
		return err
	}

	if t.Department == "" {
		t.Department = tp.Department
	}
	if t.Name == "" {
		t.Name = tp.Name
	}

	return s.labelsRepository.Update(userID, labelID, t)
}

func (s *Service) Detach(userID, labelID, reportID int) error {
	_, err := s.reportsRepository.GetOne(userID, reportID)
	if err != nil {
		return errors.New("report does not exists or does not belong to accounts")
	}

	ns, err := s.reportsRepository.GetAll(userID)
	if err != nil {
		return err
	}

	t, err := s.labelsRepository.GetOne(userID, labelID)
	s.logger.Infof("Found label %v: %v, %v", labelID, t.Name, t.Department)
	if err != nil {
		return err
	}

	for _, n := range ns {
		n.Labels, err = s.labelsRepository.GetAllByReport(userID, n.ID)
		if err != nil {
			return err
		}
		if n.HasSpecificLabel(t.Name) && n.ID != reportID {
			s.logger.Infof("Found this label at report %v", n.ID)
			err = s.labelsRepository.Detach(userID, labelID, reportID)
			return err
		}
	}

	s.logger.Info("Deleting label")
	err = s.labelsRepository.Delete(userID, labelID)
	return err
}

func (s *Service) checkIfUnique(labels []label.Label, tu label.Label) (bool, int) {
	for _, t := range labels {
		if strings.Compare(t.Name, tu.Name) == 0 {
			s.logger.Infof("Found matching label with id %v", t.ID)
			return false, t.ID
		}
	}
	return true, 0
}

func (s *Service) checkIfAssigned(labelID, reportID, userID int) (bool, error) {
	labels, err := s.labelsRepository.GetAllByReport(userID, reportID)
	if err != nil {
		return false, err
	}

	for _, t := range labels {
		if t.ID == labelID {
			s.logger.Infof("Found matching label %v assigned to report %v", labelID, reportID)
			return true, nil
		}
	}
	return false, nil
}
