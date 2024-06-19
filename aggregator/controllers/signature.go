package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Cedro-Finance/CedroAVS/aggregator/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignatureController struct {
	signatureRepository types.SignatureRepository
}

func NewSignatureController(sr types.SignatureRepository) *SignatureController {
	return &SignatureController{
		signatureRepository: sr,
	}
}

func (sc *SignatureController) AddSignature(c *gin.Context) {

	var request *types.AddSignatureRequest

	err := c.ShouldBind(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Message: err.Error()})
		return
	}

	signature := types.Signature{
		ID:        primitive.NewObjectID(),
		Chain:     request.Chain,
		RequestId: request.RequestId,
		Signer:    request.Signer,
		Signature: request.Signature,
		Message:   request.Message,
	}

	ctx, cancel := context.WithTimeout(c, time.Minute)
	defer cancel()

	err = sc.signatureRepository.CreateNewSignature(ctx, signature)

	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, []string{})
}

func (sc *SignatureController) GetSignatureFromRequestId(c *gin.Context) {

	requestId, ok := c.Params.Get("requestId")

	if !ok {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Message: "Request id not found in request"})
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Minute)
	defer cancel()

	signature, err := sc.signatureRepository.GetSignatureFromRequestId(ctx, requestId)

	if err != nil {
		c.JSON(http.StatusBadGateway, types.ErrorResponse{Message: "Error while getting signatures"})
		return
	}

	c.JSON(http.StatusOK, signature)
}
