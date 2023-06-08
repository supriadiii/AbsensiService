package rent

type RentFormatter struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	SortDescription string `json:"sort_direction"`
	Description     string `json:"description"`
	ContactPerson   string `json:"contact_person"`
	Price           int    `json:"price"`
	Quantity        int    `json:"quantity"`
	UserID          uint   `json:"user_id"`
	ImageURL        string `json:"image_url"`
}

func FormatRentFormatter(rent Rent) RentFormatter {
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

	if len(rent.RentImage) > 0 {
		formatter.ImageURL = rent.RentImage[0].FileName
	}

	return formatter
}

func FormatRents(rents []Rent) []RentFormatter {
	if len(rents) == 0 {
		return []RentFormatter{}
	}
	var RentsFormatters []RentFormatter
	for _, user := range rents {
		RentFormatter := FormatRentFormatter(user)
		RentsFormatters = append(RentsFormatters, RentFormatter)

	}
	return RentsFormatters
}
