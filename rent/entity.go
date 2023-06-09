package rent

import (
	"project_absensi/common"
	"project_absensi/user"
)

type Rent struct {
	common.Model
	UserID          uint
	User            user.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name            string
	SortDescription string
	Description     string
	ContactPerson   string
	Price           int
	Quantity        int
	Slug            string
	RentImage       []RentImage
}

type RentImage struct {
	common.Model
	RentID    uint
	FileName  string
	IsPrimary int
}
