package dmodels

import (
	"time"
)

const AccountsTable = "accounts"

type Account struct {
	Address   string    `db:"acc_address"`
	CreatedAt time.Time `db:"acc_created_at"`
}
