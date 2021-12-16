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

type Transaction struct {
	Id            string `json:id gorm:"PRIMARY_KEY"`
	Transport_id  string `json:transport_id`
	Related_batch string `json:related_batch`
	Pickup_place  string `json:pickup_place`
	Deliver_place string `json:deliver_place`
	Created_time  string `json:created_time`
	Description   string `json:description`
}

func (s *SmartContract) InitTransaction(ctx contractapi.TransactionContextInterface) error {
	trans := []Transaction{
		{Id: "1", Transport_id: "001", Related_batch: "1", Pickup_place: "Ha Noi", Deliver_place: "Bac Ninh", Created_time: "07/12/2021", Description: "nothing"},
		{Id: "2", Transport_id: "002", Related_batch: "2", Pickup_place: "TP HCM", Deliver_place: "Hai Duong", Created_time: "07/12/2021", Description: "nothing"},
		{Id: "3", Transport_id: "003", Related_batch: "3", Pickup_place: "Nam Dinh", Deliver_place: "Dong Nai", Created_time: "07/12/2021", Description: "nothing"},
	}

	fmt.Println("Init Transactions ..............")

	for _, tran := range trans {
		tranJSON, err := json.Marshal(tran)
		if err != nil {
			fmt.Println(tran)
			return err
		}

		err = ctx.GetStub().PutState(tran.Id, tranJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}

func (s *SmartContract) AddTransaction(ctx contractapi.TransactionContextInterface, id string, transport_id string, related_batch string, pickup_place string, deliver_place string, created_time string, description string) error {
	exists, err := s.TransactionExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the Transaction %s already exists", id)
	}

	tran := Transaction{
		Id:            id,
		Transport_id:  transport_id,
		Related_batch: related_batch,
		Pickup_place:  pickup_place,
		Deliver_place: deliver_place,
		Created_time:  created_time,
		Description:   description,
	}
	TranJSON, err := json.Marshal(tran)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, TranJSON)
}

func (s *SmartContract) GetAllTransactions(ctx contractapi.TransactionContextInterface) ([]*Transaction, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all Accounts in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var trans []*Transaction
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var tran Transaction
		err = json.Unmarshal(queryResponse.Value, &tran)
		if err != nil {
			return nil, err
		}
		trans = append(trans, &tran)
	}

	return trans, nil
}

func (s *SmartContract) TransactionExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	TransactionJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return TransactionJSON != nil, nil
}
