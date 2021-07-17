package node

import (
	"github.com/shopspring/decimal"
)

type (
	ChainStatus struct {
		Status struct {
			ErdHighestFinalNonce uint64 `json:"erd_highest_final_nonce"`
		} `json:"status"`
	}

	HyperBlock struct {
		HyperBlock struct {
			// block nonce is a block height
			Nonce       uint64 `json:"nonce"`
			Hash        string `json:"hash"`
			Shardblocks []struct {
				Hash  string `json:"hash"`
				Nonce uint64 `json:"nonce"`
				Shard uint64 `json:"shard"`
			} `json:"shardBlocks"`
			Transactions []Tx `json:"transactions"`
		} `json:"hyperblock"`
	}

	Block struct {
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
				DestinationShard uint64 `json:"destinationShard"`
				Hash             string `json:"hash"`
				SourceShard      uint64 `json:"sourceShard"`
				Type             string `json:"type"`
			} `json:"miniBlocks"`
		} `json:"block"`
	}

	BlockExtraData struct {
		Epoch         uint64   `json:"epoch"`
		Hash          string   `json:"hash"`
		Nonce         uint64   `json:"nonce"`
		TxCount       uint64   `json:"txCount"`
		Shard         uint64   `json:"shard"`
		Timestamp     int64    `json:"timestamp"`
		Round         uint64   `json:"round"`
		PrevHash      string   `json:"prevHash"`
		Proposer      string   `json:"proposer"`
		SizeTxs       uint64   `json:"sizeTxs"`
		Size          uint64   `json:"size"`
		StateRootHash string   `json:"stateRootHash"`
		Validators    []string `json:"validators"`
		PubKeyBitmap  string   `json:"pubKeyBitmap"`
	}

	MiniBlock struct {
		MiniBlockHash     string `json:"miniBlockHash"`
		ReceiverBlockHash string `json:"receiverBlockHash"`
		ReceiverShard     uint64 `json:"receiverShard"`
		SenderBlockHash   string `json:"senderBlockHash"`
		SenderShard       uint64 `json:"senderShard"`
		Timestamp         int64  `json:"timestamp"`
		Type              string `json:"type"`
		Error             string `json:"error"`
	}

	Tx struct {
		Data          string     `json:"data"`
		Fee           int64      `json:"fee,string"`
		GasLimit      uint64     `json:"gasLimit"`
		GasPrice      uint64     `json:"gasPrice"`
		GasUsed       uint64     `json:"gasUsed"`
		MiniBlockHash string     `json:"miniBlockHash"`
		Nonce         uint64     `json:"nonce"`
		Receiver      string     `json:"receiver"`
		ReceiverShard uint64     `json:"receiverShard"`
		Round         uint64     `json:"round"`
		ScResults     []ScResult `json:"scResults"`
		Sender        string     `json:"sender"`
		SenderShard   uint64     `json:"senderShard"`
		Signature     string     `json:"signature"`
		Status        string     `json:"status"`
		Timestamp     int64      `json:"timestamp"`
		Txhash        string     `json:"txHash"`
		Value         string     `json:"value"`
	}

	ScResult struct {
		CallType       string `json:"callType"`
		Data           string `json:"data"`
		GasLimit       uint64 `json:"gasLimit"`
		GasPrice       uint64 `json:"gasPrice"`
		Hash           string `json:"hash"`
		Nonce          uint64 `json:"nonce"`
		OriginalTxHash string `json:"originalTxHash"`
		PrevTxHash     string `json:"prevTxHash"`
		Receiver       string `json:"receiver"`
		RelayedValue   string `json:"relayedValue"`
		Sender         string `json:"sender"`
		Value          string `json:"value"`
		ReturnMessage  string `json:"returnMessage"`
	}

	TxDetails struct {
		Hash      string          `json:"hash"`
		GasLimit  uint64          `json:"gasLimit"`
		GasUsed   uint64          `json:"gasUsed"`
		Data      string          `json:"data"`
		Value     decimal.Decimal `json:"value,string"`
		ScResults []struct {
			Receiver string          `json:"receiver"`
			Data     string          `json:"data"`
			Sender   string          `json:"sender"`
			Value    decimal.Decimal `json:"value,string"`
		} `json:"scResults"`
		Timestamp int64  `json:"timestamp"`
		Status    string `json:"status"`
		Error     string `json:"error"`
	}

	Address struct {
		Account struct {
			Address  string          `json:"address"`
			Nonce    int64           `json:"nonce"`
			Balance  decimal.Decimal `json:"balance"`
			Username string          `json:"username"`
		} `json:"account"`
	}

	Stats struct {
		Shards         uint64 `json:"shards"`
		Blocks         uint64 `json:"blocks"`
		Accounts       uint64 `json:"accounts"`
		Transactions   uint64 `json:"transactions"`
		RefreshRates   uint64 `json:"refreshRate"`
		Epoch          uint64 `json:"epoch"`
		RoundsPerEpoch uint64 `json:"roundsPerEpoch"`
	}

	Identity struct {
		Identity     string                 `json:"identity"`
		Name         string                 `json:"name"`
		Description  string                 `json:"description"`
		Avatar       string                 `json:"avatar"`
		Score        uint64                 `json:"score"`
		Validators   uint64                 `json:"validators"`
		Stake        string                 `json:"stake"`
		TopUp        string                 `json:"topUp"`
		Locked       string                 `json:"locked"`
		Distribution map[string]interface{} `json:"distribution"`
		Providers    []string               `json:"providers"`
		StakePercent decimal.Decimal        `json:"stakePercent"`
		Rank         uint64                 `json:"rank"`
	}

	Economics struct {
		TotalSupply       uint64 `json:"totalSupply"`
		CirculatingSupply uint64 `json:"circulatingSupply"`
		Staked            uint64 `json:"staked"`
	}

	NetworkStatus struct {
		Status struct {
			ErdCurrentRound               int `json:"erd_current_round"`
			ErdEpochNumber                int `json:"erd_epoch_number"`
			ErdHighestFinalNonce          int `json:"erd_highest_final_nonce"`
			ErdNonce                      int `json:"erd_nonce"`
			ErdNonceAtEpochStart          int `json:"erd_nonce_at_epoch_start"`
			ErdNoncesPassedInCurrentEpoch int `json:"erd_nonces_passed_in_current_epoch"`
			ErdRoundAtEpochStart          int `json:"erd_round_at_epoch_start"`
			ErdRoundsPassedInCurrentEpoch int `json:"erd_rounds_passed_in_current_epoch"`
			ErdRoundsPerEpoch             int `json:"erd_rounds_per_epoch"`
		} `json:"status"`
	}
)
