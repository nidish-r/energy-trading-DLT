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
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

// =======================================
// Getter functions for all assets - swapssNetwork, swapping station, user and batteries
// =======================================

// ReadSwappingStation returns the network Id for TruePower stored in the world state.
func get_truePowerNetworkId(stub shim.ChaincodeStubInterface) (string, error) {
	var err error
	fmt.Println("starting read true power network id")

	truePowerIdBytes, err := stub.GetState(TruePowerNetworkPrefix) //get the var from ledger
	if err != nil {
		return "", errors.New("{\"Error\":\"Failed to get state for " + "TruePower" + "\"}")
	}

	fmt.Println("- end read")
	return string(truePowerIdBytes), nil
}

func get_ssNetwork(stub shim.ChaincodeStubInterface, id_Network string) (SSNetwork, error) {
	var swapssNetwork SSNetwork
	swapssNetworkAsBytes, err := stub.GetState(id_Network)
	if err != nil {
		return swapssNetwork, errors.New("Failed to find swapssNetwork - " + id_Network)
	}
	json.Unmarshal(swapssNetworkAsBytes, &swapssNetwork)

	if swapssNetwork.Id_Network != id_Network {
		return swapssNetwork, errors.New("swapssNetwork does not exist - " + id_Network)
	}

	return swapssNetwork, nil
}

func get_swappingStation(stub shim.ChaincodeStubInterface, id_swappingStation string) (SwappingStation, error) {
	var swappingStation SwappingStation
	swappingStationAsBytes, err := stub.GetState(id_swappingStation)
	if err != nil {
		return swappingStation, errors.New("Failed to find swappingStation - " + id_swappingStation)
	}
	json.Unmarshal(swappingStationAsBytes, &swappingStation)

	if swappingStation.Id_swappingStation != id_swappingStation {
		return swappingStation, errors.New("SwappingStation does not exist - " + id_swappingStation)
	}

	return swappingStation, nil
}

func get_user(stub shim.ChaincodeStubInterface, id_user string) (User, error) {
	var user User
	userAsBytes, err := stub.GetState(id_user)
	if err != nil {
		return user, errors.New("Failed to find user - " + id_user)
	}
	json.Unmarshal(userAsBytes, &user)

	if user.Id_user != id_user {
		return user, errors.New("User does not exist - " + id_user)
	}

	return user, nil
}

func get_battery(stub shim.ChaincodeStubInterface, id_battery string) (Battery, error) {
	var battery Battery
	batteryAsBytes, err := stub.GetState(id_battery)
	if err != nil {
		return battery, errors.New("Failed to find battery - " + id_battery)
	}
	json.Unmarshal(batteryAsBytes, &battery)

	if battery.Id_battery != id_battery {
		return battery, errors.New("Battery does not exist - " + id_battery)
	}

	return battery, nil
}

func formatResponse(status string, code string, message string, result interface{}) interface{} {
	var response Response
	response.Status = status
	response.Code = code
	response.Message = message
	response.Result = result
	responseAsBytes, _ := json.Marshal(response)
	resultAsString := "RESULT-->" + string(responseAsBytes) + "<--RESULT"

	if status == "OK" {
		//String is type casted to byte array
		return []byte(resultAsString)
	} else if status == "ERROR" {
		return resultAsString
	}
	//Code should not reach here. status can have only two values: OK and ERROR
	return nil
}

func formatSuccess(code string, message string, result interface{}) []byte {
	response := formatResponse("OK", code, message, result)
	return response.([]byte)
}

// Wrapper functions are provided as shim functions require type assertions for return values
func formatError(code string, message string, result interface{}) string {
	response := formatResponse("ERROR", code, message, result)
	return response.(string)
}

// ==============================================================
// Payment Helper functions - internal methods for use inside chaincode
// ==============================================================

