// Package httpx holds the HTTP conventions shared by every context: the response envelope and
// how a handler finds the signed-in member.
package httpx

import "github.com/gin-gonic/gin"

type PaginatedMeta struct {
	Page    int  `json:"page"`
	Limit   int  `json:"limit"`
	HasNext bool `json:"hasNext"`
}

type APIResponse struct {
	Data any            `json:"data"`
	Meta *PaginatedMeta `json:"meta,omitempty"`
}

// CurrentUser returns the ID the auth middleware put on the request, if any.
func CurrentUser(c *gin.Context) (string, bool) {
	v, ok := c.Get("userID")
	if !ok {
		return "", false
	}
	id, ok := v.(string)
	return id, ok
}
