package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/smodels"
)

func (s *ServiceFacade) GetTransactions(filter filters.Transactions) (items smodels.Pagination, err error) {
	dTxs, err := s.dao.GetTransactions(filter)
	if err != nil {
		return items, fmt.Errorf("dao.GetTransactions: %s", err.Error())
	}
	txs := make([]smodels.Tx, len(dTxs))
	for i, tx := range dTxs {
		txs[i] = smodels.Tx{
			Hash:          tx.Hash,
			Status:        tx.Status,
			From:          tx.Sender,
			To:            tx.Receiver,
			Value:         tx.Value,
			Fee:           tx.Fee,
			GasUsed:       tx.GasUsed,
			MiniblockHash: tx.MiniBlockHash,
			ShardFrom:     tx.SenderShard,
			ShardTo:       tx.ReceiverShard,
			Type:          "", // todo
			Timestamp:     smodels.NewTime(tx.CreatedAt),
		}
	}
	total, err := s.dao.GetTransactionsTotal(filter)
	if err != nil {
		return items, fmt.Errorf("dao.GetTransactionsTotal: %s", err.Error())
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
		results[i] = smodels.ScResult{
			Hash:    r.Hash,
			From:    r.From,
			To:      r.To,
			Value:   r.Value,
			Data:    r.Data,
			Message: r.Message,
		}
	}
	return smodels.Tx{
		Hash:          dTx.Hash,
		Status:        dTx.Status,
		From:          dTx.Sender,
		To:            dTx.Receiver,
		Value:         dTx.Value,
		Fee:           dTx.Fee,
		GasUsed:       dTx.GasUsed,
		MiniblockHash: dTx.MiniBlockHash,
		ShardFrom:     dTx.SenderShard,
		ShardTo:       dTx.ReceiverShard,
		Type:          "", // todo
		ScResults:     results,
		Timestamp:     smodels.NewTime(dTx.CreatedAt),
	}, nil
}
