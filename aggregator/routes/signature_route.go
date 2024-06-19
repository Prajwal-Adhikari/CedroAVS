package route

import (
	"context"

	"github.com/Cedro-Finance/CedroAVS/aggregator/controllers"
	"github.com/Cedro-Finance/CedroAVS/aggregator/repository"
	"github.com/Cedro-Finance/CedroAVS/aggregator/types"
	"github.com/gin-gonic/gin"
)

func NewSignatureRouter(ctx context.Context, group *gin.RouterGroup) {

	signRepo := repository.NewSignatureRepository(types.SignatureCollection)

	signController := controllers.NewSignatureController(signRepo)

	group.GET("/signature/:requestId", signController.GetSignatureFromRequestId)
	group.POST("/signature", signController.AddSignature)
}
