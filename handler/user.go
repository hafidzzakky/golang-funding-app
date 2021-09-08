package handler

import (
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"

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

	// Error validation handler
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Send input data to service
	newUser, err := h.userService.RegisterUser(input)
	// Error sending
	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Success sending
	formatter := user.FormatUser(newUser, "xxxxxxxx")
	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Loagin account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedInUser, err := h.userService.LoginUser(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Loagin account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loginFormatter := user.FormatUser(loggedInUser, "xxxxx")
	response := helper.APIResponse("Login successfully", http.StatusOK, "success", loginFormatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailvailability(c *gin.Context) {
	var input user.CheckEmailInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Email is not available", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// send input data to service
	isAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "An system error occurred"}
		response := helper.APIResponse("Email is not available", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Success sending
	data := gin.H{
		"is_availabile": isAvailable,
	}

	metaMessage := "Email has been registered"
	if isAvailable {
		metaMessage = "Email is available"
	}
	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	fileUpload, err := c.FormFile("avatar")
	// Failed Upload Image
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed Upload Avatar", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadGateway, response)
		return
	}
	userID := 1
	path := fmt.Sprintf("images/%d-%s", userID, fileUpload.Filename)
	err = c.SaveUploadedFile(fileUpload, path)
	// Failed Save Image
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed Save Uploaded Avatar", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadGateway, response)
		return
	}

	// Save data to DB
	_, err = h.userService.SaveAvatar(userID, path)

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed Save Uploaded Avatar to Database", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadGateway, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Success Upload Avatar", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
