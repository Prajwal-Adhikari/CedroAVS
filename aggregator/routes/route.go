package route

import (
	"context"

	"github.com/gin-gonic/gin"
)

func Setup(ctx context.Context, gin *gin.Engine) {
	publicRouter := gin.Group("")

	NewSignatureRouter(ctx, publicRouter)
}