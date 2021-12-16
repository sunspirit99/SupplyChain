/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"MyChaincode/chaincode"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	diaryChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating diaryChaincode chaincode: %v", err)
	}

	if err := diaryChaincode.Start(); err != nil {
		log.Panicf("Error starting diaryChaincode chaincode: %v", err)
	}
}
