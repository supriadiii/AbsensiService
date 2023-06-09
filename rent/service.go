package rent

import (
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetRents(userID uint) ([]Rent, error)
	CreateRent(input CreateRentInput) (Rent, error)
	SaveRentImage(input CreateRentImageInput, filelocation string) (RentImage, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetRents(userID uint) ([]Rent, error) {
	var rent []Rent
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

func (s *service) CreateRent(input CreateRentInput) (Rent, error) {
	rent := Rent{}
	rent.Name = input.Name
	rent.SortDescription = input.SortDescription
	rent.Description = input.Description
	rent.ContactPerson = input.ContactPerson
	rent.Price = input.Price
	rent.Quantity = input.Quantity
	rent.UserID = input.User.ID

	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	rent.Slug = slug.Make(slugCandidate)

	newRent, err := s.repository.Save(rent)
	if err != nil {
		return newRent, err
	}
	return newRent, nil
}

func (s *service) SaveRentImage(input CreateRentImageInput, filelocation string) (RentImage, error) {
	isPrimary := 0
	if input.IsPrimary {
		isPrimary = 1
		_, err := s.repository.MarkAllImagesAsNonPrimary(input.RentID)
		if err != nil {
			return RentImage{}, err
		}

	}

	rentImage := RentImage{}
	rentImage.RentID = input.RentID
	rentImage.IsPrimary = isPrimary
	rentImage.FileName = filelocation

	newRentImage, err := s.repository.CreateImage(rentImage)
	if err != nil {
		return newRentImage, err
	}
	return newRentImage, nil
}
