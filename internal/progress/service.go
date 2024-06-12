package progress

import (
	"errors"
	"learn-swiping-api/erro"
	progress "learn-swiping-api/internal/progress/dto"
)

type ProgressService interface {
	Create(progress.AccessRequest) error
	Progress(progress.AccessRequest) (Progress, error)
	Update(progress.UpdateRequest) error
	Delete(progress.AccessRequest) error
}

type ProgressServiceImpl struct {
	repository ProgressRepository
}

func NewProgressService(repository ProgressRepository) ProgressService {
	return &ProgressServiceImpl{repository: repository}
}

func (s *ProgressServiceImpl) Create(req progress.AccessRequest) error {
	_, err := s.repository.Create(req)
	return err
}

func (s *ProgressServiceImpl) Progress(req progress.AccessRequest) (Progress, error) {
	return s.repository.Progress(req)
}

func (s *ProgressServiceImpl) Update(req progress.UpdateRequest) error {
	if req.Ease == nil && req.Interval == nil && req.Priority == nil && req.DaysHidden == nil && req.WatchCount == nil && req.PriorityExam == nil && req.DaysHiddenExam == nil && req.AnswerCount == nil && req.CorrectCount == nil && req.IsRelearning == nil && req.IsBuried == nil {
		return erro.ErrBadField
	}

	accReq := progress.AccessRequest{Token: req.Token, CardID: req.CardID}
	_, err := s.repository.Progress(accReq)
	if err != nil {
		if errors.Is(err, erro.ErrProgressNotFound) {
			// If progress doesn't exist, then creates and later updates it.
			err = s.Create(accReq)
			if err != nil {
				return err
			}
		}
	}

	err = s.repository.Update(req)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProgressServiceImpl) Delete(req progress.AccessRequest) error {
	return s.repository.Delete(req)
}
