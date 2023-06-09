package handler

import (
	"fmt"
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

func (h *rentHandler) SaveRentImage(c *gin.Context) {
	var input rent.CreateRentImageInput
	err := c.ShouldBind(&input)

	if err != nil {
		data := gin.H{"is_uploaded": false, "message": err.Error()}
		response := helper.ResponseFormatter("failed to upload  Image", http.StatusUnprocessableEntity, "error", data)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false, "message": err.Error()}
		response := helper.ResponseFormatter("failed to upload  Image", http.StatusUnprocessableEntity, "error", data)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	currentUser := c.MustGet("CurrentUser").(user.User)
	userID := currentUser.ID
	path := fmt.Sprintf("image/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false, "message": err.Error()}
		response := helper.ResponseFormatter("failed to upload  Image", http.StatusUnprocessableEntity, "error", data)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	_, err = h.service.SaveRentImage(input, path)
	if err != nil {
		data := gin.H{"is_uploaded": false, "message": err.Error()}
		response := helper.ResponseFormatter("failed to upload  Image", http.StatusUnprocessableEntity, "error", data)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	data := gin.H{"is_uploaded": true}
	response := helper.ResponseFormatter("Image succes uploaded", http.StatusOK, "succes", data)
	c.JSON(http.StatusOK, response)
}
