package Controller

import (
	cfg "SuperBank/Config"
	m "SuperBank/Model"
	"SuperBank/fabric"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

var c = cfg.GetConfig()
var (
	channelName         string = c.GetString("channelName")
	AreaContract        string = c.GetString("AreaContract")
	DiaryContract       string = c.GetString("DiaryContract")
	TransactionContract string = c.GetString("TransactionContract")
	AccountContract     string = c.GetString("AccountContract")

	maxWorker      = c.GetInt("maxWorker")
	maxTransaction = c.GetInt("maxTransaction")
	maxAccount     = c.GetInt("maxAccount")

	processing = c.GetInt("processing")
	success    = c.GetInt("success")
	fail       = c.GetInt("fail")
)

var trans chan m.Tx
var accounts chan []m.Account
var gw *gateway.Gateway

type AccountAPI struct{}

func NewAccountAPI() *AccountAPI {
	return &AccountAPI{}
}

func (api AccountAPI) GetAllAccount() ([]m.Account, error) {

	var AccountORM = m.NewAccountORM()

	Accounts, error := AccountORM.FindAllAccount() // // Return all Accounts exists in table of database
	if error != nil {
		fmt.Println("Query error !")
		return []m.Account{}, error
	}
	return Accounts, nil

	// Returns a list of Accounts in JSON format
}

// GetAccountByID get Account with specific ID
func (api AccountAPI) GetAccountByID(id string) (m.Account, error) {

	var AccountORM = m.NewAccountORM()

	account, error := AccountORM.FindAccountByID(id) // Return Account with id = key
	if error != nil {
		fmt.Println("Query error")
		return m.Account{}, error
	}

	return account, nil
	// Returns the found Account in JSON format
}

//CreateAccount creates an Account
func (api AccountAPI) CreateAccount(Account m.Account) (m.Account, error) {

	var AccountORM = m.NewAccountORM()
	gw := fabric.InitGateway()

	network, err := gw.GetNetwork(channelName)
	if err != nil {
		fmt.Println("Failed to get network")
	}

	id := Account.Id
	name := Account.Name
	address := Account.Address
	phonenumber := Account.PhoneNumber
	balance := fmt.Sprintf("%f", Account.Balance)
	status := strconv.Itoa(Account.Status)
	createtime := Account.Createtime

	contract := network.GetContract(AccountContract)

	response, err := contract.SubmitTransaction("CreateAccount", id, name, address, phonenumber, balance, status, createtime)
	if err != nil {
		log.Fatalln("Failed to create account in fabric network !", err)
	}

	fmt.Println("Response from Fabric : ", string(response.Responses[0].Response.Payload))

	Account.Createtime = time.Now().String()
	Account, error := AccountORM.CreateAccount(Account) // Create a record in database
	if error != nil {
		fmt.Println("Failed to create account in mysql !")
		return m.Account{}, error
	}

	return Account, nil

	// return the created Account in JSON format
}

// UpdateAccountByID updates Account with respective ID, if ID does not exist, create a new Account
func (api AccountAPI) UpdateAccountByID(Account m.Account) (m.Account, error) {

	var AccountORM = m.NewAccountORM()
	Account, error := AccountORM.SaveAccount(Account) // Update database
	if error != nil {
		fmt.Println("Query Error !")
		return m.Account{}, error
	}

	fmt.Println("Updated Account successfully !")
	return Account, nil
	// Return the updated Account in JSON format
}

// DeleteAccountByID delete an Account with specific ID
func (api AccountAPI) DeleteAccountByID(id string) error {

	var AccountORM = m.NewAccountORM()

	err := AccountORM.DeleteAccount(id)
	if err != nil {
		return err
	}
	fmt.Println("[ID :", id, "] has been successfully deleted !")

	return nil
}

// AccountWithdraw withdraw money from account through a transaction
func (api AccountAPI) AccountWithdraw(tx m.Tx) (m.Account, error) {

	var Account m.Account

	var txORM = m.NewTxORM()
	var AccountORM = m.NewAccountORM()

	tx.Trace = uuid.New().String()
	// fmt.Print(tx.Trace)
	now := time.Now()
	tx.Createtime = &now
	tx.Status = processing // processing

	tx, err := txORM.CreateTx(tx) // Create a record in database
	if err != nil {
		fmt.Println("Failed to create transaction in DB !")
		return m.Account{}, err
	}

	gw := fabric.InitGateway()
	network, err := gw.GetNetwork(channelName)
	if err != nil {
		fmt.Println("Failed to get network !")
	}
	contract := network.GetContract(AccountContract)

	from := tx.To
	amount := fmt.Sprintf("%f", tx.Amount)

	response, err := contract.SubmitTransaction("AccountWithdraw", from, amount)
	if err != nil {
		fmt.Println("Failed to send proposal !")
		txORM.UpdateTxID(fail, "", tx.Trace)

		return Account, err
	} else {
		payload := response.Responses[0].ProposalResponse.Response.Payload
		txid := string(response.TransactionID)

		err = json.Unmarshal(payload, &Account)
		if err != nil {
			log.Fatalln("err :", err)
		}

		fmt.Println("You have successfully withdrew", tx.Amount, "to your account !")

		AccountORM.UpdateBalance(Account.Balance, Account.Id)
		txORM.UpdateTxID(success, txid, tx.Trace)

	}

	return Account, nil
	// Returns the new state of the Account when the transaction is done
}

// AccountDeposit deposit money into account through a transaction
func (api AccountAPI) AccountDeposit(tx m.Tx) (m.Account, error) {

	var Account m.Account
	var txORM = m.NewTxORM()
	var AccountORM = m.NewAccountORM()

	tx.Trace = uuid.New().String()
	// fmt.Print(tx.Trace)
	now := time.Now()
	tx.Createtime = &now
	tx.Status = processing // processing

	tx, err := txORM.CreateTx(tx) // Create a record in database
	if err != nil {
		fmt.Println("Failed to create transaction in DB !")
		return m.Account{}, err
	}

	gw := fabric.InitGateway()
	network, err := gw.GetNetwork(channelName)
	if err != nil {
		fmt.Println("Failed to get network !")
	}
	contract := network.GetContract(AccountContract)

	from := tx.To
	amount := fmt.Sprintf("%f", tx.Amount)
	response, err := contract.SubmitTransaction("AccountDeposit", from, amount)
	if err != nil {
		fmt.Println("Failed to send proposal !")
		txORM.UpdateTxID(fail, "", tx.Trace) // status = 2 : Failed

		return Account, err
	} else {
		payload := response.Responses[0].ProposalResponse.Response.Payload
		txid := string(response.TransactionID)

		err = json.Unmarshal(payload, &Account)
		if err != nil {
			log.Fatalln("err :", err)
		}

		fmt.Println("You have successfully deposited", tx.Amount, "to your account !")

		AccountORM.UpdateBalance(Account.Balance, Account.Id)
		txORM.UpdateTxID(success, txid, tx.Trace)
	}

	return Account, nil
	// Returns the new state of the Account when the Tx is done
}

// AccountTransfer transfer money between 2 accounts through a transaction
func (api AccountAPI) AccountTransfer(tx m.Tx) ([]m.Account, error) {

	var Accounts []m.Account
	var txORM = m.NewTxORM()
	var AccountORM = m.NewAccountORM()

	start := time.Now()
	gw = fabric.InitGateway()
	fmt.Println("OK")

	Accounts, err := AccountORM.FindAccountByTX(tx)
	if err != nil {
		return []m.Account{}, err
	}
	fmt.Println("Before :", Accounts)

	tx.Trace = uuid.New().String()
	// fmt.Print(tx.Trace)
	now := time.Now()
	tx.Createtime = &now
	tx.Status = processing // processing

	tx, err = txORM.CreateTx(tx) // Create a record in database
	if err != nil {
		fmt.Println("Failed to create Tx in DB !")
		return []m.Account{}, err
	}

	network, err := gw.GetNetwork(channelName)
	if err != nil {
		fmt.Println("Failed to get network from gateway !")
	}

	contract := network.GetContract(AccountContract)

	from := tx.From
	to := tx.To
	amount := strconv.Itoa(int(tx.Amount))

	response, err := contract.SubmitTransaction("AccountTransfer", from, to, amount)
	if err != nil {
		fmt.Println("Failed to submit transaction ! Can't update the balance !")
		txORM.UpdateTxID(fail, "", tx.Trace)
		return Accounts, err
	} else {
		txid := string(response.TransactionID)

		fmt.Println("txid :", txid)
		// fmt.Println(string(response.Responses[0].ProposalResponse.Payload))
		payload := response.Responses[0].ProposalResponse.Response.Payload

		// fmt.Println(string(payload))

		err := json.Unmarshal(payload, &Accounts)
		if err != nil {
			log.Fatalln("err :", err)
		}
		fmt.Println("After :", Accounts)

		txORM.UpdateTxID(success, txid, tx.Trace)
		AccountORM.UpdateBalance(Accounts[0].Balance, tx.From)
		AccountORM.UpdateBalance(Accounts[1].Balance, tx.To)

	}
	fmt.Println("Execution time :", time.Since(start))

	return Accounts, nil

}

func (api AccountAPI) AccountTransfer_CC(tx m.Tx) ([]m.Account, error) {

	var Accounts []m.Account
	go func() {
		fmt.Println("goroutine is processing !")
		trans <- tx
	}()

	Accounts = <-accounts

	fmt.Println("After :", Accounts)

	return Accounts, nil

	// Returns the new state of 2 Accounts when the transaction is done
}

func (api AccountAPI) AddTransaction(tran m.Transactions) (m.Transactions, error) {
	var TransORM = m.NewTransactionsORM()
	// Account.Createtime = time.Now().String()
	gw := fabric.InitGateway()

	network, err := gw.GetNetwork(channelName)
	if err != nil {
		fmt.Println("Failed to get network")
	}

	contract := network.GetContract(TransactionContract)

	created_time := time.Now().String()

	response, err := contract.SubmitTransaction("AddTransaction", tran.Id, tran.Transport_id, tran.Related_batch, tran.Pickup_place, tran.Deliver_place, created_time, tran.Description)
	if err != nil {
		log.Fatalln("Failed to create account in fabric network !", err)
	}

	fmt.Println("Response from Fabric : ", string(response.Responses[0].Response.Payload))

	tran, err = TransORM.CreateTx(tran)
	if err != nil {
		fmt.Println(err)
		return m.Transactions{}, err
	}
	return tran, nil
}

func (api AccountAPI) AddAction(dia m.Diaries) (m.Diaries, error) {
	var DiariesORM = m.NewDiariesORM()

	gw := fabric.InitGateway()

	network, err := gw.GetNetwork(channelName)
	if err != nil {
		fmt.Println("Failed to get network")
	}

	contract := network.GetContract(DiaryContract)

	created_time := time.Now().String()

	response, err := contract.SubmitTransaction("AddAction", dia.Id, dia.Title, dia.Content, dia.Created_by, created_time, dia.Related_area, dia.Related_batch, dia.Related_products)
	if err != nil {
		log.Fatalln("Failed to create account in fabric network !", err)
	}

	fmt.Println("Response from Fabric : ", string(response.Responses[0].Response.Payload))

	action, err := DiariesORM.CreateDiary(dia)
	if err != nil {
		fmt.Println(err)
		return m.Diaries{}, err
	}
	return action, nil
}

func (api AccountAPI) GetActionDetail(id string) (m.Diaries, error) {
	var DiariesORM = m.NewDiariesORM()
	diary, err := DiariesORM.GetDiaryById(id)
	if err != nil {
		fmt.Println(err)
		return m.Diaries{}, err
	}
	return diary, nil
}

func (api AccountAPI) GetDiaryByProduct(product string) ([]m.Diaries, error) {
	var DiariesORM = m.NewDiariesORM()
	diary, err := DiariesORM.GetDiaryByProduct(product)
	if err != nil {
		fmt.Println(err)
		return []m.Diaries{}, err
	}
	return diary, nil
}

func InitWorker() {
	trans = make(chan m.Tx, maxTransaction)
	accounts = make(chan []m.Account, maxAccount)

	gw := fabric.InitGateway()
	network, err := gw.GetNetwork(channelName)
	if err != nil {
		log.Fatal("Failed to get network !")
	}

	contract := network.GetContract(AccountContract)
	for i := 0; i < maxWorker; i++ {
		fmt.Printf("Worker %v is waiting ..... ! \n", i+1)
		go func() {
			Worker(trans, accounts, contract)
		}()
	}
}

func Worker(trans <-chan m.Tx, accounts chan<- []m.Account, contract *gateway.Contract) {

	var txORM m.TxORM
	var AccountORM = m.NewAccountORM()

	fmt.Println("Worker is working !")
	for tx := range trans {
		fmt.Println("Still processing _________")
		if tx.From == tx.To {
			log.Panicf("Please enter correct recipient ID ! %v and %v", tx.From, tx.To)
		}

		Accounts, err := AccountORM.FindAccountByTX(tx)
		if err != nil {
			log.Panicf("Cant find account by tx")
		}
		fmt.Println("Before :", Accounts)

		tx.Trace = uuid.New().String()
		// fmt.Print(tx.Trace)
		now := time.Now()
		tx.Createtime = &now
		tx.Status = processing // processing

		tx, err := txORM.CreateTx(tx) // Create a record in database
		if err != nil {
			log.Fatalln("Failed to create transaction in DB !")
		}

		from := tx.From
		to := tx.To
		amount := fmt.Sprintf("%f", tx.Amount)

		response, err := contract.SubmitTransaction("AccountTransfer", from, to, amount)
		if err != nil {
			fmt.Println("Failed to submit transaction ! Can't update the balance !")
			txORM.UpdateTxID(fail, "", tx.Trace)
		} else {
			txid := string(response.TransactionID)

			fmt.Println("txid :", txid)

			payload := response.Responses[0].ProposalResponse.Response.Payload

			// fmt.Println(string(payload))

			err := json.Unmarshal(payload, &Accounts)
			if err != nil {
				log.Fatalln("err :", err)
			}
			// fmt.Println(Accounts)

			txORM.UpdateTxID(success, txid, tx.Trace)
			AccountORM.UpdateBalance(Accounts[0].Balance, tx.From)
			AccountORM.UpdateBalance(Accounts[1].Balance, tx.To)

		}
		// time.Sleep(time.Second)
		accounts <- Accounts
	}
}
