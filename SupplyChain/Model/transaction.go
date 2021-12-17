package Model

import (
	"gorm.io/gorm"
)

type TxORM struct {
	Tx
	db *gorm.DB
}

func NewTxORM() *TxORM {
	return &TxORM{db: Connector}
}

func (tx TxORM) CreateTx(tran Tx) (Tx, error) {

	err := tx.db.Create(&tran).Error
	return tran, err
}

func (tx TxORM) UpdateTxID(status int, txid string, trace string) (Tx, error) {
	var tran Tx
	err := tx.db.Exec("UPDATE txes SET status = ?, tx_id = ? WHERE trace = ?", status, txid, trace).Error // status = 2 : Failed
	return tran, err
}
