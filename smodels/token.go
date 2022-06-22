package smodels

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

type (
	NFT struct {
		Name       string          `json:"name"`
		Identity   string          `json:"identity"`
		Owner      string          `json:"owner"`
		Creator    string          `json:"creator"`
		Collection string          `json:"collection"`
		Type       string          `json:"type"`
		Minted     Time            `json:"minted"`
		Royalties  decimal.Decimal `json:"royalties"`
		Assets     json.RawMessage `json:"assets"`
	}
	Token struct {
		Identity   string          `json:"identity"`
		Name       string          `json:"name"`
		Type       string          `json:"type"`
		Owner      string          `json:"owner"`
		Supply     decimal.Decimal `json:"supply"`
		Decimals   uint64          `json:"decimals"`
		Properties json.RawMessage `json:"properties"`
		Roles      json.RawMessage `json:"roles"`
	}
	NFTCollection struct {
		Name       string          `json:"name"`
		Identity   string          `json:"identity"`
		Owner      string          `json:"owner"`
		Type       string          `json:"type"`
		Properties json.RawMessage `json:"properties"`
		CreatedAt  Time            `json:"created_at"`
	}
)
