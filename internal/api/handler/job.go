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
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "User role not found")
		return
	}

	validationCtx := validator.ValidationContext{
		Role:      validator.UserRole(userRole.(string)),
		UserID:    userID.(uuid.UUID).String(),
		CompanyID: userID.(uuid.UUID).String(),
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
		Status:          "active",
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
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "User role not found")
		return
	}

	validationCtx := validator.ValidationContext{
		Role:      validator.UserRole(userRole.(string)),
		UserID:    userID.(uuid.UUID).String(),
		CompanyID: userID.(uuid.UUID).String(),
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

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 10
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

func (h *JobHandler) ChangeStatus(c *gin.Context) {
	// Parse job ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid job ID", err.Error())
		return
	}

	// Parse status from request body
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Get user context
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
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
		response.Error(c, http.StatusForbidden, "Unauthorized", "Not allowed to change job status")
		return
	}

	// Change job status
	if err := h.jobService.ChangeJobStatus(c.Request.Context(), id, req.Status); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to change job status", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Job status updated successfully", nil)
}

func (h *JobHandler) CompleteDelete(c *gin.Context) {
	// Parse job ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid job ID", err.Error())
		return
	}

	// Get user context
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
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
		response.Error(c, http.StatusForbidden, "Unauthorized", "Not allowed to delete this job")
		return
	}

	// Complete delete job and its applications
	if err := h.jobService.CompleteDeleteJob(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete job", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Job completely deleted successfully", nil)
}

func (h *JobHandler) GetJobAnalytics(c *gin.Context) {
    userID, _ := c.Get("userID")
    companyID := userID.(uuid.UUID)

    analytics, err := h.jobService.GetJobAnalytics(c.Request.Context(), companyID)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to fetch job analytics", err.Error())
        return
    }

    response.Success(c, http.StatusOK, "Job analytics retrieved", analytics)
}

func (h *JobHandler) BulkCreateJobs(c *gin.Context) {
    var jobs []domain.CreateJobRequest
    if err := c.ShouldBindJSON(&jobs); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
        return
    }

    userID, _ := c.Get("userID")
    createdJobs, err := h.jobService.BulkCreateJobs(c.Request.Context(), jobs, userID.(uuid.UUID))
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to create jobs", err.Error())
        return
    }

    response.Success(c, http.StatusCreated, "Jobs created successfully", createdJobs)
}

func (h *JobHandler) GetJobApplicationInsights(c *gin.Context) {
    idStr := c.Param("id")
    jobID, _ := uuid.Parse(idStr)

    insights, err := h.jobService.GetJobApplicationInsights(c.Request.Context(), jobID)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to fetch application insights", err.Error())
        return
    }

    response.Success(c, http.StatusOK, "Job application insights retrieved", insights)
}

func (h *JobHandler) GetRecommendedCandidates(c *gin.Context) {
    idStr := c.Param("id")
    jobID, _ := uuid.Parse(idStr)

    candidates, err := h.jobService.GetRecommendedCandidates(c.Request.Context(), jobID)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to fetch recommended candidates", err.Error())
        return
    }

    response.Success(c, http.StatusOK, "Recommended candidates retrieved", candidates)
}