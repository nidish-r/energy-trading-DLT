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
	"strconv"

	//	"strings"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// ============================================================================================================================
// write() - genric write variable into ledger
//
// Shows Off PutState() - writting a key/value into the ledger
//
// Inputs - Array of strings
//    0   ,    1
//   key  ,  value
//  "abc" , "test"
// ============================================================================================================================
func Write(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, value string
	var err error
	fmt.Println("starting write")

	if len(args) != 2 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 2. key of the variable and value to set"))
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the ledger
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end write")
	return shim.Success(formatSuccess("NA", "NA", nil))
}

/* -------------------------------------------------------------------------- */
/*                        Asset Initialization Methods                        */
/* -------------------------------------------------------------------------- */

// InitializeSSNetwork adds a new SSNetwork asset in the world state with given id.
func InitializeSSNetwork(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting InitializeSSNetwork")

	if len(args) != 3 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 3"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	var ssNetwork SSNetwork
	ssNetwork.Id_Network = args[0]
	ssNetwork.Name_Network = args[1]
	ssNetwork.Status = args[2]
	ssNetwork.DocType = "ssNetwork"
	fmt.Println(ssNetwork)

	//check if ssNetwork initialization is for TruePower
	if ssNetwork.Name_Network == "TruePower" {
		//Store the TruePower network id in ledger
		err = stub.PutState(TruePowerNetworkPrefix, []byte(ssNetwork.Id_Network))
		if err != nil {
			fmt.Println("Could not store TruePowerNetwork Id")
			return shim.Error(formatError("NA", "NA", err.Error()))
		}
	}

	//check if ssNetwork already exists
	_, err = get_ssNetwork(stub, ssNetwork.Id_Network)
	if err == nil {
		fmt.Println("This SSNetwork already exists - " + ssNetwork.Id_Network)
		return shim.Error(formatError("NA", "NA", "This swappingStation already exists - "+ssNetwork.Id_Network))
	}

	//Store the SSNetwork in ledger
	ssNetworkAsBytes, _ := json.Marshal(ssNetwork)
	err = stub.PutState(ssNetwork.Id_Network, ssNetworkAsBytes)
	if err != nil {
		fmt.Println("Could not store ssNetwork")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end init ssNetwork ")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

// InitializeSwappingStation adds a new SwappingStation asset in the world state with given id.
func InitializeSwappingStation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting InitializeSwappingStation")

	if len(args) != 9 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 9"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	var swappingStation SwappingStation
	swappingStation.Id_swappingStation = args[0]
	swappingStation.SwappingStationName = args[1]
	swappingStation.Id_Network = args[2]
	swappingStation.GeoCoordinates = args[3]
	swappingStation.Address = args[4]
	swappingStation.LicenseNumber = args[5]
	swappingStation.EmailId = args[6]
	swappingStation.ContactNumber = args[7]
	swappingStation.Company = args[8]
	swappingStation.DocType = "swappingStation"
	fmt.Println(swappingStation)

	//check if swappingStation already exists
	_, err = get_swappingStation(stub, swappingStation.Id_swappingStation)
	if err == nil {
		fmt.Println("This swappingStation already exists - " + swappingStation.Id_swappingStation)
		return shim.Error(formatError("NA", "NA", "This swappingStation already exists - "+swappingStation.Id_swappingStation))
	}

	//Store the swappingStation in ledger
	swappingStationAsBytes, _ := json.Marshal(swappingStation)
	err = stub.PutState(swappingStation.Id_swappingStation, swappingStationAsBytes)
	if err != nil {
		fmt.Println("Could not store swappingStation")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end init swappingStation ")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

// InitializeBattery adds a new Battery asset in the world state with given id.
func InitializeBattery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting InitializeBattery")

	if len(args) != 9 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 9"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	var battery Battery

	battery.Id_battery = args[0]
	battery.ModelNumber = args[1]
	SoC_uint64, _ := strconv.ParseUint(args[2], 10, 0)
	battery.SoC = uint8(SoC_uint64)
	SoH_uint64, _ := strconv.ParseUint(args[3], 10, 0)
	battery.SoH = uint8(SoH_uint64)
	EnergyContent_float64, _ := strconv.ParseFloat(args[4], 32)
	battery.EnergyContent = float32(EnergyContent_float64)
	Cdc_uint64, _ := strconv.ParseUint(args[5], 10, 0)
	battery.Cdc = uint16(Cdc_uint64)
	battery.Id_Network = args[6]
	battery.ManufacturerId = args[7]
	battery.ManufactureDate = args[8]
	battery.Status = "Undocked"
	battery.DocType = "battery"
	fmt.Println(battery)

	//Check if battery already exists
	_, err = get_battery(stub, battery.Id_battery)
	if err == nil {
		fmt.Println("This battery already exists - " + battery.Id_battery)
		return shim.Error(formatError("NA", "NA", "This battery already exists - "+battery.Id_battery))
	}

	//SS_Network (owner) needs to be initialized
	ssNetwork, err := get_ssNetwork(stub, battery.Id_Network)
	if err != nil {
		fmt.Println("The Battery Network doesnt exist - " + battery.Id_Network)
		return shim.Error(formatError("NA", "NA", "The Battery Network doesnt exist - "+battery.Id_Network))
	}

	// Increment Number of Batteries for Network by 1 as new battery is added
	networktotalBatteriesInt, err := add(int(ssNetwork.TotalBatteries), 1)
	ssNetwork.TotalBatteries = uint64(networktotalBatteriesInt)

	//Store the battery in ledger
	batteryAsBytes, _ := json.Marshal(battery)
	err = stub.PutState(battery.Id_battery, batteryAsBytes)
	if err != nil {
		fmt.Println("Could not store battery")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	//Store the ssNetwork in ledger
	ssNetworkAsBytes, _ := json.Marshal(ssNetwork)
	err = stub.PutState(ssNetwork.Id_Network, ssNetworkAsBytes)
	if err != nil {
		fmt.Println("Could not store battery Network")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end init battery ")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

// InitializeFleet adds a new Fleet asset in the world state with given id.
func InitializeFleet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting InitializeFleet")

	if len(args) != 7 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 4"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	var fleet Fleet
	fleet.Id_fleet = args[0]
	fleet.FleetName = args[1]
	fleet.Company = args[2]
	fleet.Industry = args[3]
	fleet.EmailId = args[4]
	fleet.ContactNumber = args[5]
	fleet.Address = args[6]
	fleet.DocType = "fleet"
	fmt.Println(fleet)

	//check if fleet already exists
	fleetAsBytes, err := stub.GetState(fleet.Id_fleet)
	if err == nil && fleetAsBytes != nil {
		fmt.Println("This fleet already exists - " + fleet.Id_fleet)
		return shim.Error(formatError("NA", "NA", "This fleet already exists - "+fleet.Id_fleet))
	}

	//Store the fleet in ledger
	fleetAsBytes, _ = json.Marshal(fleet)
	err = stub.PutState(fleet.Id_fleet, fleetAsBytes)
	if err != nil {
		fmt.Println("Could not store fleet")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end init fleet")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

// InitializeUser adds a new User asset in the world state with given id.
func InitializeUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting InitializeUser")

	if len(args) != 8 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 7"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	var user User
	user.Id_user = args[0]
	user.UserName = args[1]
	user.Address = args[2]
	user.AadharNumber = args[3]
	user.EmailId = args[4]
	user.FleetId = args[5]
	user.Company = args[6]
	user.MobileNumber = args[7]
	user.DocType = "user"
	fmt.Println(user)

	//check if user already exists
	_, err = get_user(stub, user.Id_user)
	if err == nil {
		fmt.Println("This user already exists - " + user.Id_user)
		return shim.Error(formatError("NA", "NA", "This user already exists - "+user.Id_user))
	}

	//Store the user in ledger
	userAsBytes, _ := json.Marshal(user)
	err = stub.PutState(user.Id_user, userAsBytes)
	if err != nil {
		fmt.Println("Could not store user")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end init user ")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

/* -------------------------------------------------------------------------- */
/*                          Asset Transfer Methods                            */
/* -------------------------------------------------------------------------- */

// AllocateBatteryToFleet allocates a battery to a fleet.
func AllocateBatteryToFleet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting AllocateBatteryToFleet")

	if len(args) != 4 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 4"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	id_fleet := args[0]
	id_battery := args[1]

	// Check if the fleet exists
	fleet, err := get_fleet(stub, id_fleet)
	if err != nil {
		fmt.Println("Fleet not found in Blockchain - " + id_fleet)
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	// Check if the battery exists
	battery, err := get_battery(stub, id_battery)
	if err != nil {
		fmt.Println("Battery not found in Blockchain - " + id_battery)
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	if battery.AllocatedToFleet != "" {
		return shim.Error(formatError("NA", "NA", "Battery is already allocated to fleet - "+battery.AllocatedToFleet))
	}

	// Update the battery status and owner
	battery.Status = args[2]
	battery.Company = args[3]
	battery.AllocatedToFleet = id_fleet

	// Decrement Number of Batteries for fleet by 1 as the battery is deallocated
	fleetTotalBatteriesInt, err := add(int(fleet.TotalBatteries), 1)
	fleet.TotalBatteries = uint64(fleetTotalBatteriesInt)

	// Store the updated battery in the ledger
	batteryAsBytes, _ := json.Marshal(battery)
	err = stub.PutState(battery.Id_battery, batteryAsBytes)
	if err != nil {
		fmt.Println("Could not store battery")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	// Store the updated battery in the ledger
	fleetAsBytes, _ := json.Marshal(fleet)
	err = stub.PutState(fleet.Id_fleet, fleetAsBytes)
	if err != nil {
		fmt.Println("Could not store fleet")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end AllocateBatteryToFleet")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

// DockBatteryOnSwappingStation will link existing swap station with new discharged/charged battery, before returning battery from service (charge).
func DockBatteryOnSwappingStation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting DockBatteryOnSwappingStation")

	if len(args) != 6 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 6"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	id_battery := args[0]

	battery, err := get_battery(stub, id_battery)
	if err != nil {
		fmt.Println("battery not found in Blockchain - " + id_battery)
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	if battery.Status != "Undocked" {
		return shim.Error(formatError("NA", "NA", "Battery needs to be undocked before first docking - "+id_battery))
	}

	SoC_uint64, _ := strconv.ParseUint(args[1], 10, 0)
	battery.SoC = uint8(SoC_uint64)
	SoH_uint64, _ := strconv.ParseUint(args[2], 10, 0)
	battery.SoH = uint8(SoH_uint64)
	EnergyContent_float64, _ := strconv.ParseFloat(args[3], 32)
	battery.EnergyContent = float32(EnergyContent_float64)
	Cdc_uint64, _ := strconv.ParseUint(args[4], 10, 0)
	battery.Cdc = uint16(Cdc_uint64)
	battery.DockedStation = args[5]
	battery.Owner = "TruePower"
	battery.Status = "In_Service"

	//Swapping station (docked station) needs to be initialized
	swappingStation, err := get_swappingStation(stub, battery.DockedStation)
	if err != nil {
		fmt.Println("The Swapping Station doesnt exist - " + battery.Id_Network)
		return shim.Error(formatError("NA", "NA", "The Swapping Station doesnt exist - "+battery.Id_Network))
	}

	if battery.Company != swappingStation.Company {
		return shim.Error(formatError("NA", "NA", "Battery can only be docked on swapping station of the same company - "+id_battery))
	}

	// Increment Number of Batteries for swapping station by 1 as new battery is added
	swappingStationTotalBatteriesInt, err := add(int(swappingStation.TotalBatteries), 1)
	swappingStation.TotalBatteries = uint64(swappingStationTotalBatteriesInt)

	// Increment Number of discharged in swapping station by 1 as new uncharged battery is added
	swappingStationDischargedBatteriesInt, err := add(int(swappingStation.DischargedBatteries), 1)
	swappingStation.DischargedBatteries = uint64(swappingStationDischargedBatteriesInt)

	//Store the battery in ledger
	batteryAsBytes, _ := json.Marshal(battery)
	err = stub.PutState(battery.Id_battery, batteryAsBytes)
	if err != nil {
		fmt.Println("Could not store battery")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	//Store the swapping station in ledger
	swappingStationAsBytes, _ := json.Marshal(swappingStation)
	err = stub.PutState(swappingStation.Id_swappingStation, swappingStationAsBytes)
	if err != nil {
		fmt.Println("Could not store swapping station")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end DockBatteryOnSwappingStation")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

// DockOONBatteryOnSwappingStationgitk will dock unverified out of network battery on swap station, before verifying and returning battery from service (charge).
func DockOONBatteryOnSwappingStation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting DockOONBatteryOnSwappingStation")

	if len(args) != 7 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 7"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	id_battery := args[0]

	_, err = get_battery(stub, id_battery)
	if err == nil {
		fmt.Println("battery already found in Blockchain - " + id_battery)
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	var battery Battery

	battery.Id_battery = id_battery
	battery.Id_Network = args[1]
	battery.DockedStation = args[2]
	battery.Owner = args[3]
	battery.ModelNumber = args[4]
	battery.ManufacturerId = args[5]
	battery.ManufactureDate = args[6]
	battery.Status = "Unverified"
	battery.DocType = "battery"

	//SS_Network needs to be initialized
	ssNetwork, err := get_ssNetwork(stub, battery.Id_Network)
	if err != nil {
		fmt.Println("The Battery Network doesnt exist - " + battery.Id_Network)
		return shim.Error(formatError("NA", "NA", "The Battery Network doesnt exist - "+battery.Id_Network))
	}

	// Increment Number of Batteries for Network by 1 as new battery is added
	networkUnverifiedBatteriesInt, err := add(int(ssNetwork.UnverifiedBatteries), 1)
	ssNetwork.UnverifiedBatteries = uint64(networkUnverifiedBatteriesInt)

	//Swapping station (docked station) needs to be initialized
	swappingStation, err := get_swappingStation(stub, battery.DockedStation)
	if err != nil {
		fmt.Println("The Swapping Station doesnt exist - " + battery.DockedStation)
		return shim.Error(formatError("NA", "NA", "The Swapping Station doesnt exist - "+battery.DockedStation))
	}

	// Increment Number of Batteries for swapping station by 1 as new battery is added
	swappingStationUnverifiedBatteriesInt, err := add(int(swappingStation.UnverifiedBatteries), 1)
	swappingStation.UnverifiedBatteries = uint64(swappingStationUnverifiedBatteriesInt)

	//Store the ssNetwork in ledger
	ssNetworkAsBytes, _ := json.Marshal(ssNetwork)
	err = stub.PutState(ssNetwork.Id_Network, ssNetworkAsBytes)
	if err != nil {
		fmt.Println("Could not store battery Network")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	//Store the battery in ledger
	batteryAsBytes, _ := json.Marshal(battery)
	err = stub.PutState(id_battery, batteryAsBytes)
	if err != nil {
		fmt.Println("Could not store battery")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	//Store the swapping station in ledger
	swappingStationAsBytes, _ := json.Marshal(swappingStation)
	err = stub.PutState(swappingStation.Id_swappingStation, swappingStationAsBytes)
	if err != nil {
		fmt.Println("Could not store swapping station")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end DockOONBatteryOnSwappingStation battery ")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

// VerifiyOONBatteryOnSS will verifiy out of network battery on swap station, before putting the battery for service (charge).
func VerifiyOONBatteryOnSS(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting VerifiyOONBatteryOnSS")

	if len(args) != 5 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 5"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	id_battery := args[0]

	battery, err := get_battery(stub, id_battery)
	if err != nil {
		fmt.Println("battery not found in Blockchain. Please dock OON battery first  - " + id_battery)
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	SoC_uint64, _ := strconv.ParseUint(args[1], 10, 0)
	battery.SoC = uint8(SoC_uint64)
	SoH_uint64, _ := strconv.ParseUint(args[2], 10, 0)
	battery.SoH = uint8(SoH_uint64)
	EnergyContent_float64, _ := strconv.ParseFloat(args[3], 32)
	battery.EnergyContent = float32(EnergyContent_float64)
	Cdc_uint64, _ := strconv.ParseUint(args[4], 10, 0)
	battery.Cdc = uint16(Cdc_uint64)
	battery.Status = "In_Service"

	//SS_Network needs to be initialized
	ssNetwork, err := get_ssNetwork(stub, battery.Id_Network)
	if err != nil {
		fmt.Println("The Battery Network doesnt exist - " + battery.Id_Network)
		return shim.Error(formatError("NA", "NA", "The Battery Network doesnt exist - "+battery.Id_Network))
	}

	// Increment Number of Batteries for Network by 1 as new battery is added
	networkUnverifiedBatteriesInt, err := sub(int(ssNetwork.UnverifiedBatteries), 1)
	ssNetwork.UnverifiedBatteries = uint64(networkUnverifiedBatteriesInt)

	// Increment Number of Batteries for swapping station by 1 as battery is verified
	networkTotalBatteriesInt, err := add(int(ssNetwork.TotalBatteries), 1)
	ssNetwork.TotalBatteries = uint64(networkTotalBatteriesInt)

	//Swapping station (docked station) needs to be initialized
	swappingStation, err := get_swappingStation(stub, battery.DockedStation)
	if err != nil {
		fmt.Println("The Swapping Station doesnt exist - " + battery.Id_Network)
		return shim.Error(formatError("NA", "NA", "The Swapping Station doesnt exist - "+battery.Id_Network))
	}

	// Decrement Number of Batteries for swapping station by 1 as battery is verified
	swappingStationUnverifiedBatteriesInt, err := sub(int(swappingStation.UnverifiedBatteries), 1)
	swappingStation.UnverifiedBatteries = uint64(swappingStationUnverifiedBatteriesInt)

	// Increment Number of Batteries for swapping station by 1 as battery is verified
	swappingStationTotalBatteriesInt, err := add(int(swappingStation.TotalBatteries), 1)
	swappingStation.TotalBatteries = uint64(swappingStationTotalBatteriesInt)

	// Increment Number of discharged in swapping station by 1 as new uncharged battery is added
	swappingStationDischargedBatteriesInt, err := add(int(swappingStation.DischargedBatteries), 1)
	swappingStation.DischargedBatteries = uint64(swappingStationDischargedBatteriesInt)

	//Store the ssNetwork in ledger
	ssNetworkAsBytes, _ := json.Marshal(ssNetwork)
	err = stub.PutState(ssNetwork.Id_Network, ssNetworkAsBytes)
	if err != nil {
		fmt.Println("Could not store battery Network")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	//Store the battery in ledger
	batteryAsBytes, _ := json.Marshal(battery)
	err = stub.PutState(battery.Id_battery, batteryAsBytes)
	if err != nil {
		fmt.Println("Could not store battery")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	//Store the swapping station in ledger
	swappingStationAsBytes, _ := json.Marshal(swappingStation)
	err = stub.PutState(swappingStation.Id_swappingStation, swappingStationAsBytes)
	if err != nil {
		fmt.Println("Could not store swapping station")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end VerifiyOONBatteryOnSS battery ")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

// DeallocateBatteryFromFleet removes the association between a battery and a fleet.
func DeallocateBatteryFromFleet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting DeallocateBatteryFromFleet")

	if len(args) != 2 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 2"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	id_battery := args[0]
	id_fleet := args[1]

	battery, err := get_battery(stub, id_battery)
	if err != nil {
		fmt.Println("battery not found in Blockchain - " + id_battery)
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fleet, err := get_fleet(stub, id_fleet)
	if err != nil {
		fmt.Println("Fleet not found in Blockchain - " + id_fleet)
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	if battery.AllocatedToFleet != id_fleet {
		return shim.Error(formatError("NA", "NA", "Battery is not allocated to the specified fleet - "+id_fleet))
	}

	battery.User = ""
	battery.AllocatedToFleet = ""
	battery.Status = "Available"

	// Decrement Number of Batteries for fleet by 1 as the battery is deallocated
	fleetTotalBatteriesInt, err := sub(int(fleet.TotalBatteries), 1)
	fleet.TotalBatteries = uint64(fleetTotalBatteriesInt)

	//Store the battery in ledger
	batteryAsBytes, _ := json.Marshal(battery)
	err = stub.PutState(battery.Id_battery, batteryAsBytes)
	if err != nil {
		fmt.Println("Could not store battery")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	//Store the fleet in ledger
	fleetAsBytes, _ := json.Marshal(fleet)
	err = stub.PutState(fleet.Id_fleet, fleetAsBytes)
	if err != nil {
		fmt.Println("Could not store fleet")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end DeallocateBatteryFromFleet")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

// TransferBatteryBetweenFleets transfers a battery from one fleet to another.
func TransferBatteryBetweenFleets(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting TransferBatteryBetweenFleets")

	if len(args) != 3 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 3"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	id_battery := args[0]
	source_fleet_id := args[1]
	target_fleet_id := args[2]

	battery, err := get_battery(stub, id_battery)
	if err != nil {
		fmt.Println("battery not found in Blockchain - " + id_battery)
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	source_fleet, err := get_fleet(stub, source_fleet_id)
	if err != nil {
		fmt.Println("Source fleet not found in Blockchain - " + source_fleet_id)
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	target_fleet, err := get_fleet(stub, target_fleet_id)
	if err != nil {
		fmt.Println("Target fleet not found in Blockchain - " + target_fleet_id)
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	if battery.AllocatedToFleet != source_fleet_id {
		return shim.Error(formatError("NA", "NA", "Battery is not allocated to the source fleet - "+source_fleet_id))
	}

	if source_fleet.Company != target_fleet.Company {
		return shim.Error(formatError("NA", "NA", "Battery can only be transferred between fleets on the same company - "+source_fleet.Company))
	}

	// Update the battery's fleet allocation
	battery.AllocatedToFleet = target_fleet_id

	// Decrement Number of Batteries for source fleet by 1 as the battery is transferred
	sourceFleetTotalBatteriesInt, err := sub(int(source_fleet.TotalBatteries), 1)
	source_fleet.TotalBatteries = uint64(sourceFleetTotalBatteriesInt)

	// Increment Number of Batteries for target fleet by 1 as the battery is transferred
	targetFleetTotalBatteriesInt, err := add(int(target_fleet.TotalBatteries), 1)
	target_fleet.TotalBatteries = uint64(targetFleetTotalBatteriesInt)

	//Store the battery in ledger
	batteryAsBytes, _ := json.Marshal(battery)
	err = stub.PutState(battery.Id_battery, batteryAsBytes)
	if err != nil {
		fmt.Println("Could not store battery")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	//Store the source fleet in ledger
	sourceFleetAsBytes, _ := json.Marshal(source_fleet)
	err = stub.PutState(source_fleet.Id_fleet, sourceFleetAsBytes)
	if err != nil {
		fmt.Println("Could not store source fleet")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	//Store the target fleet in ledger
	targetFleetAsBytes, _ := json.Marshal(target_fleet)
	err = stub.PutState(target_fleet.Id_fleet, targetFleetAsBytes)
	if err != nil {
		fmt.Println("Could not store target fleet")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end TransferBatteryBetweenFleets")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

// TransferBatteryFromSSToUser transfers charged battery from SS to user for consumption.
func TransferBatteryFromSSToUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting TransferBatteryFromSSToUser")

	if len(args) != 7 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 7"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	id_battery := args[0]

	battery, err := get_battery(stub, id_battery)
	if err != nil {
		fmt.Println("battery not found in Blockchain - " + id_battery)
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	if battery.Status == "In_Use" {
		return shim.Error(formatError("NA", "NA", "Battery already in use - "+id_battery))
	}

	if battery.DockedStation == "" {
		return shim.Error(formatError("NA", "NA", "Battery is not docked to any swapping station- "+id_battery))
	}

	//Swapping station (docked station) needs to be initialized
	swappingStation, err := get_swappingStation(stub, battery.DockedStation)
	if err != nil {
		fmt.Println("The Swapping Station doesnt exist - " + battery.DockedStation)
		return shim.Error(formatError("NA", "NA", "The Swapping Station doesnt exist - "+battery.DockedStation))
	}

	SoC_uint64, _ := strconv.ParseUint(args[1], 10, 0)
	battery.SoC = uint8(SoC_uint64)
	SoH_uint64, _ := strconv.ParseUint(args[2], 10, 0)
	battery.SoH = uint8(SoH_uint64)
	EnergyContent_float64, _ := strconv.ParseFloat(args[3], 0)
	battery.EnergyContent = float32(EnergyContent_float64)
	Cdc_uint64, _ := strconv.ParseUint(args[4], 10, 0)
	battery.Cdc = uint16(Cdc_uint64)
	battery.User = args[5]
	charge_float64, _ := strconv.ParseFloat(args[6], 32)
	battery.Status = "In_Use"

	// Decrement Number of Active Batteries for SS by 1 as battery is rented for use
	swappingStationActiveBatteriesInt, err := sub(int(swappingStation.ActiveBatteries), 1)
	swappingStation.ActiveBatteries = uint64(swappingStationActiveBatteriesInt)

	// Decrement Number of Total Batteries for SS by 1 as battery is rented for use
	swappingStationTotalBatteriesInt, err := sub(int(swappingStation.TotalBatteries), 1)
	swappingStation.TotalBatteries = uint64(swappingStationTotalBatteriesInt)

	//User (Battery User) needs to be initialized
	user, err := get_user(stub, battery.User)
	if err != nil {
		fmt.Println("The User doesnt exist - " + battery.User)
		return shim.Error(formatError("NA", "NA", "The User doesnt exist - "+battery.User))
	}

	if swappingStation.Company != user.Company {
		return shim.Error(formatError("NA", "NA", "Battery can only be swapped between swapping stations of same company - "+id_battery))
	}

	// Mark rented battery in the user object
	user.RentedBattery = id_battery

	//Escrow Amount from User to the Battery for settlement post use
	user.Wallet, err = subFloat(user.Wallet, float32(charge_float64))
	if err != nil {
		fmt.Println("Could not deduct money from User")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	battery.EscrowedAmount = float32(charge_float64)

	//Store the battery in ledger
	batteryAsBytes, _ := json.Marshal(battery)
	err = stub.PutState(battery.Id_battery, batteryAsBytes)
	if err != nil {
		fmt.Println("Could not store battery")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	//Store the swapping station in ledger
	swappingStationAsBytes, _ := json.Marshal(swappingStation)
	err = stub.PutState(swappingStation.Id_swappingStation, swappingStationAsBytes)
	if err != nil {
		fmt.Println("Could not store swapping station")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	//Store the user in ledger
	userAsBytes, _ := json.Marshal(user)
	err = stub.PutState(user.Id_user, userAsBytes)
	if err != nil {
		fmt.Println("Could not store user")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end TransferBatteryFromSSToUser battery ")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

// TransferBatteryFromUserToSS transfers discharged (used) battery from user to SS for recharging post use.
func TransferBatteryFromUserToSS(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting TransferBatteryFromUserToSS")

	if len(args) != 6 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 6"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	id_battery := args[0]

	battery, err := get_battery(stub, id_battery)
	if err != nil {
		fmt.Println("battery not found in Blockchain - " + id_battery)
		return shim.Error(formatError("NA", "NA", "battery not found in Blockchain - "+id_battery))
	}

	//User (Battery User) needs to be initialized
	user, err := get_user(stub, battery.User)
	if err != nil {
		fmt.Println("The User doesnt exist - " + battery.User)
		return shim.Error(formatError("NA", "NA", "The User doesnt exist - "+battery.User))
	}

	if battery.Status == "In_Service" {
		return shim.Error(formatError("NA", "NA", "Battery already in service- "+id_battery))
	}

	// Calculate energy consumption during battery use for refund and payment settlement
	currentEnergyContent_float64, _ := strconv.ParseFloat(args[3], 32)
	initialEnergy := battery.EnergyContent
	usedEnergy_float32, err := subFloat(initialEnergy, float32(currentEnergyContent_float64))
	consumptionCharge := usedEnergy_float32 / initialEnergy * float32(battery.EscrowedAmount)

	//Calculate Refund
	refundPercentage := float32(.9)
	refundAmountFloat32, err := subFloat(battery.EscrowedAmount, consumptionCharge)
	refundAmount := refundPercentage * refundAmountFloat32

	//Calculate Fee Split
	truePowerNetworkId, err := get_truePowerNetworkId(stub)
	if err != nil {
		fmt.Println("Could not get TruePower network id")
		return shim.Error(formatError("NA", "NA", "Could not get TruePower network id"))
	}

	//Calculate Fee Split
	oldDockedStation, err := get_swappingStation(stub, battery.DockedStation)
	if err != nil {
		fmt.Println("Could not get Old Docked Station")
		return shim.Error(formatError("NA", "NA", "Could not get Old Docked Station"))
	}

	var percentageNetwork, percentageTruePower float32

	if battery.Id_Network == truePowerNetworkId && oldDockedStation.Id_Network == truePowerNetworkId {
		percentageTruePower = 1
		percentageNetwork = 0
	} else {
		percentageNetwork = float32(0.85)
		percentageTruePower = float32(0.15)
		if battery.Id_Network != truePowerNetworkId {
			percentageNetwork += float32(0.5)
			percentageNetwork -= float32(0.5)
		}
		if oldDockedStation.Id_Network != truePowerNetworkId {
			percentageNetwork += float32(0.5)
			percentageNetwork -= float32(0.5)
		}
	}

	amountNetwork := percentageNetwork * consumptionCharge
	amountTruePower := percentageTruePower*consumptionCharge + (1-refundPercentage)*refundAmountFloat32

	//Complete refund payment
	user.Wallet, err = addFloat(user.Wallet, refundAmount)
	if err != nil {
		fmt.Println("Could not refund User")
		return shim.Error(formatError("NA", "NA", "Could not refund User"))
	}

	if percentageNetwork != 0 {
		//Complete settlement to Third Party Network
		err = _payToSSNetwork(stub, battery.Id_Network, amountNetwork)
		if err != nil {
			fmt.Println("Could not send fees to swapping station" + battery.Id_Network)
			return shim.Error(formatError("NA", "NA", "Could not send fees to swapping station"+battery.Id_Network))
		}
	}

	//Complete settlement to TruePower
	err = _payToSSNetwork(stub, truePowerNetworkId, amountTruePower)
	if err != nil {
		fmt.Println("Could not send fees to swapping station" + truePowerNetworkId)
		return shim.Error(formatError("NA", "NA", "Could not send fees to swapping station TruePower"+truePowerNetworkId))
	}

	SoC_uint64, _ := strconv.ParseUint(args[1], 10, 0)
	battery.SoC = uint8(SoC_uint64)
	SoH_uint64, _ := strconv.ParseUint(args[2], 10, 0)
	battery.SoH = uint8(SoH_uint64)
	battery.EnergyContent = float32(currentEnergyContent_float64)
	Cdc_uint64, _ := strconv.ParseUint(args[4], 10, 0)
	battery.Cdc = uint16(Cdc_uint64)
	battery.User = ""
	battery.DockedStation = args[5]
	battery.EscrowedAmount = 0
	battery.Status = "In_Service"

	//Swapping station (docked station) needs to be initialized
	swappingStation, err := get_swappingStation(stub, battery.DockedStation)
	if err != nil {
		fmt.Println("The Swapping Station doesnt exist - " + battery.Id_Network)
		return shim.Error(formatError("NA", "NA", "The Swapping Station doesnt exist - "+battery.Id_Network))
	}

	if swappingStation.Company != battery.Company {
		return shim.Error(formatError("NA", "NA", "Battery can only be swapped between swapping stations of same company - "+id_battery))
	}

	// Increment Number of Discharged Batteries for swapping station by 1 as used battery is added
	swappingStationDischargedBatteriesInt, err := add(int(swappingStation.DischargedBatteries), 1)
	swappingStation.DischargedBatteries = uint64(swappingStationDischargedBatteriesInt)

	// Increment Number of Total Batteries for swapping station by 1 as used battery is added
	swappingStationTotalBatteriesInt, err := add(int(swappingStation.TotalBatteries), 1)
	swappingStation.TotalBatteries = uint64(swappingStationTotalBatteriesInt)

	user.RentedBattery = ""

	//Store the battery in ledger
	batteryAsBytes, _ := json.Marshal(battery)
	err = stub.PutState(battery.Id_battery, batteryAsBytes)
	if err != nil {
		fmt.Println("Could not store battery" + battery.Id_battery)
		return shim.Error(formatError("NA", "NA", "Could not store battery"+battery.Id_battery))
	}

	//Store the swapping station in ledger
	swappingStationAsBytes, _ := json.Marshal(swappingStation)
	err = stub.PutState(swappingStation.Id_swappingStation, swappingStationAsBytes)
	if err != nil {
		fmt.Println("Could not store swapping station" + swappingStation.Id_swappingStation)
		return shim.Error(formatError("NA", "NA", "Could not store swapping station"+swappingStation.Id_swappingStation))
	}

	//Store the user in ledger
	userAsBytes, _ := json.Marshal(user)
	err = stub.PutState(user.Id_user, userAsBytes)
	if err != nil {
		fmt.Println("Could not store user" + user.Id_user)
		return shim.Error(formatError("NA", "NA", "Could not store user"+user.Id_user))
	}

	fmt.Println("- end transferLnd2BSS battery ")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

// ReturnBatteryFromService marks battery on SS as ready to use once charged.
func ReturnBatteryFromService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting ReturnBatteryFromService")

	if len(args) != 5 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 5"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	id_battery := args[0]

	battery, err := get_battery(stub, id_battery)
	if err != nil {
		fmt.Println("battery not found in Blockchain - " + id_battery)
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	if battery.Status == "Available" {
		return shim.Error(formatError("NA", "NA", "Battery already returned from service- "+id_battery))
	}

	if battery.DockedStation == "" {
		return shim.Error(formatError("NA", "NA", "Battery is not docked to any swapping station- "+id_battery))
	}

	//Swapping station (docked station) needs to be initialized
	swappingStation, err := get_swappingStation(stub, battery.DockedStation)
	if err != nil {
		fmt.Println("The Swapping Station doesnt exist - " + battery.Id_Network)
		return shim.Error(formatError("NA", "NA", "The Swapping Station doesnt exist - "+battery.Id_Network))
	}

	// Increment Number of Active Batteries for SS by 1 as battery is recharged
	swappingStationActiveBatteriesInt, err := add(int(swappingStation.ActiveBatteries), 1)
	swappingStation.ActiveBatteries = uint64(swappingStationActiveBatteriesInt)

	// Decrement Number of Discharged Batteries for SS by 1 as battery is recharged
	swappingStationDischargedBatteriesInt, err := sub(int(swappingStation.DischargedBatteries), 1)
	swappingStation.DischargedBatteries = uint64(swappingStationDischargedBatteriesInt)

	SoC_uint64, _ := strconv.ParseUint(args[1], 10, 0)
	battery.SoC = uint8(SoC_uint64)
	SoH_uint64, _ := strconv.ParseUint(args[2], 10, 0)
	battery.SoH = uint8(SoH_uint64)
	EnergyContent_float64, _ := strconv.ParseFloat(args[3], 0)
	battery.EnergyContent = float32(EnergyContent_float64)
	Cdc_uint64, _ := strconv.ParseUint(args[4], 10, 0)
	battery.Cdc = uint16(Cdc_uint64)
	battery.Status = "Available"

	//Store the battery in ledger
	batteryAsBytes, _ := json.Marshal(battery)
	err = stub.PutState(battery.Id_battery, batteryAsBytes)
	if err != nil {
		fmt.Println("Could not store battery")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	//Store the swapping station in ledger
	swappingStationAsBytes, _ := json.Marshal(swappingStation)
	err = stub.PutState(swappingStation.Id_swappingStation, swappingStationAsBytes)
	if err != nil {
		fmt.Println("Could not store swapping station")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end returnFromService battery ")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

/* -------------------------------------------------------------------------- */
/*                              Payment Methods                               */
/* -------------------------------------------------------------------------- */

// PayToSSNetwork adds specific amount to Battery Network wallet balance.
func PayToSSNetwork(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting PayToSSNetwork")

	if len(args) != 2 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 2"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	id_Network := args[0]
	charge, _ := strconv.ParseFloat(args[1], 32)

	err = _payToSSNetwork(stub, id_Network, float32(charge))

	if err != nil {
		fmt.Println("Could not store Battery Network")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end payToWallet Battery Network")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

// PayToUser adds specific amount to user wallet balance.
func RechargeUserWallet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting PayToUser")

	if len(args) != 2 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 2"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	id_user := args[0]
	charge, _ := strconv.ParseFloat(args[1], 32)

	err = _payToUser(stub, id_user, float32(charge))

	if err != nil {
		fmt.Println("Could not store user")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end payToWallet user ")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

// PayToUser adds specific amount to user wallet balance.
func PayToUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting PayToUser")

	if len(args) != 2 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 2"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	id_user := args[0]
	charge, _ := strconv.ParseFloat(args[1], 32)

	err = _payToUser(stub, id_user, float32(charge))

	if err != nil {
		fmt.Println("Could not store user")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end payToWallet user ")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

// PayFromUser deducts specific amount from user wallet balance.
func PayFromUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting PayFromUser")

	if len(args) != 2 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 2"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	id_user := args[0]
	charge, _ := strconv.ParseFloat(args[1], 32)

	err = _payFromUser(stub, id_user, float32(charge))

	if err != nil {
		fmt.Println("Could not store user")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end payFromWallet user ")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

/* -------------------------------------------------------------------------- */
/*                         Battery Edge Case Handling                         */
/* -------------------------------------------------------------------------- */

// MarkBatteryStolen marks a specific battery stolen to remove it from inventory.
func MarkBatteryStolen(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting MarkBatteryStolen")

	if len(args) != 1 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 1"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	id_battery := args[0]

	battery, err := get_battery(stub, id_battery)
	if err != nil {
		fmt.Println("battery not found in Blockchain - " + id_battery)
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	if battery.Status == "Stolen" {
		return shim.Error(formatError("NA", "NA", "Battery already marked stolen - "+id_battery))
	}

	battery.Status = "Stolen"

	//Store the battery in ledger
	batteryAsBytes, _ := json.Marshal(battery)
	err = stub.PutState(battery.Id_battery, batteryAsBytes)
	if err != nil {
		fmt.Println("Could not store battery")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end markStolen battery ")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

// MarkBatteryError marks an error on against a specific battery id.
func MarkBatteryError(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting MarkBatteryError")

	if len(args) != 1 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 1"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	id_battery := args[0]

	battery, err := get_battery(stub, id_battery)
	if err != nil {
		fmt.Println("battery not found in Blockchain - " + id_battery)
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	if battery.Status == "Error" {
		return shim.Error(formatError("NA", "NA", "Battery already marked error - "+id_battery))
	}

	battery.Status = "Error"

	//Store the battery in ledger
	batteryAsBytes, _ := json.Marshal(battery)
	err = stub.PutState(battery.Id_battery, batteryAsBytes)
	if err != nil {
		fmt.Println("Could not store battery")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end markError battery ")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}

// MarkBatteryExpired will mark a specific battery as expired against a specific battery id.
func MarkBatteryExpired(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting MarkBatteryExpired")

	if len(args) != 1 {
		return shim.Error(formatError("NA", "NA", "Incorrect number of arguments. Expecting 1"))
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	id_battery := args[0]

	battery, err := get_battery(stub, id_battery)
	if err != nil {
		fmt.Println("battery not found in Blockchain - " + id_battery)
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	if battery.Status == "Expired" {
		return shim.Error(formatError("NA", "NA", "Battery already marked expired - "+id_battery))
	}

	battery.Status = "Expired"

	//Swapping station (docked station) needs to be initialized
	swappingStation, err := get_swappingStation(stub, battery.DockedStation)
	if err != nil {
		fmt.Println("The Swapping Station doesnt exist - " + battery.Id_Network)
		return shim.Error(formatError("NA", "NA", "The Swapping Station doesnt exist - "+battery.Id_Network))
	}

	// Increment Number of Expired Batteries for swapping station by 1 as used battery is added
	swappingStationBatteriesInt, err := add(int(swappingStation.ExpiredBatteries), 1)
	swappingStation.ExpiredBatteries = uint64(swappingStationBatteriesInt)

	//Swapping network (owner network) needs to be initialized
	ssNetwork, err := get_ssNetwork(stub, battery.DockedStation)
	if err != nil {
		fmt.Println("The Swapping Station doesnt exist - " + battery.Id_Network)
		return shim.Error(formatError("NA", "NA", "The Swapping Station doesnt exist - "+battery.Id_Network))
	}

	// Increment Number of Expired Batteries for swapping station by 1 as used battery is added
	ssNetworkBatteriesInt, err := add(int(ssNetwork.ExpiredBatteries), 1)
	ssNetwork.ExpiredBatteries = uint64(ssNetworkBatteriesInt)

	//Store the battery in ledger
	batteryAsBytes, _ := json.Marshal(battery)
	err = stub.PutState(battery.Id_battery, batteryAsBytes)
	if err != nil {
		fmt.Println("Could not store battery")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	//Store the swapping station in ledger
	swappingStationAsBytes, _ := json.Marshal(swappingStation)
	err = stub.PutState(swappingStation.Id_swappingStation, swappingStationAsBytes)
	if err != nil {
		fmt.Println("Could not store swapping station")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	//Store the swapping station in ledger
	ssNetworkAsBytes, _ := json.Marshal(ssNetwork)
	err = stub.PutState(ssNetwork.Id_Network, ssNetworkAsBytes)
	if err != nil {
		fmt.Println("Could not store swapping station")
		return shim.Error(formatError("NA", "NA", err.Error()))
	}

	fmt.Println("- end markExpired battery ")

	txID := stub.GetTxID()
	return shim.Success(formatSuccess("NA", "NA", txID))
}
