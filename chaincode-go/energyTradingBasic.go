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
// User Definitions - The ledger with user
// ============================================================================================================================

// User represents the schema for the user table.
type User struct {
	ID        string `json:"id"`
	Category  string `json:"category"`
	CreatedOn int64  `json:"createdOn"`
	IsAdmin   bool   `json:"isAdmin"`
	Location  string `json:"location"`
	MeterID   string `json:"meterId"`
	Source    string `json:"source"`
	UpdatedOn int64  `json:"updatedOn"`
}

type PlatformContract struct {
	UserID             string `json:"userId"`
	SignedContractHash string `json:"signedContractHash"`
	CreatedOn          int64  `json:"createdOn"`
	UpdatedOn          int64  `json:"updatedOn"`
}

type TradingContract struct {
	UserID             string `json:"userId"`
	BidStatus          string `json:"bidStatus"`
	SignedContractHash string `json:"signedContractHash"`
	CreatedOn          int64  `json:"createdOn"`
	UpdatedOn          int64  `json:"updatedOn"`
}

// ============================================================================================================================
// Trading Definitions - The ledger with user
// ============================================================================================================================

// Order captures the details of an energy buy or sell bid.
// It includes attributes like total quantity, unit cost, and the total order cost.
// Struct fields are arranged alphabetically to ensure determinism across languages.
// Note: While Golang maintains field order when marshaling to JSON, it doesn't auto-sort them.
type Order struct {
	BidMatchID    string  `json:"bidMatchId"`
	BidStatus     string  `json:"bidStatus"`
	CreatedOn     int64   `json:"createdOn"`
	ID            string  `json:"id"`
	OnMarketPrice string  `json:"onMarketPrice"`
	OrderCost     float64 `json:"orderCost"`
	PaymentID     string  `json:"paymentId"`
	SlotID        string  `json:"slotId"`
	SlotExecDate  int64   `json:"slotExecDate"`
	TotalQuantity int64   `json:"totalQuantity"`
	UnitCost      float64 `json:"unitCost"`
	UpdatedOn     int64   `json:"updatedOn"`
	UserAction    string  `json:"action"`
	UserID        string  `json:"userId"`
}

// BidMatch records the details of a matched bid in the energy market.
// Struct fields are alphabetically ordered for cross-language determinism.
type BidMatch struct {
	BidMatchTms       int64   `json:"bidMatchTms"`
	BidSlot           string  `json:"bidSlot"`
	BidStatus         string  `json:"bidStatus"`
	BidUnitPrice      int64   `json:"bidUnitPrice"`
	BuyerUserId       string  `json:"buyerUserId"`
	DeliveredBidUnits float64 `json:"deliveredBidUnits"`
	ID                string  `json:"id"`
	OriginalBidUnits  float64 `json:"originalBidUnits"`
	SellerUserId      string  `json:"sellerUserId"`
	TransactionBuyID  string  `json:"transactionBuyId"`
	TransactionSellID string  `json:"transactionSellId"`
}

// EnergyBid records the details of a executed bid in the energy market.
// Struct fields are alphabetically ordered for cross-language determinism.
type EnergyBid struct {
	ID                         string  `json:"id"`
	BidMatchID                 string  `json:"bidMatchId"`
	InitialBidUnits            float64 `json:"initialBidUnits"`
	AcceptedBidUnits           float64 `json:"acceptedBidUnits"`
	BuyerMeterUnit             float64 `json:"buyerMeterUnit"`
	SellerMeterUnit            float64 `json:"sellerMeterUnit"`
	BuyerBroughtUnitFromSeller float64 `json:"buyerBroughtUnitFromSeller"`
	SellerSoldUnitToBuyer      float64 `json:"sellerSoldUnitToBuyer"`
	SellerSoldUnitToGrid       float64 `json:"sellerSoldUnitToGrid"`
	BuyerSoldUnitToGrid        float64 `json:"buyerSoldUnitToGrid"`
	BuyerBroughtUnitFromGrid   float64 `json:"buyerBroughtUnitFromGrid"`
	Reason                     string  `json:"reason"`
	CreatedOn                  int64   `json:"createdOn"`
}

// Payment logs transaction details for energy market payments.
// Struct fields are alphabetically ordered for cross-language determinism.
type Payment struct {
	BidMatchID      string  `json:"bidMatchId"`
	CreatedOn       int64   `json:"createdOn"`
	ID              string  `json:"id"`
	PaymentDetailID string  `json:"paymentDetail"`
	PaymentType     string  `json:"paymentType"`
	TotalAmount     float64 `json:"totalAmount"`
	UserID          string  `json:"userId"`
	OrderID         string  `json:"orderId"`
}

// PaymentDetail captures more granular transaction information.
// It includes attributes like the amount refunded, fees applied, and transaction parties.
// Struct fields are arranged alphabetically for consistent representation.
type PaymentDetail struct {
	ID                      string  `json:"id"`
	DebitedFrom             string  `json:"debitedFrom"`
	CreditedTo              string  `json:"creditedTo"`
	TotalUnitCost           float64 `json:"totalUnitCost"`
	PlatformFee             float64 `json:"platformFee"`
	TokenAmount             float64 `json:"tokenAmount"`
	BidRefundAmount         float64 `json:"bidRefundAmount"`
	PlatformFeeRefundAmount float64 `json:"platformFeeRefundAmount"`
	TokenAmountRefund       float64 `json:"tokenAmountRefund"`
	PenaltyFromSeller       float64 `json:"penaltyFromSeller"`
}

// ============================================================================================================================
// Prefix Definitions - For creating composite keys and avoid id overlap (for future use)
// ============================================================================================================================

const BuyBidPrefix = "BuyBid"
const SellBidPrefix = "SellBid"

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
	if function == "Write" {
		return Write(stub, args)
	} else if function == "UpdateUserProfile" {
		return UpdateUserProfile(stub, args)
	} else if function == "SignPlatformContract" {
		return SignPlatformContract(stub, args)
	} else if function == "SignTradingContract" {
		return SignTradingContract(stub, args)
	} else if function == "RecordPayment" {
		return RecordPayment(stub, args)
	} else if function == "RegisterOrder" {
		return RegisterOrder(stub, args)
	} else if function == "ProcessBidMatch" {
		return ProcessBidMatch(stub, args)
	} else if function == "ProcessEnergyBid" {
		return ProcessEnergyBid(stub, args)
	} else if function == "ReadUserProfile" {
		return ReadUserProfile(stub, args)
	} else if function == "ReadPlatformContract" {
		return ReadPlatformContract(stub, args)
	} else if function == "ReadTradingContract" {
		return ReadTradingContract(stub, args)
	} else if function == "ReadPayment" {
		return ReadPayment(stub, args)
	} else if function == "ReadPaymentDetail" {
		return ReadPaymentDetail(stub, args)
	} else if function == "ReadOrder" {
		return ReadOrder(stub, args)
	} else if function == "ReadBidMatch" {
		return ReadBidMatch(stub, args)
	} else if function == "ReadEnergyBid" {
		return ReadEnergyBid(stub, args)
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
