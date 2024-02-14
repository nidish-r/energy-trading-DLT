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
	"time"

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
		return shim.Error("Incorrect number of arguments. Expecting 2. key of the variable and value to set")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the ledger
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end write")
	return shim.Success(nil)
}

/* -------------------------------------------------------------------------- */
/*                             User Write Methods                             */
/* -------------------------------------------------------------------------- */

func UpdateUserProfile(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting UpdateUserProfile")

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	err := sanitize_arguments(args)
	if err != nil {
		return shim.Error("Invalid argument: " + err.Error())
	}

	userID := args[0]
	existingUserAsBytes, err := stub.GetState(userID)

	var user User
	if err != nil || existingUserAsBytes == nil {
		// New user creation
		user.CreatedOn = time.Now().Unix()
		user.UpdatedOn = user.CreatedOn
	} else {
		// Existing user update
		err = json.Unmarshal(existingUserAsBytes, &user)
		if err != nil {
			return shim.Error("Failed to unmarshal user: " + err.Error())
		}
		user.UpdatedOn = time.Now().Unix()
	}

	user.ID = userID
	if err != nil {
		return shim.Error("Failed to convert user ID: " + err.Error())
	}
	user.Category = args[1]
	user.Location = args[2]
	user.MeterID = args[3]
	user.Source = args[4]
	isAdminStr := args[5]
	isAdminBool, err := strconv.ParseBool(isAdminStr)
	if err != nil {
		return shim.Error("Failed to parse IsAdmin Bool: " + err.Error())
	} else {
		user.IsAdmin = isAdminBool
	}

	// Store the user in ledger
	userAsBytes, _ := json.Marshal(user)
	err = stub.PutState(user.ID, userAsBytes)
	if err != nil {
		return shim.Error("Could not store user: " + err.Error())
	}

	if existingUserAsBytes == nil {
		fmt.Println("- end CreateUser")
		return shim.Success(nil)
	} else {
		fmt.Println("- end UpdateUser")
		return shim.Success([]byte(stub.GetTxID()))
	}
}

func SignPlatformContract(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting SignPlatformContract")

	// We are assuming that the only argument is the user ID.
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1 (UserID)")
	}

	// Check if user exists.
	userID := args[0]
	existingUserAsBytes, err := stub.GetState(userID)
	if err != nil || existingUserAsBytes == nil {
		return shim.Error("User with ID " + userID + " not found")
	}

	// Creating a new platform contract for the user.
	var contract PlatformContract
	contract.UserID = userID
	if err != nil {
		return shim.Error("Failed to convert user ID: " + err.Error())
	}
	contract.CreatedOn = time.Now().Unix()
	contract.UpdatedOn = contract.CreatedOn

	// Store the contract in the ledger using a composite key for uniqueness.
	// This will use "PlatformContract" as a prefix followed by the user ID.
	contractKey := "PlatformContract_" + userID

	contractAsBytes, _ := json.Marshal(contract)
	err = stub.PutState(contractKey, contractAsBytes)
	if err != nil {
		return shim.Error("Could not store platform contract: " + err.Error())
	}

	fmt.Println("- end SignPlatformContract")
	return shim.Success([]byte(stub.GetTxID()))
}

/* -------------------------------------------------------------------------- */
/*                              Payment Methods                               */
/* -------------------------------------------------------------------------- */

