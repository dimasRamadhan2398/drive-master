package controllers

import (
	"core-service/models/dto"
	"core-service/pkg/response"
	"core-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IContactController interface {
	CreateInquiry(ctx *gin.Context)
	GetAllInquiries(ctx *gin.Context)
}

type ContactController struct {
	contactService services.IContactService
}

func NewContactController(contactService services.IContactService) IContactController {
	return &ContactController{
		contactService: contactService,
	}
}

func (c *ContactController) CreateInquiry(ctx *gin.Context) {
	var req dto.CreateContactInquiryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request payload: "+err.Error())
		return
	}

	contact, err := c.contactService.CreateInquiry(ctx.Request.Context(), &req)
	if err != nil {
		response.InternalServerError(ctx, "Failed to submit inquiry: "+err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Message submitted successfully",
		"data":    contact,
	})
}

func (c *ContactController) GetAllInquiries(ctx *gin.Context) {
	contacts, err := c.contactService.GetAllInquiries(ctx.Request.Context())
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch inquiries: "+err.Error())
		return
	}

	response.OK(ctx, "Contact inquiries fetched successfully", contacts)
}
