// Package httpapi exposes the Identity context over HTTP.
package httpapi

import (
	"errors"
	"log"
	"net/http"

	"circleoflife/internal/identity/application"
	"circleoflife/internal/identity/domain"
	"circleoflife/internal/platform/httpx"
	"circleoflife/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	app      *application.Service
	validate *validator.Validate
}

func NewHandler(app *application.Service) *Handler {
	return &Handler{app: app, validate: validator.New()}
}

// Request shapes. Their type names appear in validation messages, so they're part of the API.

type SignupRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *Handler) Signup(c *gin.Context) {
	var req SignupRequest
	if !h.bind(c, &req) {
		return
	}
	token, err := h.app.SignUp(c.Request.Context(), req.Name, req.Email, req.Password)
	switch {
	case errors.Is(err, domain.ErrEmailTaken):
		utils.JSONError(c, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrInvalidEmail), errors.Is(err, domain.ErrInvalidName), errors.Is(err, domain.ErrWeakPassword):
		utils.JSONError(c, http.StatusBadRequest, err.Error())
	case err != nil:
		log.Printf("signup failed: %v", err)
		utils.JSONError(c, http.StatusInternalServerError, "could not create user")
	default:
		c.JSON(http.StatusCreated, gin.H{"token": token})
	}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if !h.bind(c, &req) {
		return
	}
	token, err := h.app.LogIn(c.Request.Context(), req.Email, req.Password)
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		utils.JSONError(c, http.StatusUnauthorized, err.Error())
	case err != nil:
		log.Printf("login failed: %v", err)
		utils.JSONError(c, http.StatusInternalServerError, "could not log in")
	default:
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func (h *Handler) Me(c *gin.Context) {
	id, ok := httpx.CurrentUser(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	profile, err := h.app.Profile(c.Request.Context(), id)
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		utils.JSONError(c, http.StatusNotFound, "User not found")
	case err != nil:
		log.Printf("profile failed: %v", err)
		utils.JSONError(c, http.StatusInternalServerError, "could not load profile")
	default:
		c.JSON(http.StatusOK, httpx.APIResponse{Data: profile})
	}
}

func (h *Handler) bind(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "Invalid request payload")
		return false
	}
	if err := h.validate.Struct(req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}
