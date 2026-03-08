package handlers

import (
	"encoding/json"
	"fmt"
	"io"

	"circleoflife/internal/services"

	"github.com/gin-gonic/gin"
)

// HandleEvents creates an active SSE (Server-Sent Events) stream down to the client.
func HandleEvents(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	// Optional flush if a proxy doesn't chunk correctly
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// Create and register a new client with the global Hub
	clientChan := services.Hub.AddClient()
	
	// Unregister client gracefully when the HTTP request context closes (e.g. client disconnects via browser navigate)
	defer services.Hub.RemoveClient(clientChan)
	defer fmt.Println("Cleaned up SSE Context manually")

	c.Stream(func(w io.Writer) bool {
		if event, ok := <-clientChan; ok {
			eventJSON, _ := json.Marshal(event)
			
			// SSE Protocol Format demands starting strings with `data: ` followed by two newline characters.
			c.SSEvent("message", string(eventJSON))
			return true
		}
		return false
	})
}
