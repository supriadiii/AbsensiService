package handler

import (
	"net/http"
	"project_absensi/handler/auth"
	"project_absensi/helper"
	"project_absensi/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ValidationErrorFormatter(err)
		errorMassage := gin.H{"ersror": errors}
		respone := helper.ResponseFormatter("Failed to register", http.StatusUnprocessableEntity, "error", errorMassage)
		c.JSON(http.StatusUnprocessableEntity, respone)
	}

	user, err := h.userService.RegisterUser(input)
	if err != nil {
		errors := helper.ValidationErrorFormatter(err)
		errorMassage := gin.H{"error": errors}
		respone := helper.ResponseFormatter("Failed to register", http.StatusBadRequest, "error", errorMassage)
		c.JSON(http.StatusBadRequest, respone)
	}
	c.JSON(http.StatusOK, user)
}

func (h *userHandler) GetAllUsers(c *gin.Context) {
	// pagination := helper.GeneratePaginationFromRequest(c)
	var input user.UserIDInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ValidationErrorFormatter(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ResponseFormatter("Failed to get User detail ", http.StatusUnprocessableEntity, "errors", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	Users, err := h.userService.GetAllUsers(input)
	if err != nil {
		errors := helper.ValidationErrorFormatter(err)
		errorMassage := gin.H{"error": errors}
		respone := helper.ResponseFormatter("Failed to Get All Ship", http.StatusUnprocessableEntity, "error", errorMassage)
		c.JSON(http.StatusUnprocessableEntity, respone)
	}
	formatter := user.FormatUsers(Users)
	respone := helper.ResponseFormatter("List Of User", http.StatusOK, "succes", formatter)
	c.JSON(http.StatusOK, respone)
}

// login
func (h *userHandler) LoginUser(c *gin.Context) {
	token := c.GetHeader("Authorization-Agent")
	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ValidationErrorFormatter(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.ResponseFormatter("Failed to login", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	user1, err := h.userService.Login(input, token)
	if err != nil {
		response := helper.ResponseFormatter("Failed to login", http.StatusBadRequest, "error", gin.H{"errors": [1]string{err.Error()}})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err = h.authService.GenerateToken(user1.ID)
	if err != nil {
		response := helper.ResponseFormatter("Failed to login", http.StatusBadRequest, "error", gin.H{"errors": [1]string{err.Error()}})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(user1, token)

	response := helper.ResponseFormatter("Login success", http.StatusOK, "success", formatter)

	if err != nil {
		response := helper.ResponseFormatter("Failed to login", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
