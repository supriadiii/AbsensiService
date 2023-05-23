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
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	GetAllUsers(input UserIDInput) ([]User, error)
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
	user.Prodi = input.Prodi
	user.Role = "user"
	user.PT = input.PT

	// Validasi nim di API
	apiURL := fmt.Sprintf("https://api-frontend.kemdikbud.go.id/hit_mhs/%v", input.Nim)
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Println("Error calling API:", err)
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

	// Validasi nama PT
	if !strings.Contains(apiResponse.Mahasiswa[0].Text, "UNIVERSITAS NEGERI MEDAN") {
		return user, errors.New("Daftar hanya tersedia untuk UNIVERSITAS NEGERI MEDAN")
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
