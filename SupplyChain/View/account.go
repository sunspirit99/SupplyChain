package View

import (
	cfg "SuperBank/Config"
	controller "SuperBank/Controller"
	model "SuperBank/Model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var c = cfg.GetConfig()

var (
	channelName  string = c.GetString("channelName")
	contractName string = c.GetString("contractName")

	maxWorker      = c.GetInt("maxWorker")
	maxTransaction = c.GetInt("maxTransaction")
	maxAccount     = c.GetInt("maxAccount")

	processing = c.GetInt("processing")
	success    = c.GetInt("success")
	fail       = c.GetInt("fail")
)

// GetAllAccount get all Account data
func GetAllAccount(w http.ResponseWriter, r *http.Request) {

	var Accounts []model.Account

	Accounts, err := controller.NewAccountAPI().GetAllAccount()
	if err != nil {
		log.Println(err)
	}
	Response(w, Accounts)

	// Returns a list of Accounts in JSON format
}

// GetAccountByID get Account with specific ID
func GetAccountByID(w http.ResponseWriter, r *http.Request) {

	var Account model.Account

	vars := mux.Vars(r)
	id := vars["id"] // Get id from URL path

	Account, err := controller.NewAccountAPI().GetAccountByID(id)
	if err != nil {
		log.Println(err)
	}
	Response(w, Account)
	// Returns the found Account in JSON format
}

//CreateAccount creates an Account
func CreateAccount(w http.ResponseWriter, r *http.Request) {

	var Account model.Account

	w.Header().Set("Content-Type", "application/json")
	requestBody, _ := ioutil.ReadAll(r.Body) // read JSON data from Body

	json.Unmarshal(requestBody, &Account) // Convert from JSON format to Account Format

	Account, err := controller.NewAccountAPI().CreateAccount(Account)
	if err != nil {
		log.Println(err)
	}
	Response(w, Account)

	// return the created Account in JSON format
}

// UpdateAccountByID updates Account with respective ID, if ID does not exist, create a new Account
func UpdateAccountByID(w http.ResponseWriter, r *http.Request) {

	var Account model.Account

	requestBody, _ := ioutil.ReadAll(r.Body) // read JSON data from Body

	json.Unmarshal(requestBody, &Account) // Convert from JSON format to Account Format

	Account, err := controller.NewAccountAPI().UpdateAccountByID(Account)
	if err != nil {
		log.Println(err)
	}
	Response(w, Account)

	// Return the updated Account in JSON format
}

// DeleteAccountByID delete an Account with specific ID
func DeleteAccountByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r) // Get id from URL path
	id := vars["id"]
	if len(vars) == 0 {
		panic("Enter an ID !")
	}

	err := controller.NewAccountAPI().DeleteAccountByID(id)
	if err != nil {
		log.Println(err)
	}
	NoContent(w)

}

// AccountWithdraw withdraw money from account through a transaction
func AccountWithdraw(w http.ResponseWriter, r *http.Request) {

	var tx model.Tx
	var Account model.Account

	requestBody, err := ioutil.ReadAll(r.Body) // read JSON data from Body
	if err != nil {
		panic("Enter all required information !!!")
	}

	error := json.Unmarshal(requestBody, &tx) // Convert from JSON format to Account Format
	if error != nil {
		fmt.Println("Error :", error)
	}

	Account, err = controller.NewAccountAPI().AccountWithdraw(tx)
	if err != nil {
		log.Println(err)
	}

	Response(w, Account)
	// Returns the new state of the Account when the transaction is done
}

// AccountDeposit deposit money into account through a transaction
func AccountDeposit(w http.ResponseWriter, r *http.Request) {
	var tx model.Tx
	var Account model.Account

	requestBody, err := ioutil.ReadAll(r.Body) // read JSON data from Body
	if err != nil {
		panic("Enter all required information !!!")
	}

	error := json.Unmarshal(requestBody, &tx) // Convert from JSON format to Account Format
	if error != nil {
		fmt.Println("Error :", error)
	}

	Account, err = controller.NewAccountAPI().AccountDeposit(tx)
	if err != nil {
		log.Println(err)
	}

	Response(w, Account)
	// Returns the new state of the Account when the transaction is done
}

// AccountTransfer transfer money between 2 accounts through a transaction
func AccountTransfer(w http.ResponseWriter, r *http.Request) {
	var Accounts []model.Account
	var tx model.Tx

	requestBody, err := ioutil.ReadAll(r.Body) // read JSON data from Body
	if err != nil {
		panic("Enter all required information !!!")
	}

	error := json.Unmarshal(requestBody, &tx) // Convert from JSON format to Account Format
	if error != nil {
		fmt.Println("Error :", error)
	}

	Accounts, err = controller.NewAccountAPI().AccountTransfer(tx)
	if err != nil {
		log.Println(err)
	}

	Response(w, Accounts)

}

func AccountTransfer_CC(w http.ResponseWriter, r *http.Request) {

	var Accounts []model.Account
	var tx model.Tx

	requestBody, err := ioutil.ReadAll(r.Body) // read JSON data from Body
	if err != nil {
		panic("Enter all required information !!!")
	}

	error := json.Unmarshal(requestBody, &tx) // Convert from JSON format to Account Format
	if error != nil {
		fmt.Printf("can't unmarshal data ! %+v \n", requestBody)
	}

	Accounts, err = controller.NewAccountAPI().AccountTransfer_CC(tx)
	if err != nil {
		fmt.Println(err)
	}

	Response(w, Accounts)

	// Returns the new state of 2 Accounts when the transaction is done
}

func AddTransaction(w http.ResponseWriter, r *http.Request) {
	var tran model.Transactions

	w.Header().Set("Content-Type", "application/json")
	requestBody, _ := ioutil.ReadAll(r.Body) // read JSON data from Body

	json.Unmarshal(requestBody, &tran) // Convert from JSON format to Account Format

	tran, err := controller.NewAccountAPI().AddTransaction(tran)
	if err != nil {
		log.Println(err)
	}
	Response(w, tran)

}

func AddAction(w http.ResponseWriter, r *http.Request) {
	var action model.Diaries

	w.Header().Set("Content-Type", "application/json")
	requestBody, _ := ioutil.ReadAll(r.Body) // read JSON data from Body

	json.Unmarshal(requestBody, &action) // Convert from JSON format to Account Format

	action, err := controller.NewAccountAPI().AddAction(action)
	if err != nil {
		log.Println(err)
	}
	Response(w, action)

}

func GetActionDetail(w http.ResponseWriter, r *http.Request) {
	var diary model.Diaries

	vars := mux.Vars(r)
	id := vars["id"] // Get from URL path

	diary, err := controller.NewAccountAPI().GetActionDetail(id)
	if err != nil {
		log.Println(err)
	}
	Response(w, diary)
}

func GetDiaryByProduct(w http.ResponseWriter, r *http.Request) {
	var diary []model.Diaries

	vars := mux.Vars(r)
	productID := vars["productID"] // Get from URL path

	diary, err := controller.NewAccountAPI().GetDiaryByProduct(productID)
	if err != nil {
		log.Println(err)
	}
	Response(w, diary)
}

func Response(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(i)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
