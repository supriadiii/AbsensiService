package rent

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Rent, error)
	FindByUserID(UserID uint) ([]Rent, error)
	Save(Rent Rent) (Rent, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Rent, error) {
	var Rent []Rent
	err := r.db.Preload("RentImage", "rent_image.is_primary=1").Find(&Rent).Error
	if err != nil {
		return Rent, err
	}
	return Rent, nil
}

func (r *repository) FindByUserID(UserID uint) ([]Rent, error) {
	var rents []Rent
	err := r.db.Where("user_id = ?", UserID).Preload("RentImage", "rent_image.is_primary=1").Find(&rents).Error
	if err != nil {
		return rents, err
	}
	return rents, nil
}

func (r *repository) Save(Rent Rent) (Rent, error) {
	err := r.db.Create(&Rent).Error
	if err != nil {
		return Rent, err
	}
	return Rent, nil
}
