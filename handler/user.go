package handler

import (
	"net/http"
	"project_absensi/helper"
	"project_absensi/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
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
