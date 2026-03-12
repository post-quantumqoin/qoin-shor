package api

import apitypes "github.com/post-quantumqoin/qoin-shor/api/types"

func CreateEthRPCAliases(as apitypes.Aliaser) {
	// TODO: maybe use reflect to automatically register all the eth aliases
	as.AliasMethod("eth_accounts", "Qoin.EthAccounts")
	as.AliasMethod("eth_blockNumber", "Qoin.EthBlockNumber")
	as.AliasMethod("eth_getBlockTransactionCountByNumber", "Qoin.EthGetBlockTransactionCountByNumber")
	as.AliasMethod("eth_getBlockTransactionCountByHash", "Qoin.EthGetBlockTransactionCountByHash")

	as.AliasMethod("eth_getBlockByHash", "Qoin.EthGetBlockByHash")
	as.AliasMethod("eth_getBlockByNumber", "Qoin.EthGetBlockByNumber")
	as.AliasMethod("eth_getTransactionByHash", "Qoin.EthGetTransactionByHash")
	as.AliasMethod("eth_getTransactionCount", "Qoin.EthGetTransactionCount")
	as.AliasMethod("eth_getTransactionReceipt", "Qoin.EthGetTransactionReceipt")
	as.AliasMethod("eth_getTransactionByBlockHashAndIndex", "Qoin.EthGetTransactionByBlockHashAndIndex")
	as.AliasMethod("eth_getTransactionByBlockNumberAndIndex", "Qoin.EthGetTransactionByBlockNumberAndIndex")

	as.AliasMethod("eth_getCode", "Qoin.EthGetCode")
	as.AliasMethod("eth_getStorageAt", "Qoin.EthGetStorageAt")
	as.AliasMethod("eth_getBalance", "Qoin.EthGetBalance")
	as.AliasMethod("eth_chainId", "Qoin.EthChainId")
	as.AliasMethod("eth_syncing", "Qoin.EthSyncing")
	as.AliasMethod("eth_feeHistory", "Qoin.EthFeeHistory")
	as.AliasMethod("eth_protocolVersion", "Qoin.EthProtocolVersion")
	as.AliasMethod("eth_maxPriorityFeePerGas", "Qoin.EthMaxPriorityFeePerGas")
	as.AliasMethod("eth_gasPrice", "Qoin.EthGasPrice")
	as.AliasMethod("eth_sendRawTransaction", "Qoin.EthSendRawTransaction")
	as.AliasMethod("eth_getMessageCid", "Qoin.EthGetMessageCid")
	as.AliasMethod("eth_estimateGas", "Qoin.EthEstimateGas")
	as.AliasMethod("eth_call", "Qoin.EthCall")

	as.AliasMethod("eth_getLogs", "Qoin.EthGetLogs")
	as.AliasMethod("eth_getFilterChanges", "Qoin.EthGetFilterChanges")
	as.AliasMethod("eth_getFilterLogs", "Qoin.EthGetFilterLogs")
	as.AliasMethod("eth_newFilter", "Qoin.EthNewFilter")
	as.AliasMethod("eth_newBlockFilter", "Qoin.EthNewBlockFilter")
	as.AliasMethod("eth_newPendingTransactionFilter", "Qoin.EthNewPendingTransactionFilter")
	as.AliasMethod("eth_uninstallFilter", "Qoin.EthUninstallFilter")
	as.AliasMethod("eth_subscribe", "Qoin.EthSubscribe")
	as.AliasMethod("eth_unsubscribe", "Qoin.EthUnsubscribe")

	as.AliasMethod("trace_block", "Qoin.EthTraceBlock")
	as.AliasMethod("trace_replayBlockTransactions", "Qoin.EthTraceReplayBlockTransactions")

	as.AliasMethod("net_version", "Qoin.NetVersion")
	as.AliasMethod("net_listening", "Qoin.NetListening")

	as.AliasMethod("web3_clientVersion", "Qoin.Web3ClientVersion")
}