func _payToSSNetwork(stub shim.ChaincodeStubInterface, id_Network string, charge float32) error {
	var err error
	fmt.Println("starting PayToSSNetwork")

	swapssNetwork, err := get_ssNetwork(stub, id_Network)
	if err != nil {
		fmt.Println("swapssNetwork not found in Blockchain - " + id_Network)
		return err
	}

	swapssNetwork.Id_Network = id_Network
	oemWalletFloat, err := addFloat(swapssNetwork.Wallet, charge)
	swapssNetwork.Wallet = float32(oemWalletFloat)

	//Store the user in ledger
	swapssNetworkAsBytes, _ := json.Marshal(swapssNetwork)
	err = stub.PutState(swapssNetwork.Id_Network, swapssNetworkAsBytes)
	if err != nil {
		fmt.Println("Could not store Battery OEM")
		return err
	}

	fmt.Println("- end payToWallet Battery OEM")
	return err
}

func _payToUser(stub shim.ChaincodeStubInterface, id_User string, charge float32) error {
	var err error
	fmt.Println("starting PayToSSNetwork")

	user, err := get_user(stub, id_User)
	if err != nil {
		fmt.Println("swapssNetwork not found in Blockchain - " + id_User)
		return err
	}

	user.Id_user = id_User
	userWalletFloat, err := addFloat(user.Wallet, charge)
	user.Wallet = userWalletFloat
	//Store the user in ledger
	swapssNetworkAsBytes, _ := json.Marshal(user)
	err = stub.PutState(user.Id_user, swapssNetworkAsBytes)
	if err != nil {
		fmt.Println("Could not store user")
		return err
	}

	return err
}

func _payFromUser(stub shim.ChaincodeStubInterface, id_User string, charge float32) error {
	var err error
	fmt.Println("starting PayToSSNetwork")

	user, err := get_user(stub, id_User)
	if err != nil {
		fmt.Println("swapssNetwork not found in Blockchain - " + id_User)
		return err
	}

	user.Id_user = id_User
	userWalletFloat, err := subFloat(user.Wallet, charge)
	if err != nil {
		fmt.Println("Could not charge user wallet")
		return err
	}
	user.Wallet = userWalletFloat
	//Store the user in ledger
	batteryOEMAsBytes, _ := json.Marshal(user)
	err = stub.PutState(user.Id_user, batteryOEMAsBytes)
	if err != nil {
		fmt.Println("Could not store user")
		return err
	}

	return err
}

// ==============================================================
// Input Sanitation - dumb input checking, look for empty strings
// ==============================================================
func sanitize_arguments(strs []string) error {
	for i, val := range strs {
		if len(val) <= 0 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be a non-empty string")
		}
		if len(val) > 256 {
			errMsg := "Argument " + strconv.Itoa(i) + " must be <= 256 characters"
			return errors.New(errMsg)
		}
	}
	return nil
}

// ==============================================================
// Arithmetic functions to check for overflow and underflow
// ==============================================================

// add two number checking for overflow
func add(b int, q int) (int, error) {

	// Check overflow
	var sum int
	sum = q + b

	if (sum < q) == (b >= 0 && q >= 0) {
		return 0, fmt.Errorf("Math: addition overflow occurred %d + %d", b, q)
	}

	return sum, nil
}

// sub two number checking for overflow
func sub(b int, q int) (int, error) {

	// Check overflow
	var diff int
	diff = b - q

	if (diff > b) == (b >= 0 && q >= 0) {
		return 0, fmt.Errorf("Math: Subtraction overflow occurred  %d - %d", b, q)
	}

	return diff, nil
}

// add two float numbers checking for overflow
func addFloat(b float32, q float32) (float32, error) {

	// Check overflow
	var sum float32
	sum = q + b

	if (sum < q) == (b >= 0 && q >= 0) {
		return 0, fmt.Errorf("Math: addition overflow occurred %d + %d", b, q)
	}

	return sum, nil
}

// sub two float numbers checking for overflow
func subFloat(b float32, q float32) (float32, error) {

	// Check overflow
	var diff float32
	diff = b - q

	if (diff > b) == (b >= 0 && q >= 0) {
		return 0, fmt.Errorf("Math: Subtraction overflow occurred  %d - %d", b, q)
	}

	return diff, nil
}
