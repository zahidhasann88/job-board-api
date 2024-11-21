package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zahidhasann88/job-board-api/internal/domain"
	"github.com/zahidhasann88/job-board-api/internal/service"
	"github.com/zahidhasann88/job-board-api/pkg/response"
	"github.com/zahidhasann88/job-board-api/pkg/validator"
)

type JobHandler struct {
	jobService      *service.JobService
	customValidator *validator.CustomValidator
}

func NewJobHandler(jobService *service.JobService) *JobHandler {
	return &JobHandler{
		jobService:      jobService,
		customValidator: validator.NewValidator(),
	}
}

func (h *JobHandler) Create(c *gin.Context) {
	var req domain.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Get user role and context from authenticated user
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	validationCtx := validator.ValidationContext{
		Role:      validator.UserRole(userRole.(string)),
		UserID:    userID.(string),
		CompanyID: userID.(string),
	}

	// Validate the request using custom validator
	if err := h.customValidator.ValidateWithRole(req, validationCtx); err != nil {
		response.Error(c, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	job := &domain.Job{
		Title:           req.Title,
		Description:     req.Description,
		CompanyID:       userID.(uuid.UUID),
		Location:        req.Location,
		SalaryRange:     req.SalaryRange,
		JobType:         req.JobType,
		ExperienceLevel: req.ExperienceLevel,
		Skills:          req.Skills,
	}

	if err := h.jobService.CreateJob(c.Request.Context(), job); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create job", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Job created successfully", job)
}

func (h *JobHandler) Update(c *gin.Context) {
	// Parse job ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid job ID", err.Error())
		return
	}

	// Parse request body
	var req domain.UpdateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Get user role and context from authenticated user
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	validationCtx := validator.ValidationContext{
		Role:      validator.UserRole(userRole.(string)),
		UserID:    userID.(string),
		CompanyID: userID.(string),
	}

	// Validate the request using custom validator
	if err := h.customValidator.ValidateWithRole(req, validationCtx); err != nil {
		response.Error(c, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	// Fetch existing job to verify ownership
	existingJob, err := h.jobService.GetJob(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch job", err.Error())
		return
	}
	if existingJob == nil {
		response.Error(c, http.StatusNotFound, "Job not found", "")
		return
	}

	// Verify job belongs to current user
	if existingJob.CompanyID != userID {
		response.Error(c, http.StatusForbidden, "Unauthorized", "Not allowed to update this job")
		return
	}

	// Update job fields
	job := &domain.Job{
		ID:              id,
		Title:           req.Title,
		Description:     req.Description,
		Location:        req.Location,
		SalaryRange:     req.SalaryRange,
		JobType:         req.JobType,
		ExperienceLevel: req.ExperienceLevel,
		Skills:          req.Skills,
		Status:          req.Status,
	}

	// Perform update
	if err := h.jobService.UpdateJob(c.Request.Context(), job); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update job", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Job updated successfully", job)
}

func (h *JobHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid job ID", err.Error())
		return
	}

	job, err := h.jobService.GetJob(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch job", err.Error())
		return
	}
	if job == nil {
		response.Error(c, http.StatusNotFound, "Job not found", "")
		return
	}

	response.Success(c, http.StatusOK, "Job retrieved successfully", job)
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
		response.Error(c, http.StatusInternalServerError, "Failed to list jobs", err.Error())
		return
	}

	meta := response.Meta{
		Total:     total,
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		TotalPage: (total + filter.PageSize - 1) / filter.PageSize,
	}

	response.SuccessWithMeta(c, http.StatusOK, "Jobs retrieved successfully", jobs, meta)
}
