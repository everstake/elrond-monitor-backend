package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/derrors"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"github.com/shopspring/decimal"
	"net/http"
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
		fee, _ := decimal.NewFromString(tx.Fee)
		txs[i] = smodels.Tx{
			Hash:          tx.Hash,
			Status:        tx.Status,
			From:          tx.Sender,
			To:            tx.Receiver,
			Value:         node.ValueToEGLD(val),
			Fee:           node.ValueToEGLD(fee),
			GasUsed:       tx.GasUsed,
			GasPrice:      tx.GasPrice,
			MiniblockHash: tx.MBHash,
			ShardFrom:     uint64(tx.SenderShard),
			ShardTo:       uint64(tx.ReceiverShard),
			ScResults:     nil,
			Signature:     tx.Signature,
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
		if err == derrors.NotFound {
			return tx, smodels.Error{
				Err:      err.Error(),
				Msg:      "transaction not found",
				HttpCode: http.StatusNotFound,
			}
		}
		return tx, fmt.Errorf("dao.GetTransaction: %s", err.Error())
	}
	results := make([]smodels.ScResult, len(dTx.SCResults))
	for i, r := range dTx.SCResults {
		val, _ := decimal.NewFromString(r.Value)
		results[i] = smodels.ScResult{
			Hash:    r.ResultHash,
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
		Hash:          hash,
		Status:        dTx.Status,
		From:          dTx.Sender,
		To:            dTx.Receiver,
		Value:         node.ValueToEGLD(val),
		Fee:           node.ValueToEGLD(fee),
		GasUsed:       dTx.GasUsed,
		GasPrice:      dTx.GasPrice,
		MiniblockHash: dTx.MBHash,
		ShardFrom:     uint64(dTx.SenderShard),
		ShardTo:       uint64(dTx.ReceiverShard),
		ScResults:     results,
		Signature:     dTx.Signature,
		Data:          string(dTx.Data),
		Timestamp:     smodels.NewTime(time.Unix(int64(dTx.Timestamp), 0)),
	}, nil
}
