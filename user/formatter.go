package user

import "time"

type UserFormatter struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Nim       int       `json:"nim"`
	Kelas     string    `json:"kelas"`
	Prodi     string    `json:"prodi"`
	NoHp      string    `json:"no_hp"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserNoTokenFormatter struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Nim       int       `json:"nim"`
	Kelas     string    `json:"kelas"`
	Prodi     string    `json:"prodi"`
	NoHp      string    `json:"no_hp"`
	Role      string    `json:"role"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FormatUser(user User, token string) UserFormatter {
	formatter := UserFormatter{
		ID:        user.ID,
		Name:      user.Name,
		Nim:       user.Nim,
		Kelas:     user.Kelas,
		Prodi:     user.Prodi,
		NoHp:      user.NoHp,
		Token:     token,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return formatter
}

func FormatUserWithoutToken(user User) UserNoTokenFormatter {
	formatter := UserNoTokenFormatter{
		ID:        user.ID,
		Name:      user.Name,
		Nim:       user.Nim,
		Kelas:     user.Kelas,
		Prodi:     user.Prodi,
		NoHp:      user.NoHp,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return formatter
}

func FormatUsers(users []User) []UserNoTokenFormatter {
	if len(users) == 0 {
		return []UserNoTokenFormatter{}
	}
	var UserNoTokenFormatters []UserNoTokenFormatter
	for _, user := range users {
		UserNoTokenFormatter := FormatUserWithoutToken(user)
		UserNoTokenFormatters = append(UserNoTokenFormatters, UserNoTokenFormatter)

	}
	return UserNoTokenFormatters
}
