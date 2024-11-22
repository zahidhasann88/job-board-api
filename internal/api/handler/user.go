package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zahidhasann88/job-board-api/internal/domain"
	"github.com/zahidhasann88/job-board-api/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		Email:       req.Email,
		Role:        req.Role,
		FullName:    req.FullName,
		CompanyName: req.CompanyName,
	}

	if err := h.userService.CreateUser(c.Request.Context(), user, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) UpdateProfileDetails(c *gin.Context) {
	var req domain.UpdateProfileDetailsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	// Convert request to map for flexible updates
	updates := make(map[string]interface{})

	if req.Skills != nil {
		updates["skills"] = req.Skills
	}
	if req.Experience != nil {
		updates["experience"] = req.Experience
	}
	if req.Education != nil {
		updates["education"] = req.Education
	}
	if req.Bio != nil {
		updates["bio"] = req.Bio
	}
	if req.ProfilePictureURL != nil {
		updates["profile_picture_url"] = req.ProfilePictureURL
	}
	if req.Location != nil {
		updates["location"] = req.Location
	}
	if req.SocialLinks != nil {
		updates["social_links"] = req.SocialLinks
	}
	if req.ContactInfo != nil {
		updates["contact_info"] = req.ContactInfo
	}

	if err := h.userService.UpdateProfileDetails(c.Request.Context(), userID.(uuid.UUID), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile details updated successfully"})
}

func (h *UserHandler) UpdateEmploymentHistory(c *gin.Context) {
	var history []domain.EmploymentHistory
	if err := c.ShouldBindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	if err := h.userService.UpdateEmploymentHistory(c.Request.Context(), userID.(uuid.UUID), history); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employment history updated successfully"})
}
