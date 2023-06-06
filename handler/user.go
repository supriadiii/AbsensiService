package handler

import (
	"net/http"
	"project_absensi/auth"
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
		errorMassage := gin.H{"error": errors}
		respone := helper.ResponseFormatter("Failed to register", http.StatusUnprocessableEntity, "error", errorMassage)
		c.JSON(http.StatusUnprocessableEntity, respone)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		errors := helper.ValidationErrorFormatter(err)
		errorMassage := gin.H{"error": errors}
		respone := helper.ResponseFormatter("Failed to register", http.StatusBadRequest, "error", errorMassage)
		c.JSON(http.StatusBadRequest, respone)
		return
	}
	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		respone := helper.ResponseFormatter("Failed to register", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, respone)
		return
	}
	formatter := user.FormatUser(newUser, token)
	respone := helper.ResponseFormatter("Account has been registered", http.StatusOK, "succes", formatter)
	c.JSON(http.StatusOK, respone)
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

func (h *userHandler) Login(c *gin.Context) {
	var input user.Login
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ValidationErrorFormatter(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ResponseFormatter("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedInUser, err := h.userService.Login(input)
	if err != nil {
		errors := helper.ValidationErrorFormatter(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ResponseFormatter("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedInUser.ID)
	if err != nil {
		respone := helper.ResponseFormatter("Login Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, respone)
		return
	}

	formatter := user.FormatUser(loggedInUser, token) // Anda perlu menghasilkan token JWT di sini
	response := helper.ResponseFormatter("Login success", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckNimAvailability(c *gin.Context) {
	var input user.CheckNimInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ValidationErrorFormatter(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ResponseFormatter("Nim checking Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	IsNimAvailable, err := h.userService.IsNimAvailable(input)
	if err != nil {
		errorMassage := gin.H{"errors": "server error"}
		response := helper.ResponseFormatter("Nim checking Failed", http.StatusUnprocessableEntity, "error", errorMassage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": IsNimAvailable,
	}

	metaMassage := "Nim has been registered"

	if IsNimAvailable {
		metaMassage = "Nim available"
	}
	response := helper.ResponseFormatter(metaMassage, http.StatusOK, "Succes", data)
	c.JSON(http.StatusOK, response)
}
