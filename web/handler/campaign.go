package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
	userService     user.Service
}

func NewCampaignHandler(campaignService campaign.Service, userService user.Service) *campaignHandler {
	return &campaignHandler{campaignService, userService}
}

func (h *campaignHandler) Index(c *gin.Context) {
	campaigns, err := h.campaignService.GetCampaigns(0)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "campaign_index.html", gin.H{"campaigns": campaigns})
}

func (h *campaignHandler) New(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	input := campaign.FormCreateCampaignInput{
		Users: users,
	}

	c.HTML(http.StatusOK, "campaign_new.html", input)
}

func (h *campaignHandler) Create(c *gin.Context) {
	var input campaign.FormCreateCampaignInput

	err := c.ShouldBind(&input)
	if err != nil {
		allUsers, e := h.userService.GetAllUsers()

		if e != nil {
			c.HTML(http.StatusInternalServerError, "error.html", nil)
			return
		}

		input.Error = err
		input.Users = allUsers

		c.HTML(http.StatusOK, "campaign_new.html", input)
		return
	}

	user, err := h.userService.GetUserByID(input.UserID)

	if err != nil {
		// input.Error = err
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}

	campaignInput := campaign.CreateCampaignInput{
		Name:             input.Name,
		Description:      input.Description,
		ShortDescription: input.ShortDescription,
		GoalAmount:       input.GoalAmount,
		Perks:            input.Perks,
		User:             user,
	}

	_, err = h.campaignService.CreateCampaign(campaignInput)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) NewImage(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	c.HTML(http.StatusOK, "campaign_image.html", gin.H{"ID": id})
}

func (h *campaignHandler) CreateImage(c *gin.Context) {
	// get form file uploaded
	file, err := c.FormFile("file")

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadGateway, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	existingCampaign, err := h.campaignService.GetCampaignById(campaign.GetCampaignDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	userId := existingCampaign.UsersID

	path := fmt.Sprintf("images/campaign/%d-%d-%s", userId, id, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	user, err := h.userService.GetUserByID(userId)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	createCampaignImageInput := campaign.CreateCampaignImageInput{
		CampaignId: id,
		IsPrimary:  true,
		User:       user,
	}

	_, err = h.campaignService.SaveCampaignImage(createCampaignImageInput, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) Edit(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	existingCampaign, err := h.campaignService.GetCampaignById(campaign.GetCampaignDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	input := campaign.FormUpdateCampaignInput{
		ID:               id,
		Name:             existingCampaign.Name,
		ShortDescription: existingCampaign.ShortDescription,
		Description:      existingCampaign.Description,
		GoalAmount:       existingCampaign.GoalAmount,
		Perks:            existingCampaign.Perks,
	}
	c.HTML(http.StatusOK, "campaign_edit.html", input)
}

func (h *campaignHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var input campaign.FormUpdateCampaignInput
	err := c.ShouldBind(&input)
	if err != nil {
		input.Error = err
		input.ID = id
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	existingCampaign, err := h.campaignService.GetCampaignById(campaign.GetCampaignDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	userCampaign, err := h.userService.GetUserByID(existingCampaign.UsersID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	campaignDetailById := campaign.GetCampaignDetailInput{
		ID: id,
	}

	campaignInput := campaign.CreateCampaignInput{
		Name:             input.Name,
		Description:      input.Description,
		ShortDescription: input.ShortDescription,
		GoalAmount:       input.GoalAmount,
		Perks:            input.Perks,
		User:             userCampaign,
	}

	_, err = h.campaignService.UpdateCampaign(campaignDetailById, campaignInput)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) Show(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)
	existingCampaign, err := h.campaignService.GetCampaignById(campaign.GetCampaignDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "campaign_show.html", existingCampaign)
}
