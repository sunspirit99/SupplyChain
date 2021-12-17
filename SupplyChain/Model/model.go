package Model

import "time"

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
type Tx struct {
	Trace      string     `json:trace gorm:"PRIMARY_KEY"`
	TxID       string     `json:txid`
	From       string     `json:"from"`
	To         string     `json:"to"`
	Amount     float32    `json:"amount"`
	Status     int        `json:"status"`
	Createtime *time.Time `json:"createtime"`
}

type Areas struct {
	Id           string  `json:id gorm:"PRIMARY_KEY"`
	Owner_id     string  `json:owner_id`
	Acreage      float32 `json:acreage`
	Acreage_unit string  `json:acreage_unit`
	Seed_id      string  `json:seed_id`
	Info         string  `json:info`
	Address      string  `json:address`
}

type Diaries struct {
	Id               string     `json:id gorm:"PRIMARY_KEY"`
	Title            string     `json:title`
	Content          string     `json:content`
	Created_by       string     `json:created_by`
	Created_time     *time.Time `json:created_time`
	Related_area     string     `json:related_area`
	Related_batch    string     `json:related_batch`
	Related_products string     `json:related_products`
}

type Diary_Fabric struct {
	Id               string `json:id gorm:"PRIMARY_KEY"`
	Title            string `json:title`
	Content          string `json:content`
	Created_by       string `json:created_by`
	Created_time     string `json:created_time`
	Related_area     string `json:related_area`
	Related_batch    string `json:related_batch`
	Related_products string `json:related_products`
}

type Transactions struct {
	Id            string     `json:id gorm:"PRIMARY_KEY"`
	Transport_id  string     `json:transport_id`
	Related_batch string     `json:related_batch`
	Pickup_place  string     `json:pickup_place`
	Deliver_place string     `json:deliver_place`
	Created_time  *time.Time `json:created_time`
	Description   string     `json:description`
}

type DBConfig struct {
	ServerName string
	User       string
	Password   string
	DB         string
}
