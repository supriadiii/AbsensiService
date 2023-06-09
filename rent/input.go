package rent

import (
	"project_absensi/user"
)

type CreateRentInput struct {
	Name            string `json:"name" binding:"required"`
	SortDescription string `json:"sort_description" binding:"required"`
	Description     string `json:"description" binding:"required"`
	ContactPerson   string `json:"contact_person" binding:"required"`
	Price           int    `json:"price" binding:"required"`
	Quantity        int    `json:"quantity" binding:"required"`
	User            user.User
}

type CreateRentImageInput struct {
	RentID    uint `form:"id" binding:"required"`
	IsPrimary bool `form:"is_primary"`
}
