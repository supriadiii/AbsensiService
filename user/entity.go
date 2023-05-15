package user

import "project_absensi/common"

type User struct {
	common.Model
	Nim      int
	Password string
	Name     string
	Kelas    string
	Prodi    string
	NoHp     string
	Role     string
}
