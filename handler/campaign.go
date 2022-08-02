package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"net/http"
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

// api/v1/campaigns
func (h *campaignHandler) FindCampaigns(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userId)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if len(campaigns) == 0 {
		response := helper.APIResponse("Campaigns not found", http.StatusNotFound, "failed", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}
	// success
	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaigns)
	c.JSON(http.StatusOK, response)
	return
}
