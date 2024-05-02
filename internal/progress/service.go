package progress

import (
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
	if req.Priority == nil && req.DaysHidden == nil && req.WatchCount == nil && req.PriorityExam == nil && req.DaysHiddenExam == nil && req.AnswerCount == nil && req.CorrectCount == nil && req.IsBuried == nil {
		return erro.ErrBadField
	}

	return s.repository.Update(req)
}

func (s *ProgressServiceImpl) Delete(req progress.AccessRequest) error {
	return s.repository.Delete(req)
}
