package handlers

import (
	"net/http"

	"circleoflife/internal/models"
	"circleoflife/internal/services"
	"circleoflife/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService services.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(s services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: s,
		validator:   validator.New(),
	}
}

type SignupRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	token, err := h.authService.Signup(c.Request.Context(), user)
	if err != nil {
		utils.JSONError(c, http.StatusConflict, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		utils.JSONError(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
