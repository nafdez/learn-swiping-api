package progress

type ProgressService interface {
}

type ProgressServiceImpl struct {
	repository ProgressRepository
}

func NewProgressService(repository ProgressRepository) ProgressService {
	return &ProgressServiceImpl{repository: repository}
}
