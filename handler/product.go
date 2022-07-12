package handler

import (
	"bwastartup/helper"
	"bwastartup/product"
	"net/http"

	"github.com/gin-gonic/gin"
)

type productHandler struct {
	productService product.ServiceProduct
}

func NewProductHandler(productService product.ServiceProduct) *productHandler {
	return &productHandler{productService}
}

func (h *productHandler) RegisterProduct(c *gin.Context) {

	var input product.RegisterProductInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse("Register product failed", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	NewProduct, err := h.productService.RegisterProduct(input)
	if err != nil {
		response := helper.APIResponse("Register product failed", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	formatter := product.FormatProduct(NewProduct)

	response := helper.APIResponse("Product has been created", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)

}
