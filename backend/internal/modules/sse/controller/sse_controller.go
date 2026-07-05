package controller

import (
	"fmt"
	"io"
	"time"

	"github.com/Kar-Su/uas-mobile.git/internal/package/sse"
	"github.com/gin-gonic/gin"
)

type SSEController interface {
	Stream(ctx *gin.Context)
}

type sseController struct{}

func NewSSEController() SSEController {
	return &sseController{}
}

func (c *sseController) Stream(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("X-Accel-Buffering", "no")

	hub := sse.Default()
	ch := hub.Subscribe()
	defer hub.Unsubscribe(ch)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	ctx.Stream(func(w io.Writer) bool {
		select {
		case event, ok := <-ch:
			if !ok {
				return false
			}
			fmt.Fprintf(w, "event: %s\ndata: changed\n\n", event)
			return true
		case <-ticker.C:
			fmt.Fprintf(w, ": ping\n\n")
			return true
		case <-ctx.Request.Context().Done():
			return false
		}
	})
}

