package rent

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Rent, error)
	FindByUserID(UserID uint) ([]Rent, error)
	Save(Rent Rent) (Rent, error)
	CreateImage(rentImage RentImage) (RentImage, error)
	MarkAllImagesAsNonPrimary(rentID uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Rent, error) {
	var Rent []Rent
	err := r.db.Preload("User").Preload("RentImage", "rent_image.is_primary=1").Find(&Rent).Error
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
	err := r.db.Create(&Rent).Preload("User").Error
	if err != nil {
		return Rent, err
	}
	return Rent, nil
}

func (r *repository) CreateImage(rentImage RentImage) (RentImage, error) {
	err := r.db.Create(&rentImage).Preload("User").Error
	if err != nil {
		return rentImage, err
	}
	return rentImage, nil
}

func (r *repository) MarkAllImagesAsNonPrimary(rentID uint) (bool, error) {
	err := r.db.Model(&RentImage{}).Where("rent_id=?", rentID).Update("is_primary", false).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
