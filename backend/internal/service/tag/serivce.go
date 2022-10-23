package tag

import (
	"errors"
	"neatly/internal/model/tag"
	"neatly/internal/repository"
	"neatly/pkg/logging"
	"strings"
)

type Service struct {
	tagsRepository  repository.Tag
	reportsRepository repository.Report
	logger          logging.Logger
}

func NewService(tr repository.Tag, nr repository.Report, l logging.Logger) *Service {
	return &Service{tagsRepository: tr, reportsRepository: nr, logger: l}
}

func (s *Service) Create(userID, reportID int, t *tag.Tag) error {
	_, err := s.reportsRepository.GetOne(userID, reportID)
	if err != nil {
		return errors.New("report does not exists or does not belong to accounts")
	}

	tags, err := s.tagsRepository.GetAll(userID)
	if err != nil {
		return err
	}

	unique, tuID := s.checkIfUnique(tags, *t)
	if !unique {
		s.logger.Infof("Tag with ID %v is not unique", tuID)
		assigned, err := s.checkIfAssigned(tuID, reportID, userID)
		if err != nil {
			return err
		}
		if !assigned {
			s.logger.Infof("Tag with ID %v is not assigned to report %v", tuID, reportID)
			t.ID = tuID
			err := s.tagsRepository.Assign(tuID, reportID, userID)
			return err
		}
		t.ID = tuID
		return nil
	}

	s.logger.Infof("Tag with ID %v is inuque and will be assigned to report with ID %v", t.ID, reportID)
	err = s.tagsRepository.Create(userID, reportID, t)
	if err != nil {
		return err
	}
	err = s.tagsRepository.Assign(t.ID, reportID, userID)
	return err
}

func (s *Service) GetAll(userID int) ([]tag.Tag, error) {
	return s.tagsRepository.GetAll(userID)
}

func (s *Service) GetAllByReport(userID, reportID int) ([]tag.Tag, error) {
	return s.tagsRepository.GetAllByReport(userID, reportID)
}

func (s *Service) GetOne(userID, tagID int) (tag.Tag, error) {
	return s.tagsRepository.GetOne(userID, tagID)
}

func (s *Service) Delete(userID, tagID int) error {
	return s.tagsRepository.Delete(userID, tagID)
}

func (s *Service) Update(userID, tagID int, t tag.Tag) error {
	tp, err := s.tagsRepository.GetOne(userID, tagID)
	if err != nil {
		return err
	}

	if t.Department == "" {
		t.Department = tp.Department
	}
	if t.Name == "" {
		t.Name = tp.Name
	}

	return s.tagsRepository.Update(userID, tagID, t)
}

func (s *Service) Detach(userID, tagID, reportID int) error {
	_, err := s.reportsRepository.GetOne(userID, reportID)
	if err != nil {
		return errors.New("report does not exists or does not belong to accounts")
	}

	ns, err := s.reportsRepository.GetAll(userID)
	if err != nil {
		return err
	}

	t, err := s.tagsRepository.GetOne(userID, tagID)
	s.logger.Infof("Found tag %v: %v, %v", tagID, t.Name, t.Department)
	if err != nil {
		return err
	}

	for _, n := range ns {
		n.Tags, err = s.tagsRepository.GetAllByReport(userID, n.ID)
		if err != nil {
			return err
		}
		if n.HasSpecificTag(t.Name) && n.ID != reportID {
			s.logger.Infof("Found this tag at report %v", n.ID)
			err = s.tagsRepository.Detach(userID, tagID, reportID)
			return err
		}
	}

	s.logger.Info("Deleting tag")
	err = s.tagsRepository.Delete(userID, tagID)
	return err
}

func (s *Service) checkIfUnique(tags []tag.Tag, tu tag.Tag) (bool, int) {
	for _, t := range tags {
		if strings.Compare(t.Name, tu.Name) == 0 {
			s.logger.Infof("Found matching tag with id %v", t.ID)
			return false, t.ID
		}
	}
	return true, 0
}

func (s *Service) checkIfAssigned(tagID, reportID, userID int) (bool, error) {
	tags, err := s.tagsRepository.GetAllByReport(userID, reportID)
	if err != nil {
		return false, err
	}

	for _, t := range tags {
		if t.ID == tagID {
			s.logger.Infof("Found matching tag %v assigned to report %v", tagID, reportID)
			return true, nil
		}
	}
	return false, nil
}
