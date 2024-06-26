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

type AgentPodCreatedQuery struct {
	AgentPods []struct {
		Agent         graphql.String `graphql:"agent"`
		AgentContract graphql.String `graphql:"agentContract"`
		ChainId       graphql.String `graphql:"chainID"`
	} `graphql:"agentPodCreateds(where: {agent_contains: $agent})"`
}

type AgentPodMappingQuery struct {
	AgentPods struct {
		Agent         graphql.String `graphql:"id"`
		AgentContract graphql.String `graphql:"agentContract"`
		ChainId       graphql.String `graphql:"chainID"`
	} `graphql:"agentOwnerMapping(id:$id)"`
}

type AgentPodChainListsQuery struct {
	AgendPodChains []struct {
		ChainIDs []graphql.String `graphql:"chainIDs"`
		Agent    graphql.String   `graphql:"agent"`
	} `graphql:"agentPodChainLists(where: {agent: $agent}, orderBy: blockTimestamp, orderDirection: desc)"`
}

type DeposisInfoQuery struct {
	IsClaimed       graphql.Boolean `graphql:"isClaimed"`
	InterestRate    graphql.String  `graphql:"interestRate"`
	TxnID           graphql.String  `graphql:"txnID"`
	PrincipalAmount graphql.String  `graphql:"principalAmount"`
	BlockTimestamp  graphql.String  `graphql:"blockTimestamp"`
	TimeInterval    graphql.String  `graphql:"timeInterval"`
	AgentID         graphql.String  `graphql:"agentId"`
}

type UnclaimedDepositListQuery struct {
	Deposits []DeposisInfoQuery `graphql:"deposits(where: {agentId: $agentId, isClaimed:false})"`
}