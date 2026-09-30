package realtime

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

// Stream serves the Server-Sent Events endpoint: each event on the bus becomes a "message".
func (b *Bus) Stream(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	events := b.hub.add()
	defer b.hub.remove(events)

	c.Stream(func(w io.Writer) bool {
		e, ok := <-events
		if !ok {
			return false
		}
		payload, _ := json.Marshal(e)
		c.SSEvent("message", string(payload))
		return true
	})
}
