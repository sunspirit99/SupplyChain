package Type

import (
	model "SuperBank/Model"
	"SuperBank/database"

	"gorm.io/gorm"
)

type TransactionsORM struct {
	model.Transactions
	db *gorm.DB
}

func NewTransactionsORM() *TransactionsORM {
	return &TransactionsORM{db: database.Connector}
}

func (tx TransactionsORM) CreateTx(tran model.Transactions) (model.Transactions, error) {
	err := tx.db.Create(&tran).Error
	return tran, err
}
