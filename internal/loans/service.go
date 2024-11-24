package loans

import "time"

func NewService(repo Repo) Service {
	return &implService{
		repo: repo,
	}
}

type implService struct {
	repo Repo
}

func (s *implService) TakeBook(authToken string, bookID string) error {
	return nil
}

func (s *implService) ReturnBook(authToken string, bookID string) error {
	return nil
}

func (s *implService) CountAvailableBook(authToken string, bookID string) (uint, error) {
	return 0, nil
}

func (s *implService) ListReservations(authToken string, at time.Time) ([]LentBook, error) {
	return nil, nil
}

func (s *implService) ListOverdue(authToken string, at time.Time) ([]LentBook, error) {
	return nil, nil
}
