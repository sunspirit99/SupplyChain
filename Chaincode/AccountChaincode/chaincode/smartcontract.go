package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Account
type SmartContract struct {
	contractapi.Contract
}

// Account Object
type Account struct {
	Id          string  `json:id gorm:"primary_key"`
	Name        string  `json:name`
	Address     string  `json:address`
	PhoneNumber string  `json:phonenumber`
	Balance     float32 `json:balance`
	Status      int     `json:status`
	Createtime  string  `json:createtime`
}

// Transaction Object
type Transaction struct {
	Trace      string  `json:trace gorm:"PRIMARY_KEY"`
	TxID       string  `json:txid`
	From       string  `json:"from"`
	To         string  `json:"to"`
	Amount     float32 `json:"amount"`
	Status     int     `json:"status"`
	Createtime string  `json:"createtime"`
}

func (s *SmartContract) InitAccount(ctx contractapi.TransactionContextInterface) error {
	accounts := []Account{
		{Id: "1", Name: "Linh1", Address: "Ha Noi", PhoneNumber: "0123456789", Balance: 10000, Status: 1, Createtime: "09/12/2021"},
		{Id: "2", Name: "Linh2", Address: "Ha Noi", PhoneNumber: "0123456789", Balance: 10000, Status: 1, Createtime: "09/12/2021"},
		{Id: "3", Name: "Linh3", Address: "Ha Noi", PhoneNumber: "0123456789", Balance: 10000, Status: 1, Createtime: "09/12/2021"},
	}

	fmt.Println("Init Accounts ..............")

	for _, account := range accounts {
		accountJSON, err := json.Marshal(account)
		if err != nil {
			fmt.Println(account)
			return err
		}

		err = ctx.GetStub().PutState(account.Id, accountJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}

func (s *SmartContract) CreateAccount(ctx contractapi.TransactionContextInterface, id string, name string, address string, phonenumber string, balance float32, status int, createtime string) error {
	exists, err := s.AccountExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the Account %s already exists", id)
	}

	account := Account{
		Id:          id,
		Name:        name,
		Address:     address,
		PhoneNumber: phonenumber,
		Balance:     balance,
		Status:      status,
		Createtime:  createtime,
	}
	AccountJSON, err := json.Marshal(account)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, AccountJSON)
}

func (s *SmartContract) AccountWithdraw(ctx contractapi.TransactionContextInterface, id string, amount float32) (*Account, error) {
	accountJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if accountJSON == nil {
		return nil, fmt.Errorf("the account %s does not exist", id)
	}

	var account Account
	err = json.Unmarshal(accountJSON, &account)
	if err != nil {
		return nil, err
	}

	account.Balance = account.Balance - amount

	accountJSON, err = json.Marshal(account)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(id, accountJSON)

	return &account, err
}

func (s *SmartContract) AccountDeposit(ctx contractapi.TransactionContextInterface, id string, amount float32) (*Account, error) {
	accountJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if accountJSON == nil {
		return nil, fmt.Errorf("the account %s does not exist", id)
	}

	var account Account
	err = json.Unmarshal(accountJSON, &account)
	if err != nil {
		return nil, err
	}

	account.Balance = account.Balance + amount

	accountJSON, err = json.Marshal(account)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(id, accountJSON)

	return &account, err
}

func (s *SmartContract) AccountTransfer(ctx contractapi.TransactionContextInterface, from string, to string, amount float32) ([]*Account, error) {
	senderJSON, err := ctx.GetStub().GetState(from)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if senderJSON == nil {
		return nil, fmt.Errorf("the account %s does not exist", from)
	}

	var sender Account
	err = json.Unmarshal(senderJSON, &sender)
	if err != nil {
		return nil, err
	}

	receiverJSON, err := ctx.GetStub().GetState(to)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if senderJSON == nil {
		return nil, fmt.Errorf("the account %s does not exist", to)
	}

	var receiver Account
	err = json.Unmarshal(receiverJSON, &receiver)
	if err != nil {
		return nil, err
	}

	sender.Balance = sender.Balance - amount
	receiver.Balance = receiver.Balance + amount

	senderJSON, err = json.Marshal(sender)
	if err != nil {
		return nil, err
	}

	receiverJSON, err = json.Marshal(receiver)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(from, senderJSON)
	if err != nil {
		return nil, err
	}
	err = ctx.GetStub().PutState(to, receiverJSON)

	var accounts []*Account
	accounts = append(accounts, &sender, &receiver)

	return accounts, err
}

func (s *SmartContract) AccountExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	AccountJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return AccountJSON != nil, nil
}
