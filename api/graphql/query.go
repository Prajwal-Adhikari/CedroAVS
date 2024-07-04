package graphql

import "github.com/shurcooL/graphql"

type OperatorQuery struct {
	Operator struct {
		Operator   graphql.String
		PubkeyG1_X graphql.String   `graphql:"pubkeyG1_X"`
		PubkeyG1_Y graphql.String   `graphql:"pubkeyG1_Y"`
		PubkeyG2_X []graphql.String `graphql:"pubkeyG2_X"`
		PubkeyG2_Y []graphql.String `graphql:"pubkeyG2_Y"`
		ID         graphql.String   `graphql:"id"`
	} `graphql:"operator(id: $id)"`
}

type EnrolledAgentsListQuery struct {
	AgendPodChains []struct {
		Agent graphql.String `graphql:"agentWallet"`
	} `graphql:"agentPodChainLists(orderBy: blockTimestamp, orderDirection: desc)"`
}

type AgentPodChainListsQuery struct {
	AgendPodChains []struct {
		ChainIDs []graphql.String `graphql:"chainIDs"`
		Agent    graphql.String   `graphql:"agentWallet"`
	} `graphql:"agentPodChainLists(where: {agentWallet: $agent}, orderBy: blockTimestamp, orderDirection: desc)"`
}

type DeposisInfoQuery struct {
	IsClaimed       graphql.Boolean `graphql:"isClaimed"`
	InterestRate    graphql.String     `graphql:"interestRate"`
	TxnID           graphql.String  `graphql:"txnID"`
	PrincipalAmount graphql.String     `graphql:"principalAmount"`
	BlockTimestamp  graphql.String     `graphql:"blockTimestamp"`
	TimeInterval    graphql.String     `graphql:"timeInterval"`
	Agent           graphql.String  `graphql:"agentWallet"`
	AgentContract   graphql.String  `graphql:"agentPodContract"`
	BlockNumber     graphql.String     `graphql:"blockNumber"`
}

type AgentSummaryQuery struct {
	AgentSummary struct {
		AgentContract     graphql.String `graphql:"agentPodContract"`
		Agent             graphql.String `graphql:"agentWallet"`
		TotalDepositCount graphql.String `graphql:"totalDepositCount"`
		TotalClaimCount   graphql.String `graphql:"totalClaimCount"`
		ChainId           graphql.String `graphql:"chainID"`
	} `graphql:"agentSummary(id:$id)"`
}

type UnclaimedDepositListQuery struct {
	Deposits []DeposisInfoQuery `graphql:"deposits(where: {agentWallet: $agentId, isClaimed:false, blockTimestamp_lte:$timestamp}, from: $from, skip:$skip)"`
}
