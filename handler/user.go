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
		errorMassage := gin.H{"error": errors}
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