func RecordPayment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting RecordPayment")

	// Basic argument validation. We expect 12 arguments.
	if len(args) != 12 {
		return shim.Error("Incorrect number of arguments. Expecting 12.")
	}

	// Extracting required arguments.
	paymentID := args[0]
	paymentType := args[1]
	totalAmount, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return shim.Error("Failed to parse total amount: " + err.Error())
	}
	userID := args[3]
	if err != nil {
		return shim.Error("Failed to parse user ID: " + err.Error())
	}
	paymentDetailID := args[4]
	debitedFrom := args[5]
	creditedTo := args[6]
	totalUnitCost, _ := strconv.ParseFloat(args[7], 64)
	platformFee, _ := strconv.ParseFloat(args[8], 64)
	tokenAmount, _ := strconv.ParseFloat(args[9], 64)
	bidRefundAmount, _ := strconv.ParseFloat(args[10], 64)
	platformFeeRefundAmount, _ := strconv.ParseFloat(args[11], 64)
	penaltyFromSeller, _ := strconv.ParseFloat(args[12], 64)

	// Create PaymentDetail entry.
	pd := PaymentDetail{
		ID:                      paymentDetailID, // Unique ID based on current timestamp.
		DebitedFrom:             debitedFrom,
		CreditedTo:              creditedTo,
		TotalUnitCost:           totalUnitCost,
		PlatformFee:             platformFee,
		TokenAmount:             tokenAmount,
		BidRefundAmount:         bidRefundAmount,
		PlatformFeeRefundAmount: platformFeeRefundAmount,
		TokenAmountRefund:       penaltyFromSeller,
		PenaltyFromSeller:       penaltyFromSeller,
	}

	// Store the PaymentDetail in the ledger.
	pdAsBytes, _ := json.Marshal(pd)
	err = stub.PutState("PaymentDetail_"+pd.ID, pdAsBytes)
	if err != nil {
		return shim.Error("Could not store payment detail: " + err.Error())
	}

	// Create and store the Payment entry, using the PaymentDetail ID.
	p := Payment{
		CreatedOn:       time.Now().Unix(),
		ID:              paymentID,
		PaymentDetailID: pd.ID,
		PaymentType:     paymentType,
		TotalAmount:     totalAmount,
		UserID:          userID,
	}

	pAsBytes, _ := json.Marshal(p)
	err = stub.PutState("Payment_"+paymentID, pAsBytes)
	if err != nil {
		return shim.Error("Could not store payment: " + err.Error())
	}

	fmt.Println("- end RecordPayment")
	return shim.Success([]byte(stub.GetTxID()))
}

/* -------------------------------------------------------------------------- */
/*                            Energy Bid  Methods                             */
/* -------------------------------------------------------------------------- */

func RegisterOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting RegisterOrder")

	// We expect 12 arguments.
	if len(args) != 12 {
		return shim.Error("Incorrect number of arguments. Expecting 12.")
	}

	// Parsing ID first to check existence.
	orderID := args[2]

	// Check if order with given ID already exists.
	existingOrderAsBytes, err := stub.GetState("Order_" + orderID)
	if err != nil {
		return shim.Error("Error accessing state: " + err.Error())
	}

	var order Order
	if existingOrderAsBytes != nil {
		// Order exists, so we will update it.
		err = json.Unmarshal(existingOrderAsBytes, &order)
		if err != nil {
			return shim.Error("Failed to unmarshal existing order: " + err.Error())
		}
	} else {
		// Order doesn't exist, so we will create a new one.
		order.CreatedOn = time.Now().Unix()
		order.ID = orderID
	}

	bidMatchID := args[0]

	bidStatus := args[1]

	// BidStatus check
	if existingOrderAsBytes == nil {
		if bidStatus != "BidCreated" && bidStatus != "BidAccepted" {
			return shim.Error("Invalid BidStatus provided for new Order. It should be BidCreated or BidAccepted.")
		}
	}

	onMarketPrice := args[3]

	orderCost, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		return shim.Error("Failed to parse OrderCost: " + err.Error())
	}

	paymentID := args[5]
	if err != nil {
		return shim.Error("Failed to parse PaymentID: " + err.Error())
	}

	slotID := args[6]
	totalQuantity, err := strconv.ParseInt(args[7], 10, 64)
	if err != nil {
		return shim.Error("Failed to parse TotalQuantity: " + err.Error())
	}

	unitCost, err := strconv.ParseFloat(args[8], 64)
	if err != nil {
		return shim.Error("Failed to parse UnitCost: " + err.Error())
	}

	userID := args[9]

	slotExecDate, err := strconv.ParseInt(args[10], 10, 64) // Add this line to parse SlotExecDate
	if err != nil {
		return shim.Error("Failed to parse SlotExecDate: " + err.Error())
	}

	action := args[11]
	if err != nil {
		return shim.Error("Failed to parse action: " + err.Error())
	}

	// Assign parsed values to the order struct
	order.BidMatchID = bidMatchID
	order.BidStatus = bidStatus
	order.OnMarketPrice = onMarketPrice
	order.OrderCost = orderCost
	order.PaymentID = paymentID
	order.SlotID = slotID
	order.TotalQuantity = totalQuantity
	order.UnitCost = unitCost
	order.UpdatedOn = time.Now().Unix()
	order.UserID = userID
	order.SlotExecDate = slotExecDate // Set the SlotExecDate
	order.UserAction = action

	// Store the order back in the ledger.
	orderAsBytes, _ := json.Marshal(order)
	err = stub.PutState("Order_"+order.ID, orderAsBytes)

	fmt.Println("- end RegisterOrder")
	return shim.Success([]byte(stub.GetTxID()))
}

func ProcessBidMatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting ProcessBidMatch")

	// We expect 10 arguments.
	if len(args) != 11 {
		return shim.Error("Incorrect number of arguments. Expecting 11.")
	}

	// Parsing ID first to check existence.
	bidMatchID := args[6]

	// Check if EnergyBid with the given ID already exists.
	existingBidMatchAsBytes, err := stub.GetState("BidMatch_" + bidMatchID)
	if err != nil {
		return shim.Error("Error accessing state: " + err.Error())
	}

	var bidMatch BidMatch
	if existingBidMatchAsBytes != nil {
		// BidMatch exists, so we will update it.
		err = json.Unmarshal(existingBidMatchAsBytes, &bidMatch)
		if err != nil {
			return shim.Error("Failed to unmarshal existing BidMatch: " + err.Error())
		}
	} else {
		// BidMatch doesn't exist, so we will create a new one.
		bidMatch.ID = bidMatchID
	}

	bidMatchTms, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return shim.Error("Failed to parse BidStatus: " + err.Error())
	}
	bidSlot := args[1]
	bidStatus := args[2]

	bidUnitPrice, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return shim.Error("Failed to parse BidUnitPrice: " + err.Error())
	}
	buyerUserID := args[4]
	deliveredBidUnits, err := strconv.ParseFloat(args[5], 64)
	if err != nil {
		return shim.Error("Failed to parse DeliveredBidUnits: " + err.Error())
	}
	originalBidUnits, err := strconv.ParseFloat(args[7], 64)
	if err != nil {
		return shim.Error("Failed to parse OriginalBidUnits: " + err.Error())
	}
	sellerUserID := args[8]
	transactionBuyID := args[9]
	transactionSellID := args[10]

	// Assign parsed values to bidMatch
	bidMatch.BidMatchTms = bidMatchTms
	bidMatch.BidSlot = bidSlot
	bidMatch.BidStatus = bidStatus
	bidMatch.BidUnitPrice = bidUnitPrice
	bidMatch.BuyerUserId = buyerUserID
	bidMatch.DeliveredBidUnits = deliveredBidUnits
	bidMatch.ID = bidMatchID
	bidMatch.OriginalBidUnits = originalBidUnits
	bidMatch.SellerUserId = sellerUserID
	bidMatch.TransactionBuyID = transactionBuyID
	bidMatch.TransactionSellID = transactionSellID

	// Store the bidMatch back in the ledger.
	bidMatchAsBytes, _ := json.Marshal(bidMatch)
	err = stub.PutState("BidMatch_"+bidMatch.ID, bidMatchAsBytes)
	if err != nil {
		return shim.Error("Could not store BidMatch: " + err.Error())
	}

	fmt.Println("- end ProcessBidMatch")
	//return shim.Success(nil)
	return shim.Success([]byte(stub.GetTxID()))
}

