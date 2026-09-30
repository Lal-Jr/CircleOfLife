// Package httpapi exposes the Community context over HTTP. Handlers check the shape of the
// request, call one use case and translate its result or domain error into a response.
package httpapi

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"circleoflife/internal/community/application"
	"circleoflife/internal/community/domain"
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

// Register mounts the Community routes on a group that already requires authentication.
// Every route shares a general rate limit; writes get tighter ones of their own.
func (h *Handler) Register(posts *gin.RouterGroup, limit func(perMinute int) gin.HandlerFunc) {
	posts.Use(limit(60))
	posts.GET("", h.feed)
	posts.POST("", limit(5), h.publishPost)
	posts.GET("/:id", h.post)
	posts.POST("/:id/like", limit(30), h.toggleHelpful)
	posts.GET("/:id/comments", h.comments)
	posts.POST("/:id/comments", limit(10), h.addComment)
}

// Request shapes. Their type names appear in validation messages, so they're part of the API.

type CreatePostInput struct {
	Title       string     `json:"title" validate:"required,min=5,max=120"`
	Description string     `json:"description" validate:"required,min=10,max=1000"`
	Type        string     `json:"type" validate:"required,oneof=help meetup"`
	Lat         float64    `json:"lat" validate:"required,latitude"`
	Lng         float64    `json:"lng" validate:"required,longitude"`
	MeetupTime  *time.Time `json:"meetupTime"`
}

type CreateCommentInput struct {
	Content string `json:"content" validate:"required,max=500"`
}

type helpfulResponse struct {
	Liked        bool `json:"liked"`
	HelpfulCount int  `json:"helpfulCount"`
}

func (h *Handler) publishPost(c *gin.Context) {
	member, ok := httpx.CurrentUser(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var req CreatePostInput
	if !h.bind(c, &req) {
		return
	}
	view, err := h.app.PublishPost(c.Request.Context(), application.PublishPost{
		Author:      domain.MemberID(member),
		Title:       req.Title,
		Description: req.Description,
		Kind:        req.Type,
		Lat:         req.Lat,
		Lng:         req.Lng,
		MeetupTime:  req.MeetupTime,
	})
	if err != nil {
		h.fail(c, err, "create post", "Failed to create post")
		return
	}
	c.JSON(http.StatusCreated, httpx.APIResponse{Data: view})
}

func (h *Handler) feed(c *gin.Context) {
	latStr, lngStr := c.Query("lat"), c.Query("lng")
	if latStr == "" || lngStr == "" {
		utils.JSONError(c, http.StatusBadRequest, "Latitude and longitude are required")
		return
	}
	lat, errLat := strconv.ParseFloat(latStr, 64)
	lng, errLng := strconv.ParseFloat(lngStr, 64)
	radius, errRad := strconv.Atoi(c.DefaultQuery("radius", "5")) // km
	page, errPage := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if errLat != nil || errLng != nil || errRad != nil || errPage != nil || errLimit != nil {
		utils.JSONError(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}
	viewer, _ := httpx.CurrentUser(c)

	views, hasNext, err := h.app.Feed(c.Request.Context(), lat, lng, radius, page, limit, domain.MemberID(viewer))
	if err != nil {
		h.fail(c, err, "get nearby posts", "Failed to fetch nearby posts")
		return
	}
	if views == nil {
		views = []application.PostView{} // an empty feed is [], not null
	}
	c.JSON(http.StatusOK, httpx.APIResponse{
		Data: views,
		Meta: &httpx.PaginatedMeta{Page: page, Limit: limit, HasNext: hasNext},
	})
}

func (h *Handler) post(c *gin.Context) {
	// Coordinates are optional here; without them the distance is measured from (0, 0).
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	viewer, _ := httpx.CurrentUser(c)

	view, err := h.app.Post(c.Request.Context(), c.Param("id"), lat, lng, domain.MemberID(viewer))
	if err != nil {
		h.fail(c, err, "get post", "Failed to fetch post")
		return
	}
	c.JSON(http.StatusOK, httpx.APIResponse{Data: view})
}

func (h *Handler) toggleHelpful(c *gin.Context) {
	member, ok := httpx.CurrentUser(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	voted, total, err := h.app.ToggleHelpful(c.Request.Context(), c.Param("id"), domain.MemberID(member))
	if err != nil {
		h.fail(c, err, "toggle like", "Failed to update helpful vote")
		return
	}
	c.JSON(http.StatusOK, httpx.APIResponse{Data: helpfulResponse{Liked: voted, HelpfulCount: total}})
}

func (h *Handler) addComment(c *gin.Context) {
	member, ok := httpx.CurrentUser(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var req CreateCommentInput
	if !h.bind(c, &req) {
		return
	}
	view, err := h.app.AddComment(c.Request.Context(), c.Param("id"), domain.MemberID(member), req.Content)
	if err != nil {
		h.fail(c, err, "create comment", "Failed to post comment")
		return
	}
	c.JSON(http.StatusCreated, httpx.APIResponse{Data: view})
}

func (h *Handler) comments(c *gin.Context) {
	views, err := h.app.Comments(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.fail(c, err, "get comments", "Failed to retrieve comments")
		return
	}
	if views == nil {
		views = []application.CommentView{}
	}
	c.JSON(http.StatusOK, httpx.APIResponse{Data: views})
}

// bind decodes and validates a JSON body, answering 400 itself when it's unusable.
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

// fail maps a use case's error to a response: a missing post is 404, a broken domain rule is
// 400 with the rule's message, and anything else is a logged 500.
func (h *Handler) fail(c *gin.Context, err error, action, message string) {
	switch {
	case errors.Is(err, domain.ErrPostNotFound):
		utils.JSONError(c, http.StatusNotFound, "Post not found")
	case errors.Is(err, domain.ErrInvalidLocation), errors.Is(err, domain.ErrInvalidKind),
		errors.Is(err, domain.ErrInvalidTitle), errors.Is(err, domain.ErrInvalidText), errors.Is(err, domain.ErrInvalidComment):
		utils.JSONError(c, http.StatusBadRequest, err.Error())
	default:
		log.Printf("%s failed: %v", action, err)
		utils.JSONError(c, http.StatusInternalServerError, message)
	}
}
