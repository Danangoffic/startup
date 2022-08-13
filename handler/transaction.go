package handler

import (
	"bwastartup/helper"
	"bwastartup/transaction"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *TransactionHandler {
	return &TransactionHandler{service}
}

func (h *TransactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetCampaignTransactionsInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to get Campaign's Transactions", http.StatusBadRequest, "failed", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// get user authentication
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	transactions, err := h.service.GetCampaignTransactions(input)
	if err != nil {
		response := helper.APIResponse(err.Error(), http.StatusNotFound, "failed", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := helper.APIResponse("Campaign's Transactions", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
	return
}

// parameterdi uri
// tangkap parameter mapping input struct
// call service, input struct as parameter
// service bawa campaign id, call repository
// repository find campaign transaction
