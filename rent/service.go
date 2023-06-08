package rent

type Service interface {
	FindRents(userID uint) ([]Rent, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) FindRents(userID uint) ([]Rent, error) {
	if userID != 0 {
		rent, err := s.repository.FindByUserID(userID)
		if err != nil {
			return rent, err
		}
		return rent, nil
	}
	rent, err := s.repository.FindAll()
	if err != nil {
		return rent, err
	}
	return rent, nil
}
