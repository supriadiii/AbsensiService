package user

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Nim      int    `json:"nim" binding:"required"`
	Password string `json:"password" binding:"required,alphanumunicode,min=8"`
	// CPassword string `json:"cpassword" binding:"required,eqfield=Password"`
	Role  string `json:"role" `
	Kelas string `json:"kelas" binding:"required"`
	Prodi string `json:"prodi" binding:"required"`
	NoHp  string `json:"no_hp" binding:"required"`
}

type Login struct {
	Nim      string `json:"nim" binding:"required"`
	Password string `json:"password" binding:"required"`
}
