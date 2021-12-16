package Model

import (
	"gorm.io/gorm"
)

type TransactionsORM struct {
	Transactions
	db *gorm.DB
}

func NewTransactionsORM() *TransactionsORM {
	return &TransactionsORM{db: Connector}
}

func (tx TransactionsORM) CreateTx(tran Transactions) (Transactions, error) {
	err := tx.db.Create(&tran).Error
	return tran, err
}
