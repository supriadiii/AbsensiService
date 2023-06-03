package user

import "gorm.io/gorm"

type Repository interface {
	Save(User User) (User, error)
	GetAllUsers(ID uint) ([]User, error)
	FindByNim(nim string) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(User User) (User, error) {
	err := r.db.Create(&User).Error
	if err != nil {
		return User, err
	}
	return User, nil
}

func (r *repository) GetAllUsers(ID uint) ([]User, error) {
	var Users []User
	err := r.db.Limit(10).Where("id>= ?", ID).Find(&Users).Error
	if err != nil {
		return nil, err
	}
	return Users, nil
}

func (r *repository) FindByNim(nim string) (User, error) {
	var user User
	if err := r.db.Where("nim = ?", nim).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
