package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zahidhasann88/job-board-api/internal/domain"
	"github.com/zahidhasann88/job-board-api/internal/service"
)

type ApplicationHandler struct {
	applicationService *service.ApplicationService
}

func NewApplicationHandler(applicationService *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{applicationService: applicationService}
}

type CreateApplicationRequest struct {
	JobID       uuid.UUID `json:"job_id" binding:"required"`
	CoverLetter string    `json:"cover_letter"`
	ResumeURL   string    `json:"resume_url" binding:"required"`
}

func (h *ApplicationHandler) Create(c *gin.Context) {
	var req CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	application := &domain.Application{
		JobID:       req.JobID,
		ApplicantID: userID.(uuid.UUID),
		CoverLetter: req.CoverLetter,
		ResumeURL:   req.ResumeURL,
		Status:      "pending",
	}

	if err := h.applicationService.Create(c.Request.Context(), application); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, application)
}

func (h *ApplicationHandler) List(c *gin.Context) {
	userID, _ := c.Get("userID")
	applications, err := h.applicationService.ListByUser(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, applications)
}