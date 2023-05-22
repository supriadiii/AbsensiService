package user

import "errors"

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Nim      int    `json:"nim" binding:"required"`
	Password string `json:"password" binding:"required,alphanumunicode,min=8"`
	// CPassword string `json:"cpassword" binding:"required,eqfield=Password"`
	Role  string `json:"role" `
	Kelas string `json:"kelas" binding:"required"`
	Prodi string `json:"prodi" binding:"required"`
	NoHp  string `json:"no_hp" binding:"required"`
	PT    string `json:"pt" binding:"required"`
}

type Login struct {
	Nim      string `json:"nim" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (input *RegisterUserInput) Validate() error {
	// Lakukan validasi sesuai kebutuhan
	if input.Name == "" {
		return errors.New("Name is required ff")
	}
	if input.Nim <= 0 {
		return errors.New("Nim is required")
	}
	// ... tambahkan validasi lainnya
	return nil
}
