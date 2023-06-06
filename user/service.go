package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	GetAllUsers(input UserIDInput) ([]User, error)
	Login(input Login) (User, error)
	IsNimAvailable(input CheckNimInput) (bool, error)
	GetUserByID(ID uint) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}

}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Nim = input.Nim
	Password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.Password = string(Password)
	user.Kelas = input.Kelas
	user.NoHp = input.NoHp
	user.Prodi = input.Prodi
	user.Role = "user"

	// Validasi nim di API
	apiURL := fmt.Sprintf("https://api-frontend.kemdikbud.go.id/hit_mhs/%v", input.Nim)
	resp, err := http.Get(apiURL)
	if err != nil {
		return user, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Invalid response from API:", resp.StatusCode)
		return user, errors.New("Data nim tidak valid")
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Invalid nim data:", input.Nim)
		return user, err
	}

	type ApiResponse struct {
		Mahasiswa []struct {
			Text        string `json:"text"`
			WebsiteLink string `json:"website-link"`
		} `json:"mahasiswa"`
	}

	var apiResponse ApiResponse

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return user, err
	}

	if len(apiResponse.Mahasiswa) <= 0 || apiResponse.Mahasiswa[0].WebsiteLink == "/data_mahasiswa/" {
		return user, errors.New("Data nim tidak valid")
	}

	var targetMahasiswa struct {
		Text        string `json:"text"`
		WebsiteLink string `json:"website-link"`
	}

	for _, mhs := range apiResponse.Mahasiswa {
		if strings.Contains(mhs.Text, "UNIVERSITAS NEGERI MEDAN") {
			targetMahasiswa = mhs
			break
		}
	}

	if targetMahasiswa.WebsiteLink == "" {
		return user, errors.New("Daftar hanya tersedia untuk UNIVERSITAS NEGERI MEDAN")
	}

	ptInfo := strings.Split(targetMahasiswa.Text, ", PT : ")
	if len(ptInfo) < 2 {
		return user, errors.New("Invalid PT data from API")
	}

	pt := strings.TrimSpace(ptInfo[1])
	prodiInfo := strings.Split(pt, ", Prodi: ")
	if len(prodiInfo) < 2 {
		return user, errors.New("Invalid Prodi data from API")
	}

	user.PT = strings.TrimPrefix(prodiInfo[0], "PT : ")

	existingUser, err := s.repository.FindByNim(input.Nim)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return user, err
	}
	if existingUser.ID != 0 {
		return user, errors.New("NIM sudah terdaftar")
	}

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func (s *service) GetAllUsers(input UserIDInput) ([]User, error) {
	Ships, err := s.repository.GetAllUsers(input.ID)
	if err != nil {
		return nil, err
	}
	return Ships, nil
}

func (s *service) Login(input Login) (User, error) {
	// mencari user berdasarkan Nim
	nim := input.Nim
	passwword := input.Password

	user, err := s.repository.FindByNim(nim)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("No user found on that NIM")
	}
	// membandingkan password user dengan password input
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwword))
	if err != nil {
		return user, err
	}

	// jika semuanya baik-baik saja, kembalikan user
	return user, nil
}

func (s *service) IsNimAvailable(input CheckNimInput) (bool, error) {
	nim := input.Nim
	user, err := s.repository.FindByNim(nim)
	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}
	return false, nil
}

func (s *service) GetUserByID(ID uint) (User, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}
	if user.ID == 0 {
		return user, errors.New("No user found on with that ID")
	}
	return user, nil
}
