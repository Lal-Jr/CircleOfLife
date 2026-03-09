package handlers

import (
	"net/http"
	"strconv"

	"circleoflife/internal/events"
	"circleoflife/internal/models"
	"circleoflife/internal/services"
	"circleoflife/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PostHandler struct {
	postService services.PostService
	validator   *validator.Validate
}

func NewPostHandler(s services.PostService) *PostHandler {
	return &PostHandler{
		postService: s,
		validator:   validator.New(),
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.JSONError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreatePostInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	post := models.Post{
		UserID:      userID.(string),
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		MeetupTime:  req.MeetupTime,
	}

	if err := h.postService.CreatePost(c.Request.Context(), &post, req.Lat, req.Lng); err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "Failed to create post")
		return
	}

	// Broadcast that a new post was created nearby via Redis Pub/Sub
	events.PublishPostCreated(post.ID, req.Lat, req.Lng)

	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetNearbyPosts(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.DefaultQuery("radius", "5") // Default 5 km
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	if latStr == "" || lngStr == "" {
		utils.JSONError(c, http.StatusBadRequest, "Latitude and longitude are required")
		return
	}

	lat, errLat := strconv.ParseFloat(latStr, 64)
	lng, errLng := strconv.ParseFloat(lngStr, 64)
	radius, errRad := strconv.Atoi(radiusStr)
	page, errPage := strconv.Atoi(pageStr)
	limit, errLimit := strconv.Atoi(limitStr)

	if errLat != nil || errLng != nil || errRad != nil || errPage != nil || errLimit != nil {
		utils.JSONError(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	posts, err := h.postService.GetNearbyPosts(c.Request.Context(), lat, lng, radius, page, limit)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch nearby posts")
		return
	}

	// Make sure we always return an array, even if empty, instead of null
	if posts == nil {
		posts = []models.Post{}
	}

	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	id := c.Param("id")
	
	latStr := c.Query("lat")
	lngStr := c.Query("lng")

	// Even if fetching exactly one post, it helps to pass lat/lng natively to calculate their relative distance
	lat, _ := strconv.ParseFloat(latStr, 64)
	lng, _ := strconv.ParseFloat(lngStr, 64)

	post, err := h.postService.GetPostByID(c.Request.Context(), id, lat, lng)
	if err != nil {
		utils.JSONError(c, http.StatusNotFound, "Post not found")
		return
	}

	c.JSON(http.StatusOK, post)
}
