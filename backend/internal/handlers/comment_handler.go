package handlers

import (
	"net/http"

	"circleoflife/internal/models"
	"circleoflife/internal/services"
	"circleoflife/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
)

type CommentHandler struct {
	commentService services.CommentService
	validator      *validator.Validate
}

func NewCommentHandler(s services.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: s,
		validator:      validator.New(),
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	postID := c.Param("id")
	userID, exists := c.Get("userID")
	if !exists {
		utils.JSONError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreateCommentInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	p := bluemonday.StrictPolicy()
	comment := models.Comment{
		PostID:  postID,
		UserID:  userID.(string),
		Content: p.Sanitize(req.Content),
	}

	if err := h.commentService.CreateComment(c.Request.Context(), &comment); err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "Failed to post comment")
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{Data: comment})
}

func (h *CommentHandler) GetComments(c *gin.Context) {
	postID := c.Param("id")

	comments, err := h.commentService.GetCommentsByPostID(c.Request.Context(), postID)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "Failed to retrieve comments")
		return
	}

	if comments == nil {
		comments = []models.Comment{}
	}

	c.JSON(http.StatusOK, models.APIResponse{Data: comments})
}
