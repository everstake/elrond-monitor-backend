package node

import (
	"github.com/shopspring/decimal"
	"time"
)

type (
	HyperBlock struct {
		// block nonce is a block height
		Nonce       uint64 `json:"nonce"`
		Hash        string `json:"hash"`
		Shardblocks []struct {
			Hash  string `json:"hash"`
			Nonce uint64 `json:"nonce"`
			Shard uint64 `json:"shard"`
		} `json:"shardBlocks"`
		Transactions []BlockTx `json:"transactions"`
	}

	Block struct {
		// block nonce is a block height
		AccumulatedFees int64  `json:"accumulatedFees,string"`
		DeveloperFees   int64  `json:"developerFees,string"`
		Epoch           uint64 `json:"epoch"`
		Hash            string `json:"hash"`
		Nonce           uint64 `json:"nonce"`
		NumTxs          uint64 `json:"numTxs"`
		Shard           uint64 `json:"shard"`
		Timestamp       int64  `json:"timestamp"`
		Round           uint64 `json:"round"`
		PrevBlockHash   string `json:"prevBlockHash"`
		Status          string `json:"status"`
		Miniblocks      []struct {
			DestinationShard uint64    `json:"destinationShard"`
			Hash             string    `json:"hash"`
			SourceShard      uint64    `json:"sourceShard"`
			Type             string    `json:"type"`
			Transactions     []BlockTx `json:"transactions"`
		} `json:"miniBlocks"`
	}

	BlockTx struct {
		Hash string `json:"hash"`
	}

	Tx struct {
		Type                              string                `json:"type"`
		Nonce                             uint64                `json:"nonce"`
		Round                             uint64                `json:"round"`
		Epoch                             uint64                `json:"epoch"`
		Value                             string                `json:"value"`
		Receiver                          string                `json:"receiver"`
		Sender                            string                `json:"sender"`
		GasPrice                          uint64                `json:"gasPrice"`
		GasLimit                          uint64                `json:"gasLimit"`
		Data                              string                `json:"data"`
		Signature                         string                `json:"signature"`
		SourceShard                       uint64                `json:"sourceShard"`
		DestinationShard                  uint64                `json:"destinationShard"`
		BlockNonce                        uint64                `json:"blockNonce"`
		BlockHash                         string                `json:"blockHash"`
		NotarizedAtSourceInMetaNonce      uint64                `json:"notarizedAtSourceInMetaNonce"`
		NotarizedAtSourceInMetaHash       string                `json:"NotarizedAtSourceInMetaHash"`
		NotarizedAtDestinationInMetaNonce uint64                `json:"notarizedAtDestinationInMetaNonce"`
		NotarizedAtDestinationInMetaHash  string                `json:"notarizedAtDestinationInMetaHash"`
		MiniblockType                     string                `json:"miniblockType"`
		MiniblockHash                     string                `json:"miniblockHash"`
		Status                            string                `json:"status"`
		HyperblockNonce                   uint64                `json:"hyperblockNonce"`
		HyperblockHash                    string                `json:"hyperblockHash"`
		SmartContractResults              []SmartContractResult `json:"smartContractResults"`
	}

	SmartContractResult struct {
		Hash           string          `json:"hash"`
		Nonce          uint64          `json:"nonce"`
		Value          decimal.Decimal `json:"value"`
		Receiver       string          `json:"receiver"`
		Sender         string          `json:"sender"`
		Data           string          `json:"data,omitempty"`
		PrevTxHash     string          `json:"prevTxHash"`
		OriginalTxHash string          `json:"originalTxHash"`
		GasLimit       uint64          `json:"gasLimit"`
		GasPrice       uint64          `json:"gasPrice"`
		CallType       uint64          `json:"callType"`
		OriginalSender string          `json:"originalSender,omitempty"`
		ReturnMessage  string          `json:"returnMessage"`
	}

	Address struct {
		Address  string          `json:"address"`
		Nonce    uint64          `json:"nonce"`
		Balance  decimal.Decimal `json:"balance"`
		Username string          `json:"username"`
	}

	NetworkStatus struct {
		ErdCurrentRound               uint64 `json:"erd_current_round"`
		ErdEpochNumber                uint64 `json:"erd_epoch_number"`
		ErdHighestFinalNonce          uint64 `json:"erd_highest_final_nonce"`
		ErdNonce                      uint64 `json:"erd_nonce"`
		ErdNonceAtEpochStart          uint64 `json:"erd_nonce_at_epoch_start"`
		ErdNoncesPassedInCurrentEpoch uint64 `json:"erd_nonces_passed_in_current_epoch"`
		ErdRoundAtEpochStart          uint64 `json:"erd_round_at_epoch_start"`
		ErdRoundsPassedInCurrentEpoch uint64 `json:"erd_rounds_passed_in_current_epoch"`
		ErdRoundsPerEpoch             uint64 `json:"erd_rounds_per_epoch"`
	}

	Account struct {
		Address  string          `json:"address"`
		Nonce    uint64          `json:"nonce"`
		Balance  decimal.Decimal `json:"balance"`
		Code     string          `json:"code,omitempty"`
		CodeHash string          `json:"code_hash,omitempty"`
		RootHash string          `json:"root_hash,omitempty"`
		TxCount  uint64          `json:"tx_count"`
	}

	ValidatorStatistics map[string]ValidatorStatistic
	ValidatorStatistic  struct {
		TempRating                         float64 `json:"tempRating"`
		NumLeaderSuccess                   int64   `json:"numLeaderSuccess"`
		NumLeaderFailure                   int64   `json:"numLeaderFailure"`
		NumValidatorSuccess                int64   `json:"numValidatorSuccess"`
		NumValidatorFailure                int64   `json:"numValidatorFailure"`
		NumValidatorIgnoredSignatures      int64   `json:"numValidatorIgnoredSignatures"`
		Rating                             float64 `json:"rating"`
		RatingModifier                     float64 `json:"ratingModifier"`
		TotalNumLeaderSuccess              int64   `json:"totalNumLeaderSuccess"`
		TotalNumLeaderFailure              int64   `json:"totalNumLeaderFailure"`
		TotalNumValidatorSuccess           int64   `json:"totalNumValidatorSuccess"`
		TotalNumValidatorFailure           int64   `json:"totalNumValidatorFailure"`
		TotalNumValidatorIgnoredSignatures int64   `json:"totalNumValidatorIgnoredSignatures"`
		ShardID                            uint64  `json:"shardId"`
		ValidatorStatus                    string  `json:"validatorStatus"`
	}

	HeartbeatStatus struct {
		TimeStamp        time.Time `json:"timeStamp"`
		PublicKey        string    `json:"publicKey"`
		VersionNumber    string    `json:"versionNumber"`
		NodeDisplayName  string    `json:"nodeDisplayName"`
		Identity         string    `json:"identity"`
		TotalUpTimeSec   uint64    `json:"totalUpTimeSec"`
		TotalDownTimeSec uint64    `json:"totalDownTimeSec"`
		MaxInactiveTime  string    `json:"maxInactiveTime"`
		ReceivedShardID  uint64    `json:"receivedShardID"`
		ComputedShardID  uint64    `json:"computedShardID"`
		PeerType         string    `json:"peerType"`
		IsActive         bool      `json:"isActive"`
		Nonce            uint64    `json:"nonce"`
		NumInstances     uint64    `json:"numInstances"`
	}

	NetworkEconomics struct {
		ErdDevRewards            decimal.Decimal `json:"erd_dev_rewards"`
		ErdEpochForEconomicsData uint64          `json:"erd_epoch_for_economics_data"`
		ErdInflation             decimal.Decimal `json:"erd_inflation"`
		ErdTotalBaseStakedValue  decimal.Decimal `json:"erd_total_base_staked_value"`
		ErdTotalFees             decimal.Decimal `json:"erd_total_fees"`
		ErdTotalSupply           decimal.Decimal `json:"erd_total_supply"`
		ErdTotalTopUpValue       decimal.Decimal `json:"erd_total_top_up_value"`
	}

	NetworkConfig struct {
		ErdChainID                   string `json:"erd_chain_id"`
		ErdDenomination              int64  `json:"erd_denomination"`
		ErdGasPerDataByte            int64  `json:"erd_gas_per_data_byte"`
		ErdGasPriceModifier          string `json:"erd_gas_price_modifier"`
		ErdLatestTagSoftwareVersion  string `json:"erd_latest_tag_software_version"`
		ErdMetaConsensusGroupSize    int64  `json:"erd_meta_consensus_group_size"`
		ErdMinGasLimit               int64  `json:"erd_min_gas_limit"`
		ErdMinGasPrice               int64  `json:"erd_min_gas_price"`
		ErdMinTransactionVersion     int64  `json:"erd_min_transaction_version"`
		ErdNumMetachainNodes         int64  `json:"erd_num_metachain_nodes"`
		ErdNumNodesInShard           int64  `json:"erd_num_nodes_in_shard"`
		ErdNumShardsWithoutMeta      int64  `json:"erd_num_shards_without_meta"`
		ErdRewardsTopUpGradientPoint string `json:"erd_rewards_top_up_gradient_point"`
		ErdRoundDuration             uint64 `json:"erd_round_duration"`
		ErdRoundsPerEpoch            uint64 `json:"erd_rounds_per_epoch"`
		ErdShardConsensusGroupSize   int64  `json:"erd_shard_consensus_group_size"`
		ErdStartTime                 uint64 `json:"erd_start_time"`
		ErdTopUpFactor               string `json:"erd_top_up_factor"`
	}

	ContractResp struct {
		ReturnData    []string `json:"returnData"`
		ReturnCode    string   `json:"returnCode"`
		ReturnMessage string   `json:"returnMessage"`
	}

	ContractReq struct {
		SCAddress string   `json:"scAddress"`
		FuncName  string   `json:"funcName"`
		Caller    string   `json:"caller,omitempty"`
		Args      []string `json:"args,omitempty"`
	}

	UserStake struct {
		WithdrawOnlyStake    decimal.Decimal
		WaitingStake         decimal.Decimal
		ActiveStake          decimal.Decimal
		UnstakedStake        decimal.Decimal
		DeferredPaymentStake decimal.Decimal
	}

	ProviderConfig struct {
		Owner         string
		ServiceFee    decimal.Decimal
		DelegationCap decimal.Decimal
	}

	ProviderMeta struct {
		Name      string
		Website   string
		Iidentity string
	}

	QueueItem struct {
		BLS      string
		Provider string
		Nonce    uint64
		Position int64
	}

	StakeTopup struct {
		TopUp    decimal.Decimal
		Stake    decimal.Decimal
		Locked   decimal.Decimal
		NumNodes uint64
		Address  string
		Blses    []string
	}
)
