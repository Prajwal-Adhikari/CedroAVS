package repository

import (
	"context"

	db "github.com/Cedro-Finance/CedroAVS/aggregator/services"
	"github.com/Cedro-Finance/CedroAVS/aggregator/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type signatureRepository struct {
	database   mongo.Database
	collection string
}

func NewSignatureRepository(collection string) types.SignatureRepository {
	return &signatureRepository{
		database:   *db.DbClient,
		collection: collection,
	}
}

func (sr *signatureRepository) CreateNewSignature(ctx context.Context, signData types.Signature) error {
	collection := sr.database.Collection(sr.collection)

	_, err := collection.InsertOne(ctx, signData)

	return err
}

func (sr *signatureRepository) GetSignatureFromRequestId(ctx context.Context, reqId string) ([]types.Signature, error) {
	collection := sr.database.Collection(sr.collection)

	cur, err := collection.Find(ctx, bson.M{"ReqId": reqId})

	if err != nil {
		return []types.Signature{}, err
	}

	var signatures []types.Signature

	err = cur.All(ctx, &signatures)

	if signatures == nil {
		return []types.Signature{}, err
	}

	return signatures, err
}
