package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const index = "product~id"

// SmartContract provides functions for managing an Account
type SmartContract struct {
	contractapi.Contract
}

type Diary struct {
	Id               string `json:id`
	Title            string `json:title`
	Content          string `json:content`
	Created_by       string `json:created_by`
	Created_time     string `json:created_time`
	Related_area     string `json:related_area`
	Related_batch    string `json:related_batch`
	Related_products string `json:relatedproducts`
}

func (s *SmartContract) InitDiary(ctx contractapi.TransactionContextInterface) error {
	diaries := []Diary{
		{Id: "1", Title: "tuoi nuoc", Content: "3 lan moi lan 30 giay", Created_by: "linh1", Created_time: "2021-12-16T11:45:26.371Z", Related_area: "Ha Noi", Related_batch: "1", Related_products: "xoai"},
		{Id: "2", Title: "bon phan", Content: "1kg / cay", Created_by: "linh2", Created_time: "2021-12-16T11:45:26.371Z", Related_area: "TP HCM", Related_batch: "2", Related_products: "buoi"},
		{Id: "3", Title: "thu hoach", Content: "thu hoach qua chin", Created_by: "linh3", Created_time: "2021-12-16T11:45:26.371Z", Related_area: "Nam Dinh", Related_batch: "3", Related_products: "cam"},
	}

	fmt.Println("Init Diary ..............")

	for _, diary := range diaries {

		err := s.AddAction(ctx, diary.Id, diary.Title, diary.Content, diary.Created_by, diary.Created_time, diary.Related_area, diary.Related_batch, diary.Related_products)
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

	err = ctx.GetStub().PutState(id, DiaryJSON)
	if err != nil {
		return err
	}

	productIDIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{diary.Related_products, diary.Id})
	if err != nil {
		return err
	}

	value := []byte{0x00}
	return ctx.GetStub().PutState(productIDIndexKey, value)
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

func (s *SmartContract) GetActionDetail(ctx contractapi.TransactionContextInterface, id string) (*Diary, error) {

	diaryJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if diaryJSON == nil {
		return nil, fmt.Errorf("the diary %s does not exist", id)
	}

	var diary Diary
	err = json.Unmarshal(diaryJSON, &diary)
	if err != nil {
		return nil, err
	}

	return &diary, nil
}

func (s *SmartContract) GetDiaryByProduct(ctx contractapi.TransactionContextInterface, productID string) ([]*Diary, error) {

	ResultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(index, []string{productID})
	if err != nil {
		return nil, err
	}
	defer ResultsIterator.Close()

	var diaries []*Diary
	for ResultsIterator.HasNext() {
		responseRange, err := ResultsIterator.Next()
		if err != nil {
			return nil, err
		}

		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return nil, err
		}

		if len(compositeKeyParts) > 1 {
			returnedProductID := compositeKeyParts[1]

			diary, err := s.GetActionDetail(ctx, returnedProductID)
			if err != nil {
				return nil, err
			}

			diaries = append(diaries, diary)
		}
	}

	return diaries, nil
}

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
