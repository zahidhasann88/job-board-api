package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zahidhasann88/job-board-api/internal/domain"
	"github.com/zahidhasann88/job-board-api/internal/service"
)

type JobHandler struct {
	jobService *service.JobService
}

func NewJobHandler(jobService *service.JobService) *JobHandler {
	return &JobHandler{jobService: jobService}
}

type CreateJobRequest struct {
	Title           string   `json:"title" binding:"required"`
	Description     string   `json:"description" binding:"required"`
	Location        string   `json:"location" binding:"required"`
	SalaryRange     *string  `json:"salary_range"`
	JobType         string   `json:"job_type" binding:"required"`
	ExperienceLevel string   `json:"experience_level" binding:"required"`
	Skills          []string `json:"skills" binding:"required"`
}

func (h *JobHandler) Create(c *gin.Context) {
	var req CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get company ID from authenticated user
	userID, _ := c.Get("userID")
	companyID := userID.(uuid.UUID)

	job := &domain.Job{
		Title:           req.Title,
		Description:     req.Description,
		CompanyID:       companyID,
		Location:        req.Location,
		SalaryRange:     req.SalaryRange,
		JobType:         req.JobType,
		ExperienceLevel: req.ExperienceLevel,
		Skills:          req.Skills,
	}

	if err := h.jobService.CreateJob(c.Request.Context(), job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, job)
}

func (h *JobHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job ID"})
		return
	}

	job, err := h.jobService.GetJob(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if job == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (h *JobHandler) List(c *gin.Context) {
	var filter domain.JobFilter
	// Get query parameters
	if loc := c.Query("location"); loc != "" {
		filter.Location = &loc
	}
	if jobType := c.Query("job_type"); jobType != "" {
		filter.JobType = &jobType
	}
	if expLevel := c.Query("experience_level"); expLevel != "" {
		filter.ExperienceLevel = &expLevel
	}
	filter.Skills = c.QueryArray("skills")

	// Parse pagination
	if pageStr := c.Query("page"); pageStr != "" {
		filter.Page, _ = strconv.Atoi(pageStr)
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		filter.PageSize, _ = strconv.Atoi(pageSizeStr)
	}

	jobs, total, err := h.jobService.ListJobs(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      jobs,
		"total":     total,
		"page":      filter.Page,
		"page_size": filter.PageSize,
	})
}
