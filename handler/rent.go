package handler

import (
	"net/http"
	"project_absensi/helper"
	"project_absensi/rent"
	"project_absensi/user"
	"strconv"

	"github.com/gin-gonic/gin"
)

type rentHandler struct {
	service rent.Service
}

func NewRentHandler(service rent.Service) *rentHandler {
	return &rentHandler{service}
}

func (h *rentHandler) GetRents(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	rents, err := h.service.GetRents(uint(userID))
	if err != nil {
		response := helper.ResponseFormatter("Error To Get Rent", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ResponseFormatter("list of rent", http.StatusOK, "succes", rent.FormatRents(rents))
	c.JSON(http.StatusOK, response)

}

func (h *rentHandler) CreateRent(c *gin.Context) {
	var input rent.CreateRentInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ValidationErrorFormatter(err)
		errorMassage := gin.H{"error": errors}
		respone := helper.ResponseFormatter("Failed to create Rent", http.StatusUnprocessableEntity, "error", errorMassage)
		c.JSON(http.StatusUnprocessableEntity, respone)
		return
	}
	currentUser := c.MustGet("CurrentUser").(user.User)
	input.User = currentUser
	newRent, err := h.service.CreateRent(input)
	if err != nil {
		respone := helper.ResponseFormatter("Failed to create Rent", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, respone)
		return
	}
	response := helper.ResponseFormatter("Succes create Rent", http.StatusOK, "succes", rent.FormatRentFormatter(newRent))
	c.JSON(http.StatusOK, response)

}