func ProcessEnergyBid(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting ProcessEnergyBid")

	// We expect 12 arguments.
	if len(args) != 12 {
		return shim.Error("Incorrect number of arguments. Expecting 12.")
	}

	// Parsing ID first to check existence.
	energyBidID := args[0]

	// Check if EnergyBid with the given ID already exists.
	existingEnergyBidAsBytes, err := stub.GetState("EnergyBid_" + energyBidID)
	if err != nil {
		return shim.Error("Error accessing state: " + err.Error())
	}

	var energyBid EnergyBid
	if existingEnergyBidAsBytes != nil {
		// EnergyBid exists, so we will update it.
		err = json.Unmarshal(existingEnergyBidAsBytes, &energyBid)
		if err != nil {
			return shim.Error("Failed to unmarshal existing EnergyBid: " + err.Error())
		}
	} else {
		// EnergyBid doesn't exist, so we will create a new one.
		energyBid.ID = energyBidID
	}

	bidMatchID := args[1]
	initialBidUnits, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return shim.Error("Failed to parse InitialBidUnits: " + err.Error())
	}
	acceptedBidUnits, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		return shim.Error("Failed to parse AcceptedBidUnits: " + err.Error())
	}
	buyerMeterUnit, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		return shim.Error("Failed to parse BuyerMeterUnit: " + err.Error())
	}
	sellerMeterUnit, err := strconv.ParseFloat(args[5], 64)
	if err != nil {
		return shim.Error("Failed to parse SellerMeterUnit: " + err.Error())
	}
	buyerBroughtUnitFromSeller, err := strconv.ParseFloat(args[6], 64)
	if err != nil {
		return shim.Error("Failed to parse BuyerBroughtUnitFromSeller: " + err.Error())
	}
	sellerSoldUnitToBuyer, err := strconv.ParseFloat(args[7], 64)
	if err != nil {
		return shim.Error("Failed to parse SellerSoldUnitToBuyer: " + err.Error())
	}
	sellerSoldUnitToGrid, err := strconv.ParseFloat(args[8], 64)
	if err != nil {
		return shim.Error("Failed to parse SellerSoldUnitToGrid: " + err.Error())
	}
	buyerSoldUnitToGrid, err := strconv.ParseFloat(args[9], 64)
	if err != nil {
		return shim.Error("Failed to parse BuyerSoldUnitToGrid: " + err.Error())
	}
	buyerBroughtUnitFromGrid, err := strconv.ParseFloat(args[10], 64)
	if err != nil {
		return shim.Error("Failed to parse BuyerBroughtUnitFromGrid: " + err.Error())
	}
	reason := args[11]

	// Assign parsed values to energyBid
	energyBid.BidMatchID = bidMatchID
	energyBid.InitialBidUnits = initialBidUnits
	energyBid.AcceptedBidUnits = acceptedBidUnits
	energyBid.BuyerMeterUnit = buyerMeterUnit
	energyBid.SellerMeterUnit = sellerMeterUnit
	energyBid.BuyerBroughtUnitFromSeller = buyerBroughtUnitFromSeller
	energyBid.SellerSoldUnitToBuyer = sellerSoldUnitToBuyer
	energyBid.SellerSoldUnitToGrid = sellerSoldUnitToGrid
	energyBid.BuyerSoldUnitToGrid = buyerSoldUnitToGrid
	energyBid.BuyerBroughtUnitFromGrid = buyerBroughtUnitFromGrid
	energyBid.Reason = reason
	energyBid.CreatedOn = time.Now().Unix()

	// Store the energyBid back in the ledger.
	energyBidAsBytes, _ := json.Marshal(energyBid)
	err = stub.PutState("EnergyBid_"+energyBid.ID, energyBidAsBytes)
	if err != nil {
		return shim.Error("Could not store EnergyBid: " + err.Error())
	}

	fmt.Println("- end ProcessEnergyBid")
	//return shim.Success(nil)
	return shim.Success([]byte(stub.GetTxID()))
}
