package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Area struct {
	Id           string  `json:id gorm:"PRIMARY_KEY"`
	Owner_id     string  `json:owner_id`
	Acreage      float32 `json:acreage`
	Acreage_unit string  `json:acreage_unit`
	Seed_id      string  `json:seed_id`
	Info         string  `json:info`
	Address      string  `json:address`
}

func (s *SmartContract) InitAreas(ctx contractapi.TransactionContextInterface) error {
	areas := []Area{
		{Id: "1", Owner_id: "Linh1", Acreage: 100, Acreage_unit: "ha", Seed_id: "Hat giong 1", Info: "Xuat su Han Quoc", Address: "Ha Noi"},
		{Id: "2", Owner_id: "Linh2", Acreage: 150, Acreage_unit: "ha", Seed_id: "Hat giong 2", Info: "Xuat su Viet Nam", Address: "TP HCM"},
		{Id: "3", Owner_id: "Linh3", Acreage: 200, Acreage_unit: "ha", Seed_id: "Hat giong 3", Info: "Xuat su Trung Quoc", Address: "Nam Dinh"},
	}

	fmt.Println("Init Areas ..............")

	for _, area := range areas {
		areaJSON, err := json.Marshal(area)
		if err != nil {
			fmt.Println(area)
			return err
		}

		err = ctx.GetStub().PutState(area.Id, areaJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}

func (s *SmartContract) GetAllAreas(ctx contractapi.TransactionContextInterface) ([]*Area, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all Accounts in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var areas []*Area
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var area Area
		err = json.Unmarshal(queryResponse.Value, &area)
		if err != nil {
			return nil, err
		}
		areas = append(areas, &area)
	}

	return areas, nil
}
