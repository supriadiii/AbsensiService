package rent

import "project_absensi/user"

type RentFormatter struct {
	ID              uint                      `json:"id"`
	Name            string                    `json:"name"`
	SortDescription string                    `json:"sort_direction"`
	Description     string                    `json:"description"`
	ContactPerson   string                    `json:"contact_person"`
	Price           int                       `json:"price"`
	Quantity        int                       `json:"quantity"`
	UserID          uint                      `json:"user_id"`
	ImageURL        string                    `json:"image_url"`
	User            user.UserNoTokenFormatter `json:"user"`
}

func FormatRentFormatter(rent Rent) RentFormatter {
	user := user.FormatUserWithoutToken(rent.User)
	formatter := RentFormatter{}
	formatter.ID = rent.ID
	formatter.Name = rent.Name
	formatter.SortDescription = rent.SortDescription
	formatter.Description = rent.Description
	formatter.ContactPerson = rent.ContactPerson
	formatter.Price = rent.Price
	formatter.Quantity = rent.Quantity
	formatter.UserID = rent.UserID
	formatter.ImageURL = ""
	formatter.User = user

	if len(rent.RentImage) > 0 {
		formatter.ImageURL = rent.RentImage[0].FileName
	}

	// Mengisi data user ke dalam formatter

	return formatter
}

func FormatRents(rents []Rent) []RentFormatter {
	if len(rents) == 0 {
		return []RentFormatter{}
	}
	var RentsFormatters []RentFormatter
	for _, rent := range rents {
		RentFormatter := FormatRentFormatter(rent)
		RentsFormatters = append(RentsFormatters, RentFormatter)
	}
	return RentsFormatters
}
