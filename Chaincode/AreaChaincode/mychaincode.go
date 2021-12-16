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
	areaChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating AreaChaincode chaincode: %v", err)
	}

	if err := areaChaincode.Start(); err != nil {
		log.Panicf("Error starting AreaChaincode chaincode: %v", err)
	}
}
