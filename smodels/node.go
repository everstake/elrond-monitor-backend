package smodels

import (
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/shopspring/decimal"
)

const (
	NodeTypeObserver  = "observer"
	NodeTypeValidator = "validator"
)

type Node struct {
	node.HeartbeatStatus
	node.ValidatorStatistic
	Type     string          `json:"type"`
	Status   string          `json:"status"`
	UpTime   float64         `json:"upTime"`
	DownTime float64         `json:"downTime"`
	Owner    string          `json:"owner"`
	Provider string          `json:"provider"`
	Stake    decimal.Decimal `json:"stake"`
	TopUp    decimal.Decimal `json:"topUp"`
	Locked   decimal.Decimal `json:"locked"`
}
