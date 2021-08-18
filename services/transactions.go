package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"github.com/shopspring/decimal"
	"time"
)

func (s *ServiceFacade) GetTransactions(filter filters.Transactions) (items smodels.Pagination, err error) {
	dTxs, err := s.dao.GetTransactions(filter)
	if err != nil {
		return items, fmt.Errorf("dao.GetTransactions: %s", err.Error())
	}
	txs := make([]smodels.Tx, len(dTxs))
	for i, tx := range dTxs {
		val, _ := decimal.NewFromString(tx.Value)
		txs[i] = smodels.Tx{
			Hash:          tx.Hash,
			Status:        tx.Status,
			From:          tx.Sender,
			To:            tx.Receiver,
			Value:         node.ValueToEGLD(val),
			MiniblockHash: tx.MBHash,
			ShardFrom:     uint64(tx.SenderShard),
			ShardTo:       uint64(tx.ReceiverShard),
			Timestamp:     smodels.NewTime(time.Unix(int64(tx.Timestamp), 0)),
		}
	}
	total, err := s.dao.GetTransactionsCount(filter)
	if err != nil {
		return items, fmt.Errorf("dao.GetTransactionsCount: %s", err.Error())
	}
	return smodels.Pagination{
		Items: txs,
		Count: total,
	}, nil
}

func (s *ServiceFacade) GetTransaction(hash string) (tx smodels.Tx, err error) {
	dTx, err := s.dao.GetTransaction(hash)
	if err != nil {
		return tx, fmt.Errorf("dao.GetTransaction: %s", err.Error())
	}
	dResults, err := s.dao.GetSCResults(dTx.Hash)
	if err != nil {
		return tx, fmt.Errorf("dao.GetSCResults: %s", err.Error())
	}
	results := make([]smodels.ScResult, len(dResults))
	for i, r := range dResults {
		val, _ := decimal.NewFromString(dTx.Value)
		results[i] = smodels.ScResult{
			Hash:    r.Hash,
			From:    r.Sender,
			To:      r.Receiver,
			Value:   node.ValueToEGLD(val),
			Data:    string(r.Data),
			Message: r.ReturnMessage,
		}
	}
	val, _ := decimal.NewFromString(dTx.Value)
	fee, _ := decimal.NewFromString(dTx.Fee)
	return smodels.Tx{
		Hash:          dTx.Hash,
		Status:        dTx.Status,
		From:          dTx.Sender,
		To:            dTx.Receiver,
		Value:         val,
		Fee:           fee,
		GasUsed:       dTx.GasUsed,
		GasPrice:      dTx.GasPrice,
		MiniblockHash: dTx.MBHash,
		ShardFrom:     uint64(dTx.SenderShard),
		ShardTo:       uint64(dTx.ReceiverShard),
		ScResults:     results,
		Signature:     dTx.Signature,
		Timestamp:     smodels.NewTime(time.Unix(int64(dTx.Timestamp), 0)),
	}, nil
}
