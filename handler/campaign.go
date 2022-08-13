package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// tangkap parameter di handler
// handler send ke service
// service menentukan repository mana yg akan di call
// repository: GetAllCampaign, GetCampaignByUserId
// db

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

// http://localhost:8080/api/v1/campaigns?user_id=1
func (h *campaignHandler) FindCampaigns(c *gin.Context) {
	// to bind uri query and parse int from string
	userId, _ := strconv.Atoi(c.Query("user_id"))

	// passing userId to service GetCampaigns
	campaigns, err := h.service.GetCampaigns(userId)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	message := "List of campaigns"
	if len(campaigns) == 0 {
		message = "Campaigns is empty"
	}
	// success
	response := helper.APIResponse(message, http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
	return
}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	// bind data in uri and get input validation
	// http://localhost:8080/api/v1/campaign/{id}
	err := c.ShouldBindUri(&input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// passing input struct to service GetCampaignById
	campaignDetail, err := h.service.GetCampaignById(input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
	return
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	// bind data in JSON and get input validation
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse("Failed to create a campaign", http.StatusUnprocessableEntity, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// get user authentication
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	// store input struct to service CreateCampaign
	newCampaign, err := h.service.CreateCampaign(input)

	if err != nil {
		response := helper.APIResponse("Failed to create a campaign", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Succes to create a campaign", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
	return
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput

	// bind data in uri and get inputId validation
	// http://localhost:8080/api/v1/campaigns/{inputID}
	err := c.ShouldBindUri(&inputID)

	if err != nil {
		response := helper.APIResponse("Failed to update detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput

	// bind data in JSON and get inputData validation
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to update detail of campaign", http.StatusUnprocessableEntity, "failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	// get current user authentication
	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	// passing inputID and inputData to service UpdateCampaign
	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update detail of campaign", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Succes to update campaign", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
	return
}

func (h *campaignHandler) UploadImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput

	// bind data and get input validation error
	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusUnprocessableEntity, "failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	// get current user authentication
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	userId := input.User.ID

	// get form file uploaded
	file, err := c.FormFile("file")

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadGateway, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// store file to specified directory and set/get name
	// path := "images/campaign/" + file.Filename
	// update file name to "images/campaign/{userId}-{input.CampaignId}-{file.Filename}"
	path := fmt.Sprintf("images/campaign/%d-%d-%s", userId, input.CampaignId, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse(err.Error(), http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// passing input struct and path data to service
	_, err = h.service.SaveCampaignImage(input, path)
	if err != nil {
		errPath := os.Remove(path)
		if errPath != nil {
			fmt.Println("failed to remove file caused " + errPath.Error())
		}
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image. "+err.Error(), http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Successfuly to upload campaign image", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) DeleteCampaignImage(c *gin.Context) {
	var input campaign.GetCampaignImageDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to delete campaign image", http.StatusBadRequest, "failed", errorMessage)
		c.JSON(http.StatusBadRequest, response)
	}

	// get current user authentication
	currentUser := c.MustGet("currentUser").(user.User)
	fmt.Println(currentUser)
	userId := currentUser.ID

	campaignImage, err := h.service.GetCampaignImageById(input.ID)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Campaign image is not found", http.StatusNotFound, "failed", errorMessage)
		c.JSON(http.StatusNotFound, response)
	}

	var campaignDetail campaign.GetCampaignDetailInput
	campaignDetail.ID = campaignImage.ID
	campaign, err := h.service.GetCampaignById(campaignDetail)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Campaign image is not related with any campaign", http.StatusNotFound, "failed", errorMessage)
		c.JSON(http.StatusNotFound, response)
		return
	}

	fmt.Println("campaign user id : ", campaign.UsersID)
	fmt.Println("user id", userId)
	fmt.Println("is same : ", campaign.UsersID == userId)

	if campaign.UsersID != userId {
		errorMessage := gin.H{"errors": "Unauthorized"}
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "failed", errorMessage)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	_, err = h.service.DeleteCampaignImage(input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to delete campaign image", http.StatusBadRequest, "failed", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	errRemoveFile := os.Remove(campaignImage.FileName)
	if errRemoveFile != nil {
		errorMessage := gin.H{"errors": "failed to removing file"}
		response := helper.APIResponse("Failed to removing file", http.StatusConflict, "failed", errorMessage)
		c.JSON(http.StatusConflict, response)
	}
	data := gin.H{"is_deleted": true}
	response := helper.APIResponse("Successfuly to delete campaign image", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
