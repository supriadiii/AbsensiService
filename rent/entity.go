package rent

import "project_absensi/common"

type Rent struct {
	common.Model
	UserID          uint
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
