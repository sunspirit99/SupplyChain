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

type Diary struct {
	Id               string `json:id gorm:"PRIMARY_KEY"`
	Title            string `json:title`
	Content          string `json:content`
	Created_by       string `json:created_by`
	Created_time     string `json:created_time`
	Related_area     string `json:related_area`
	Related_batch    string `json:related_batch`
	Related_products string `json:related_products`
}

func (s *SmartContract) InitDiary(ctx contractapi.TransactionContextInterface) error {
	diaries := []Diary{
		{Id: "1", Title: "tuoi nuoc", Content: "3 lan moi lan 30 giay", Created_by: "linh1", Created_time: "07/12/2021", Related_area: "Ha Noi", Related_batch: "1", Related_products: "Hat giong 1"},
		{Id: "2", Title: "bon phan", Content: "1kg / cay", Created_by: "linh2", Created_time: "07/12/2021", Related_area: "TP HCM", Related_batch: "2", Related_products: "Hat giong 2"},
		{Id: "3", Title: "thu hoach", Content: "thu hoach qua chin", Created_by: "linh3", Created_time: "07/12/2021", Related_area: "Nam Dinh", Related_batch: "3", Related_products: "Hat giong 3"},
	}

	fmt.Println("Init Diary ..............")

	for _, diary := range diaries {
		diaryJSON, err := json.Marshal(diary)
		if err != nil {
			fmt.Println(diary)
			return err
		}

		err = ctx.GetStub().PutState(diary.Id, diaryJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}

func (s *SmartContract) AddAction(ctx contractapi.TransactionContextInterface, id string, title string, content string, created_by string, created_time string, related_area string, related_batch string, related_products string) error {
	exists, err := s.ActionExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the diary %s already exists", id)
	}

	diary := Diary{
		Id:               id,
		Title:            title,
		Content:          content,
		Created_by:       created_by,
		Created_time:     created_time,
		Related_area:     related_area,
		Related_batch:    related_batch,
		Related_products: related_products,
	}
	DiaryJSON, err := json.Marshal(diary)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, DiaryJSON)
}

func (s *SmartContract) GetAllDiaries(ctx contractapi.TransactionContextInterface) ([]*Diary, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all Accounts in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var diaries []*Diary
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var diary Diary
		err = json.Unmarshal(queryResponse.Value, &diary)
		if err != nil {
			return nil, err
		}
		diaries = append(diaries, &diary)
	}

	return diaries, nil
}

func (s *SmartContract) ActionExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	actionJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return actionJSON != nil, nil
}

// func (s *SmartContract) GetActionDetail(ctx contractapi.TransactionContextInterface, title string) ([]*Diary, error) {
// 	queryString := fmt.Sprintf(`{"selector":{"docType":"diary","title":"%s"}}`, title)
// 	return getQueryResultForQueryString(ctx, queryString)
// }

// func (s *SmartContract) GetDiaryByProduct(ctx contractapi.TransactionContextInterface, productID string) ([]*Diary, error) {
// 	queryString := fmt.Sprintf(`{"selector":{"docType":"diary","related_products":"%s"}}`, productID)
// 	return getQueryResultForQueryString(ctx, queryString)
// }

// func (s *SmartContract) GetAllDiaries(ctx contractapi.TransactionContextInterface) ([]*Diary, error) {
// 	// range query with empty string for startKey and endKey does an
// 	// open-ended query of all Accounts in the chaincode namespace.
// 	resultsIterator, err := ctx.GetStub().GetQueryResult(fmt.Sprintf(`{"selector":{"docType":"diary"}}`))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	var diaries []*Diary
// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}

// 		var diary Diary
// 		err = json.Unmarshal(queryResponse.Value, &diary)
// 		if err != nil {
// 			return nil, err
// 		}
// 		diaries = append(diaries, &diary)
// 	}

// 	return diaries, nil
// }

// func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Diary, error) {

// 	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
// 	if err != nil {
// 		fmt.Println("Failed to get Query Result !")
// 	}
// 	defer resultsIterator.Close()

// 	var diaries []*Diary

// 	for resultsIterator.HasNext() {
// 		queryResult, err := resultsIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}
// 		var diary Diary
// 		err = json.Unmarshal(queryResult.Value, &diary)
// 		if err != nil {
// 			return nil, err
// 		}
// 		diaries = append(diaries, &diary)
// 	}

// 	return diaries, nil
// }
