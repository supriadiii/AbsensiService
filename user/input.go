package user

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Nim      int    `json:"nim" binding:"required"`
	Password string `json:"password" binding:"required,alphanumunicode,min=8"`
	Kelas    string `json:"kelas" binding:"required"`
	Prodi    string `json:"prodi" binding:"required"`
	NoHp     string `json:"no_hp" binding:"required"`
}

type Login struct {
	Nim      int    `json:"nim" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UserIDInput struct {
	ID uint `json:"id" binding:"required"`
}

type CheckNimInput struct {
	Nim int `json:"nim" binding:"required"`
}
