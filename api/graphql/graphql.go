package graphql

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/Layr-Labs/eigensdk-go/logging"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/shurcooL/graphql"
)

type GraphQLClient struct {
	client *graphql.Client
	logger logging.Logger
}

func NewGraphQLClient(url string, logger logging.Logger) *GraphQLClient {
	client := graphql.NewClient(url, nil)

	return &GraphQLClient{
		client: client,
		logger: logger,
	}
}

// To Fetch Registered Opeator Public Key Information from theGraph
func (c *GraphQLClient) GetOperatorInfo(ctx context.Context, opId sdktypes.OperatorId) (sdktypes.OperatorInfo, error) {
	operatorId := hex.EncodeToString(opId[:])

	vars := map[string]interface{}{
		"id": graphql.String(fmt.Sprintf("0x%s", operatorId)),
	}

	var query OperatorQuery

	err := c.client.Query(ctx, &query, vars)

	if err != nil {
		c.logger.Error("Error while fetching operator info from graph", "err", err)
		return sdktypes.OperatorInfo{}, err
	}

	parsed, err := ConvertIndexedOperatorInfoGqlToOperatorInfo(&query)

	return *parsed, err
}

func (c *GraphQLClient) GetAgentPodsForAgent(ctx context.Context, agent string) (AgentPodsInfo, error) {
	// Define the query variables
	vars := map[string]interface{}{
		"id": graphql.String(strings.ToLower(agent)),
	}

	var query AgentPodMappingQuery

	err := c.client.Query(ctx, &query, vars)

	if err != nil {
		c.logger.Error("Error while fetching operator info from graph", "err", err)
		return AgentPodsInfo{}, err
	}

	pods, err := ConvertAgentPodsInfo(&query)

	return pods, err
}
