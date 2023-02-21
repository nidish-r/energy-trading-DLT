/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

/* -------------------------------------------------------------------------- */
/*                     Swapping Network Related Methods                       */
/* -------------------------------------------------------------------------- */

// ReadSSNetwork returns the ssNetwork stored in the world state with given id.
func ReadSSNetwork(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var key, jsonResp string
	var err error
	fmt.Println("starting read")

	if len(args) != 1 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting key of the var to query"))
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	key = args[0]
	ssNetworkAsbytes, err := stub.GetState(key) //get the var from ledger
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(formatError("NA", "NA", jsonResp))
	}

	var ssNetwork SSNetwork
	json.Unmarshal(ssNetworkAsbytes, &ssNetwork)
	if ssNetwork.DocType != "ssNetwork" {
		jsonResp = "{\"Error\":\"No ssNetwork was found with ssNetwork id " + key + "\"}"
		return shim.Error(formatError("NA", "NA", jsonResp))
	}

	fmt.Println("- end read")
	return shim.Success(formatSuccess("NA", "NA", ssNetwork))
}

/* -------------------------------------------------------------------------- */
/*                       Swapping Station Related Methods                     */
/* -------------------------------------------------------------------------- */

// ReadSwappingStation returns the swappingStation stored in the world state with given id.
func ReadSwappingStation(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var key, jsonResp string
	var err error
	fmt.Println("starting read")

	if len(args) != 1 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting key of the var to query"))
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	key = args[0]
	swappingStationAsbytes, err := stub.GetState(key) //get the var from ledger
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(formatError("NA", "NA", jsonResp))
	}

	var swappingStation SwappingStation
	json.Unmarshal(swappingStationAsbytes, &swappingStation)
	if swappingStation.DocType != "swappingStation" {
		jsonResp = "{\"Error\":\"No swappingStation was found with swappingStation id " + key + "\"}"
		return shim.Error(formatError("NA", "NA", jsonResp))
	}

	fmt.Println("- end read")
	return shim.Success(formatSuccess("NA", "NA", swappingStation))
}

/* -------------------------------------------------------------------------- */
/*                            User Related Methods                            */
/* -------------------------------------------------------------------------- */

// ReadUser returns the user stored in the world state with given id.
func ReadUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, jsonResp string
	var err error
	fmt.Println("starting read")

	if len(args) != 1 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting key of the var to query"))
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	key = args[0]
	userAsbytes, err := stub.GetState(key) //get the var from ledger
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(formatError("NA", "NA", jsonResp))
	}

	var user User
	json.Unmarshal(userAsbytes, &user)
	if user.DocType != "user" {
		jsonResp = "{\"Error\":\"No user was found with user id " + key + "\"}"
		return shim.Error(formatError("NA", "NA", jsonResp))

	}

	fmt.Println("- end read")
	return shim.Success(formatSuccess("NA", "NA", user))
}

/* -------------------------------------------------------------------------- */
/*                          Battery Related Methods                           */
/* -------------------------------------------------------------------------- */

// ReadBattery returns the battery stored in the world state with given id.
func ReadBattery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, jsonResp string
	var err error
	fmt.Println("starting read")

	if len(args) != 1 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting key of the var to query"))
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	key = args[0]
	batteryAsbytes, err := stub.GetState(key) //get the var from ledger
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(formatError("NA", "NA", jsonResp))
	}

	var battery Battery
	json.Unmarshal(batteryAsbytes, &battery)
	if battery.DocType != "battery" {
		jsonResp = "{\"Error\":\"No battery was found with battery id " + key + "\"}"
		return shim.Error(formatError("NA", "NA", jsonResp))

	}

	fmt.Println("- end read")
	return shim.Success(formatSuccess("NA", "NA", battery))
}

// ReadBatteryHistory returns the battery audit history of a given id from the world state.
func ReadBatteryHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AuditHistory struct {
		TxId      string    `json:"txId"`
		Timestamp time.Time `json:"timestamp"`
		Value     Battery   `json:"value"`
	}

	var history []AuditHistory
	var battery Battery

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	id_battery := args[0]
	fmt.Printf("- start getHistoryForBattery: %s\n", id_battery)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(id_battery)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		historicValue, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		txID := ""
		timestamp, err := ptypes.Timestamp(historicValue.Timestamp)
		if err != nil {
			return shim.Error(err.Error())
		}
		var tx AuditHistory
		tx.TxId = txID                                //copy transaction id over
		json.Unmarshal(historicValue.Value, &battery) //un stringify it aka JSON.parse()
		if historicValue == nil {                     //battery has been deleted
			var emptyMarble Battery
			tx.Value = emptyMarble //copy nil battery
		} else {
			json.Unmarshal(historicValue.Value, &battery) //un stringify it aka JSON.parse()
			tx.Value = battery                            //copy battery over
			tx.Timestamp = timestamp
			tx.TxId = historicValue.TxId
		}
		history = append(history, tx) //add this tx to the list
	}
	fmt.Printf("- getHistoryForBattery returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history) //convert to array of bytes
	return shim.Success(historyAsBytes)
}
