/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the

wit"License"); you may not use this file except in compliance the License.  You may obtain a copy of the License at

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
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
	contractapi.Contract
}

// ============================================================================================================================
// Asset Definitions - The ledger will battery, swappingStation and user
// ============================================================================================================================

//
// Structure for response
//
type Response struct {
	Status  string      `json:"status"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

//
// Structure for `SSNetwork`
//
type SSNetwork struct {
	Id_Network          string  `json:"Id_Network"`
	Name_Network        string  `json:"name_network"`
	UnverifiedBatteries uint64  `json:"unverifiedBatteries"`
	TotalBatteries      uint64  `json:"totalBatteries"`
	ExpiredBatteries    uint64  `json:"expiredBatteries"`
	Status              string  `json:"status"`
	Wallet              float32 `json:"wallet"`
	DocType             string  `json:"docType"`
}

//
// Structure for `Battery`
//
type Battery struct {
	Id_battery       string  `json:"id_battery"`
	ModelNumber      string  `json:"modelNumber"`
	SoC              uint8   `json:"soC"`
	SoH              uint8   `json:"soH"`
	EnergyContent    float32 `json:"energyContent"`
	Cdc              uint16  `json:"cdc"`
	DockedStation    string  `json:"dockedStation"`
	AllocatedToFleet string  `json:"allocatedToFleet"`
	Company          string  `json:"company"`
	EscrowedAmount   float32 `json:"escrowedAmount"`
	Id_Network       string  `json:"id_Network"`
	User             string  `json:"user"`
	Owner            string  `json:"owner"`
	Status           string  `json:"status"`
	ManufacturerId   string  `json:"manufacturerId"`
	ManufactureDate  string  `json:"manufactureDate"`
	DocType          string  `json:"docType"`
}

//
// Structure for `SwappingStation`
//
type SwappingStation struct {
	Id_swappingStation  string `json:"id_swappingStation"`
	SwappingStationName string `json:"swappingStationName"`
	Id_Network          string `json:"id_Network"`
	UnverifiedBatteries uint64 `json:"unverifiedBatteries"`
	TotalBatteries      uint64 `json:"totalBatteries"`
	ActiveBatteries     uint64 `json:"activeBatteries"`
	ExpiredBatteries    uint64 `json:"expiredBatteries"`
	DischargedBatteries uint64 `json:"dischargedBatteries"`
	GeoCoordinates      string `json:"geoCoordinates"`
	Address             string `json:"address"`
	LicenseNumber       string `json:"licenseNumber"`
	EmailId             string `json:"emailId"`
	ContactNumber       string `json:"contactNumber"`
	Company             string `json:"company"`
	DocType             string `json:"docType"`
}

//
// Structure for `Fleet`
//
type Fleet struct {
	Id_fleet       string `json:"id_fleet"`
	FleetName      string `json:"fleetName"`
	TotalBatteries uint64 `json:"totalBatteries"`
	Company        string `json:"company"`
	Industry       string `json:"industry"` // can be used to differentiate fleets (ninjacart, zepto), last-mile (delhivery, amazon), ride-hailing (rapido, ola bike), and corporate campuses (avis india, etc.)
	EmailId        string `json:"emailId"`
	ContactNumber  string `json:"contactNumber"`
	Address        string `json:"address"`
	DocType        string `json:"docType"`
}

//
// Structure for `User`
//
type User struct {
	Id_user       string  `json:"id_user"`
	UserName      string  `json:"userName"`
	Address       string  `json:"address"`
	AadharNumber  string  `json:"aadharNumber"`
	EmailId       string  `json:"emailId"`
	FleetId       string  `json:"fleetId"`
	Company       string  `json:"company"`
	MobileNumber  string  `json:"mobileNumber"`
	RentedBattery string  `json:"rentedBattery"`
	Wallet        float32 `json:"wallet"`
	DocType       string  `json:"docType"`
}

// ============================================================================================================================
// Prefix Definitions - For creating composite keys and avoid id overlap (for future use)
// ============================================================================================================================

const SSNetworkPrefix = "OEM"
const SwappingStationPrefix = "SS"
const BatteryPrefix = "Battery"
const UserPrefix = "User"
const TruePowerNetworkPrefix = "TruePowerNetwork"
const FleetPrefix = "Fleet"

// ============================================================================================================================
// Enum Definitions - Absolute states of allowed status for different assets (WIP)
// ============================================================================================================================

// For use with specific enums for indexing
type BatteryStatus int64

const (
	Prebook    BatteryStatus = iota // = 0
	Cancelled                       // = 1
	InProgress                      // = 2
	Completed                       // = 3
)

var (
	sessionMap = map[string]BatteryStatus{
		"Prebook":    Prebook,
		"Cancelled":  Cancelled,
		"InProgress": InProgress,
		"Completed":  Completed,
	}
)

func SessionStatusString(status BatteryStatus) string {
	return []string{"Prebook", "Cancelled", "InProgress", "Completed"}[status]
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode - %s", err)
	}
}

// ============================================================================================================================
// Init - initialize the chaincode (needed for interface) - returning empty success response
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("starting invoke, for - " + function)

	// Handle different functions
	if function == "Write" { //generic writes to ledger
		return Write(stub, args)
	} else if function == "InitializeSSNetwork" { //write functions
		return InitializeSSNetwork(stub, args)
	} else if function == "InitializeSwappingStation" {
		return InitializeSwappingStation(stub, args)
	} else if function == "InitializeBattery" {
		return InitializeBattery(stub, args)
	} else if function == "InitializeUser" {
		return InitializeUser(stub, args)
	} else if function == "RechargeUserWallet" {
		return RechargeUserWallet(stub, args)
	} else if function == "DockOONBatteryOnSwappingStation" {
		return DockOONBatteryOnSwappingStation(stub, args)
	} else if function == "VerifiyOONBatteryOnSS" {
		return VerifiyOONBatteryOnSS(stub, args)
	} else if function == "DockBatteryOnSwappingStation" {
		return DockBatteryOnSwappingStation(stub, args)
	} else if function == "TransferBatteryFromSSToUser" {
		return TransferBatteryFromSSToUser(stub, args)
	} else if function == "TransferBatteryFromUserToSS" {
		return TransferBatteryFromUserToSS(stub, args)
	} else if function == "ReturnBatteryFromService" {
		return ReturnBatteryFromService(stub, args)
	} else if function == "MarkBatteryStolen" {
		return MarkBatteryStolen(stub, args)
	} else if function == "MarkBatteryError" {
		return MarkBatteryError(stub, args)
	} else if function == "MarkBatteryExpired" {
		return MarkBatteryExpired(stub, args)
	} else if function == "ReadSSNetwork" { // read functions
		return ReadSSNetwork(stub, args)
	} else if function == "ReadSwappingStation" {
		return ReadSwappingStation(stub, args)
	} else if function == "ReadUser" {
		return ReadUser(stub, args)
	} else if function == "ReadBattery" {
		return ReadBattery(stub, args)
	} else if function == "ReadBatteryHistory" {
		return ReadBatteryHistory(stub, args)
	} else if function == "InitializeFleet" {
		return InitializeFleet(stub, args)
	} else if function == "AllocateBatteryToFleet" {
		return AllocateBatteryToFleet(stub, args)
	} else if function == "DeallocateBatteryFromFleet" {
		return DeallocateBatteryFromFleet(stub, args)
	} else if function == "TransferBatteryBetweenFleets" {
		return TransferBatteryBetweenFleets(stub, args)
	} else if function == "GenerateFleetReport" {
		return GenerateFleetReport(stub, args)
	}

	// error out
	fmt.Println("Received unknown invoke function name - " + function)
	return shim.Error("Received unknown invoke function name - '" + function + "'")
}

// ============================================================================================================================
// Query - legacy function (needed for interface)
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Unknown supported call - Query()")
}
