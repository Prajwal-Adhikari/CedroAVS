package types

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const SignatureCollection = "signature"

type Signature struct {
	ID        primitive.ObjectID `bson:"_id"`
	Chain     string             `bson:"chain"`
	RequestId string             `bson:"requestId"`
	Signer    string             `bson:"signer"`
	Signature string             `bson:"signature"`
	Message   string             `bson:"message"`
	IsDone    bool               `bson:"isDone"` //If the signature associated with requestid request is already done
}

type AddSignatureRequest struct {
	Chain     string `form:"chain"  binding:"required"`
	RequestId string `form:"requestId" binding:"required"`
	Signer    string `form:"signer" binding:"required"`
	Signature string `form:"signature" binding:"required"`
	Message   string `form:"message" binding:"required"`
}

type ErrorResponse struct {
	Message string
}

type SignatureRepository interface {
	CreateNewSignature(context.Context, Signature) error
	GetSignatureFromRequestId(context.Context, string) ([]Signature, error)
}
