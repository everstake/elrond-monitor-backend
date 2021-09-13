package smodels

import "github.com/shopspring/decimal"

type RankingRange struct {
	Amount decimal.Decimal `json:"amount"`
	Count  uint64          `json:"count"`
}

type Ranking struct {
	Name      string       `json:"name"`
	Address   string       `json:"address"`
	T100      RankingRange `json:"t_100"`
	F100T1k   RankingRange `json:"f_100_t_1k"`
	F1kT10k   RankingRange `json:"f_1k_t_10k"`
	F10kT100k RankingRange `json:"f_10k_t_100k"`
	F100k     RankingRange `json:"f_100k"`
}
