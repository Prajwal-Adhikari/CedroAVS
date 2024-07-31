package graphql

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	// "time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shurcooL/graphql"
)

type GraphQLClient struct {
	client *graphql.Client
	logger logging.Logger
}


func NewGraphQLClient(url string, logger logging.Logger) *GraphQLClient {
	client := graphql.NewClient(url, nil)

	// querier := NewRetryQuerier(client, time.Second*5, 3)

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

func (c *GraphQLClient) GetAgentPodsChainLists(ctx context.Context, agent string) (AgentPodChainListInfo, error) {
	// Define the query variables
	vars := map[string]interface{}{
		"agent": graphql.String(strings.ToLower(agent)),
	}

	var query AgentPodChainListsQuery

	err := c.client.Query(ctx, &query, vars)

	if err != nil {
		c.logger.Error("Error while fetching operator info from graph", "err", err)
		return AgentPodChainListInfo{}, err
	}

	pods, err := ConvertAgentPodsChains(&query)

	return pods, err
}

func (c *GraphQLClient) GetEnrolledAgentsList(ctx context.Context) ([]common.Address, error) {

	var query EnrolledAgentsListQuery

	err := c.client.Query(ctx, &query, nil)

	if err != nil {
		c.logger.Error("Error while fetching operator info from graph", "err", err)
		return []common.Address{}, err
	}

	pods, err := ConvertEnrolledAgent(&query)

	return pods, err
}

func (c *GraphQLClient) GetAgentSummaries(ctx context.Context, agent string) (AgentSummary, error) {

	vars := map[string]interface{}{
		"id": graphql.String(strings.ToLower(agent)),
	}

	var query AgentSummaryQuery

	err := c.client.Query(ctx, &query, vars)

	if err != nil {
		return AgentSummary{}, err
	}

	summary, err := ConvertAgentSummaryInfo(&query)

	return summary, err
}

func (c *GraphQLClient) GetUnclaimedDepositsLists(ctx context.Context, agent string, timestamp uint64) ([]DepositInfo, error) {

	summary, err := c.GetAgentSummaries(ctx, agent)

	if err != nil {
		c.logger.Error("Error while fetching info from graph", "err", err)

		return []DepositInfo{}, err
	}

	totalUnclaimedDeposits := summary.TotalDepositCount - summary.TotalClaimCount

	deposits := make([]DeposisInfoQuery, 0)

	for {
		vars := map[string]interface{}{
			"agentId":   graphql.String(strings.ToLower(agent)),
			"timestamp": graphql.String(fmt.Sprint(timestamp)),
			"skip":      graphql.Int(len(deposits)),
			"first":     graphql.Int(1000),
		}

		var query UnclaimedDepositListQuery

		err = c.client.Query(ctx, &query, vars)

		if err != nil {
			c.logger.Error("Error while fetching info from graph", "err", err)
			return []DepositInfo{}, err
		}

		if len(deposits) == totalUnclaimedDeposits {
			break
		}

		deposits = append(deposits, query.Deposits...)
	}

	pDeps, err := ConvertDepositListInfo(deposits)

	return pDeps, err
}
