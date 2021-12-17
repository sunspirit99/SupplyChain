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
	SupplyChain_Chaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating SupplyChain_Chaincode chaincode: %v", err)
	}

	if err := SupplyChain_Chaincode.Start(); err != nil {
		log.Panicf("Error starting SupplyChain_Chaincode chaincode: %v", err)
	}

}
