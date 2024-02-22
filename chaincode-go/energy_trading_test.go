package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/stretchr/testify/assert"
)

// func TestWrite(t *testing.T) {
// 	// Create a mock stub
// 	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

// 	// Positive Test Case: Writing a key-value pair
// 	t.Run("Successfully Write key-value pair", func(t *testing.T) {
// 		// Invoke the 'Write' function with a key and a value
// 		response := stub.MockInvoke("1", [][]byte{[]byte("Write"), []byte("TestKey"), []byte("TestValue")})

// 		// Assert that the function completed successfully
// 		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

// 		// Fetch the value for the key from the ledger and check it
// 		value, err := stub.GetState("TestKey")
// 		assert.NoError(t, err, "Error getting value from ledger")
// 		assert.Equal(t, "TestValue", string(value), "Incorrect value retrieved from ledger")
// 	})

// 	// Negative Test Case: Wrong number of arguments
// 	t.Run("Invalid number of arguments", func(t *testing.T) {
// 		response := stub.MockInvoke("2", [][]byte{[]byte("Write"), []byte("OnlyOneArg")})

// 		// Assert that the function did not complete successfully
// 		assert.NotEqual(t, shim.OK, response.GetStatus(), "Function unexpectedly succeeded")
// 	})
// }

func TestUpdateUserProfile(t *testing.T) {
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	// Test Case 1: Successfully Update User Profile
	// Test Case 1: Successfully Update User Profile
	t.Run("Successfully Update User Profile", func(t *testing.T) {
		response := stub.MockInvoke("1", [][]byte{
			[]byte("UpdateUserProfile"),
			[]byte("1"),          // UserID field
			[]byte("Prosumer"),   // Category field
			[]byte("Location 1"), // Location field
			[]byte("MeterId 1"),  // MeterID field
			[]byte("Solar"),      // Energy Source field
			[]byte("true"),       // IsAdmin field
		})

		// Assert the function completed successfully
		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		// Assert the function completed successfully
		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		// Test Case 1.1: Read Updated User Profile
		response = stub.MockInvoke("2", [][]byte{
			[]byte("ReadUserProfile"),
			[]byte("1"), // UserID field
		})

		// Assert the function completed successfully
		assert.Equal(t, int32(shim.OK), response.GetStatus(), "Failed to read user profile")

		// Check if the read user profile matches the updated values
		var readUser User
		err := json.Unmarshal(response.GetPayload(), &readUser)
		assert.NoError(t, err, "Error unmarshalling read user")
		assert.Equal(t, "Location 1", readUser.Location, "Incorrect value retrieved from ledger")
		assert.True(t, readUser.IsAdmin, "Incorrect value retrieved from ledger for IsAdmin")
	})

	// Test Case 2: Incorrect Number of Arguments
	t.Run("Incorrect Number of Arguments", func(t *testing.T) {
		response := stub.MockInvoke("2", [][]byte{
			[]byte("UpdateUserProfile"),
			[]byte("1"),
			[]byte("Prosumer"),
		})

		// Assert the function did not complete successfully
		assert.NotEqual(t, shim.OK, response.GetStatus(), "Function unexpectedly succeeded")
	})
}

func TestUpdateEnterpriseUserProfile(t *testing.T) {
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	// Test Case 1: Successfully Update Enterprise User Profile
	t.Run("Successfully Update Enterprise User Profile", func(t *testing.T) {
		response := stub.MockInvoke("1", [][]byte{
			[]byte("UpdateEnterpriseUserProfile"),
			[]byte("1"),                          // UserID field
			[]byte("Prosumer"),                   // Category field
			[]byte("Location 1"),                 // Location field
			[]byte(`["MeterId 1", "MeterId 2"]`), // MeterIDs field, updated to support multiple meter IDs
			[]byte("Solar"),                      // Energy Source field
			[]byte("true"),                       // IsAdmin field
		})

		// Assert the function completed successfully
		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		// Test Case 1.1: Read Updated Enterprise User Profile
		response = stub.MockInvoke("2", [][]byte{
			[]byte("ReadEnterpriseUserProfile"),
			[]byte("1"), // UserID field
		})

		// Assert the function completed successfully
		assert.Equal(t, int32(shim.OK), response.GetStatus(), "Failed to read enterprise user profile")

		// Check if the read enterprise user profile matches the updated values
		var readUser EnterpriseUser
		err := json.Unmarshal(response.GetPayload(), &readUser)
		assert.NoError(t, err, "Error unmarshalling read enterprise user")
		assert.Equal(t, "Location 1", readUser.Location, "Incorrect value retrieved from ledger for location")
		assert.True(t, readUser.IsAdmin, "Incorrect value retrieved from ledger for IsAdmin")
		assert.Len(t, readUser.MeterIDs, 2, "Incorrect number of MeterIDs after reading")
		assert.Contains(t, readUser.MeterIDs, "MeterId 1", "MeterId 1 not found in MeterIDs after reading")
		assert.Contains(t, readUser.MeterIDs, "MeterId 2", "MeterId 2 not found in MeterIDs after reading")
	})

	// Test Case 2: Incorrect Number of Arguments
	t.Run("Incorrect Number of Arguments", func(t *testing.T) {
		response := stub.MockInvoke("2", [][]byte{
			[]byte("UpdateEnterpriseUserProfile"),
			[]byte("1"),
			[]byte("Prosumer"),
			// Notice we're not providing enough arguments as expected
		})

		// Assert the function did not complete successfully
		assert.NotEqual(t, shim.OK, response.GetStatus(), "Function unexpectedly succeeded")
	})
}

func TestSignPlatformContract(t *testing.T) {
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	key := "12345"
	id := key

	// Creating a dummy user to be used in the tests.
	user := User{ID: id}
	userBytes, _ := json.Marshal(user)
	// Start a transaction
	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	err := stub.PutState(key, userBytes)
	if err != nil {
		t.Fatalf("Failed to put the user into the stub: %s", err.Error())
	}

	// Test Case 1: Successfully Sign Platform Contract
	t.Run("Successfully Sign Platform Contract", func(t *testing.T) {
		response := stub.MockInvoke("1", [][]byte{[]byte("SignPlatformContract"), []byte("12345"), []byte("A7324HAS7234SADF734JSDF")}) // UserID field

		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		contractAsBytes, err := stub.GetState("PlatformContract_12345")
		assert.NoError(t, err, "Error getting value from ledger")

		var contract PlatformContract
		err = json.Unmarshal(contractAsBytes, &contract)
		assert.NoError(t, err, "Error unmarshalling contract")
		assert.Equal(t, string("12345"), contract.UserID, "Incorrect UserID in contract")
	})

	// Test Case 2: Incorrect Number of Arguments
	t.Run("Incorrect Number of Arguments", func(t *testing.T) {
		response := stub.MockInvoke("2", [][]byte{[]byte("SignPlatformContract"), []byte("12345"), []byte("A7324HAS7234SADF734JSDF"), []byte("ExtraArg")})

		assert.NotEqual(t, shim.OK, response.GetStatus(), "Function unexpectedly succeeded")
	})

	// Test Case 3: Non-existing User
	t.Run("Non-existing User", func(t *testing.T) {
		response := stub.MockInvoke("3", [][]byte{[]byte("SignPlatformContract"), []byte("54321"), []byte("A7324HAS7234SADF734JSDF")})

		assert.NotEqual(t, shim.OK, response.GetStatus(), "Function unexpectedly succeeded")
	})

	// Test Case 4: Invalid User ID
	t.Run("Invalid User ID", func(t *testing.T) {
		response := stub.MockInvoke("4", [][]byte{[]byte("SignPlatformContract"), []byte("InvalidUserID"), []byte("A7324HAS7234SADF734JSDF")})

		assert.NotEqual(t, shim.OK, response.GetStatus(), "Function unexpectedly succeeded")
	})
}

// TestReadPlatformContract tests the ReadPlatformContract function
func TestReadPlatformContract(t *testing.T) {
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	key := "12345"
	id := key

	// Creating a dummy user and platform contract to be used in the tests.
	user := User{ID: id}
	userBytes, _ := json.Marshal(user)
	contract := PlatformContract{
		UserID:             id,
		SignedContractHash: "A7324HAS7234SADF734JSDF",
	}
	contractBytes, _ := json.Marshal(contract)

	// Start a transaction to put user and platform contract into the ledger
	stub.MockTransactionStart("tx1")
	err := stub.PutState(key, userBytes)
	if err != nil {
		t.Fatalf("Failed to put the user into the stub: %s", err.Error())
	}
	contractKey := "PlatformContract_" + id
	err = stub.PutState(contractKey, contractBytes)
	if err != nil {
		t.Fatalf("Failed to put the platform contract into the stub: %s", err.Error())
	}
	stub.MockTransactionEnd("tx1")

	// Test Case: Successfully Read Platform Contract
	t.Run("Successfully Read Platform Contract", func(t *testing.T) {
		response := stub.MockInvoke("tx2", [][]byte{[]byte("ReadPlatformContract"), []byte(id)})

		assert.Equal(t, int32(shim.OK), response.GetStatus(), "Unexpected error during contract read")

		var readContract PlatformContract
		err := json.Unmarshal(response.GetPayload(), &readContract)
		assert.NoError(t, err, "Error unmarshalling read contract")
		assert.Equal(t, id, readContract.UserID, "UserID does not match")
		assert.Equal(t, "A7324HAS7234SADF734JSDF", readContract.SignedContractHash, "Contract hash does not match")
	})
}

func TestTradingContract(t *testing.T) {
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	userID := "12345"
	contractHash := "HASH1234"

	// Start a transaction to put user into state
	stub.MockTransactionStart("txSetUp")
	user := User{ID: userID}
	userBytes, _ := json.Marshal(user)
	err := stub.PutState(userID, userBytes)
	if err != nil {
		t.Fatalf("Failed to put the user into the stub: %s", err.Error())
	}
	stub.MockTransactionEnd("txSetUp")

	contractStatuses := []string{"BidCreated", "BidAccepted", "BidExecuted", "BidTerminated", "BidRejected"}

	for _, status := range contractStatuses {
		t.Run(fmt.Sprintf("Sign and read trading contract with status %s", status), func(t *testing.T) {
			// Sign Trading Contract within a transaction context
			stub.MockTransactionStart(fmt.Sprintf("txSign_%s", status))
			signResp := stub.MockInvoke("1", [][]byte{
				[]byte("SignTradingContract"),
				[]byte(userID),
				[]byte(contractHash),
				[]byte(status),
			})
			stub.MockTransactionEnd(fmt.Sprintf("txSign_%s", status))
			assert.Equal(t, int32(shim.OK), signResp.GetStatus(), "Signing contract failed")

			// Read Trading Contract
			readResp := stub.MockInvoke("1", [][]byte{
				[]byte("ReadTradingContract"),
				[]byte(userID),
				[]byte(status),
			})
			assert.Equal(t, int32(shim.OK), readResp.GetStatus(), "Reading contract failed")

			var readContract TradingContract
			err := json.Unmarshal(readResp.GetPayload(), &readContract)
			assert.NoError(t, err, "Failed to unmarshal contract")
			assert.Equal(t, userID, readContract.UserID, "UserID mismatch")
			assert.Equal(t, contractHash, readContract.SignedContractHash, "Contract hash mismatch")
			assert.Equal(t, status, readContract.BidStatus, "Contract status mismatch")
		})
	}

	// Test Case: Attempt to read non-existing contract
	t.Run("Attempt to read non-existing contract", func(t *testing.T) {
		response := stub.MockInvoke("3", [][]byte{
			[]byte("ReadTradingContract"),
			[]byte(userID),
			[]byte("NonExistingStatus"),
		})
		assert.NotEqual(t, shim.OK, response.GetStatus(), "Expected failure when reading non-existing contract")
	})
}

func TestRegisterOrder(t *testing.T) {
	// Mock stub creation
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	// Test Case 1: Successfully register a new order
	t.Run("Successfully Register a New Order", func(t *testing.T) {
		response := stub.MockInvoke("1", [][]byte{
			[]byte("RegisterOrder"),
			[]byte("1"),          // bidMatchID
			[]byte("BidCreated"), // bidStatus
			[]byte("4"),          // orderID
			[]byte("0"),          // onMarketPrice
			[]byte("200"),        // orderCost
			[]byte("payment5"),   // paymentID
			[]byte("slot1234"),   // slotID
			[]byte("300"),        // totalQuantity
			[]byte("3.5"),        // unitCost
			[]byte("6"),          // userID
			[]byte("50"),         // slotExecDate
			[]byte("Buy"),        // action
		})

		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		orderAsBytes, err := stub.GetState("Order_4")
		assert.NoError(t, err, "Error getting order from ledger")

		var order Order
		err = json.Unmarshal(orderAsBytes, &order)
		assert.NoError(t, err, "Error unmarshalling order")

		assert.Equal(t, string("4"), order.ID, "Order ID mismatch")
		assert.Equal(t, string("1"), order.BidMatchID, "BidMatchID mismatch")
	})

	// Test Case 2: Provide incorrect number of arguments
	t.Run("Incorrect Number of Arguments", func(t *testing.T) {
		response := stub.MockInvoke("2", [][]byte{
			[]byte("RegisterOrder"),
			[]byte("1"),
		})

		assert.Equal(t, int32(shim.ERROR), response.GetStatus(), "Function unexpectedly succeeded")
		assert.Contains(t, response.GetMessage(), "Incorrect number of arguments")
	})

	// Test Case 3: Successfully update an existing order
	t.Run("Successfully Register a New Order", func(t *testing.T) {
		response := stub.MockInvoke("1", [][]byte{
			[]byte("RegisterOrder"),
			[]byte("1"),          // bidMatchID
			[]byte("BidCreated"), // bidStatus
			[]byte("4"),          // orderID
			[]byte("0"),          // onMarketPrice
			[]byte("200"),        // orderCost
			[]byte("payment5"),   // paymentID
			[]byte("slot1236"),   // slotID
			[]byte("300"),        // totalQuantity
			[]byte("3.5"),        // unitCost
			[]byte("6"),          // userID
			[]byte("50"),         // slotExecDate
			[]byte("Buy"),        // action
		})

		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		orderAsBytes, err := stub.GetState("Order_4")
		assert.NoError(t, err, "Error getting order from ledger")

		var order Order
		err = json.Unmarshal(orderAsBytes, &order)
		assert.NoError(t, err, "Error unmarshalling order")

		assert.Equal(t, string("4"), order.ID, "Order ID mismatch")
		assert.Equal(t, string("1"), order.BidMatchID, "BidMatchID mismatch")

		response = stub.MockInvoke("1", [][]byte{
			[]byte("RegisterOrder"),
			[]byte("1"),          // bidMatchID
			[]byte("BidCreated"), // bidStatus
			[]byte("4"),          // orderID
			[]byte("0"),          // onMarketPrice
			[]byte("200"),        // orderCost
			[]byte("payment5"),   // paymentID
			[]byte("slot1235"),   // slotID
			[]byte("300"),        // totalQuantity
			[]byte("3.5"),        // unitCost
			[]byte("6"),          // userID
			[]byte("50"),         // slotExecDate
			[]byte("Buy"),        // action
		})

		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		orderAsBytes, err = stub.GetState("Order_4")
		assert.NoError(t, err, "Error getting order from ledger")

		err = json.Unmarshal(orderAsBytes, &order)
		assert.NoError(t, err, "Error unmarshalling order")

		assert.Equal(t, string("4"), order.ID, "Order ID mismatch")
		assert.Equal(t, string("1"), order.BidMatchID, "BidMatchID mismatch")
	})
}

func TestProcessBidMatch(t *testing.T) {
	// Mock stub creation
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	// Test Case 1: Successfully process a new BidMatch
	t.Run("Successfully Process a New BidMatch", func(t *testing.T) {
		response := stub.MockInvoke("1", [][]byte{
			[]byte("ProcessBidMatch"),
			[]byte("120"),        // bidMatchTms
			[]byte("Slot1"),      // slotID
			[]byte("BidCreated"), // bidStatus
			[]byte("100"),        // bidUnitPrice
			[]byte("Buyer4"),     // buyerUserID
			[]byte("2.5"),        // deliveredBidUnits
			[]byte("BidMatch1"),  // bidMatchID
			[]byte("3.5"),        // originalBidUnits
			[]byte("Seller5"),    // sellerUserID
			[]byte("Buy6"),       // transactionBuyID
			[]byte("Sell7"),      // transactionSellID
		})

		assert.Equal(t, int32(shim.OK), response.GetStatus(), "Unexpected error: "+response.GetMessage())

		bidMatchAsBytes, err := stub.GetState("BidMatch_BidMatch1")
		assert.NoError(t, err, "Error getting BidMatch from ledger")

		var bidMatch BidMatch
		err = json.Unmarshal(bidMatchAsBytes, &bidMatch)
		assert.NoError(t, err, "Error unmarshalling BidMatch")

		assert.Equal(t, string("BidMatch1"), bidMatch.ID, "BidMatch ID mismatch")
		assert.Equal(t, string("Buyer4"), bidMatch.BuyerUserId, "BuyerUserId mismatch")
	})

	// Test Case 2: Provide incorrect number of arguments
	t.Run("Incorrect Number of Arguments", func(t *testing.T) {
		response := stub.MockInvoke("2", [][]byte{
			[]byte("ProcessBidMatch"),
			[]byte("1"),
		})

		assert.Equal(t, int32(shim.ERROR), response.GetStatus(), "Function unexpectedly succeeded")
		assert.Contains(t, response.GetMessage(), "Incorrect number of arguments")
	})
}

func TestProcessEnergyBid(t *testing.T) {
	// Mock stub creation
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	// Test Case 1: Successfully process an EnergyBid
	t.Run("Successfully Process EnergyBid", func(t *testing.T) {
		response := stub.MockInvoke("1", [][]byte{
			[]byte("ProcessEnergyBid"),
			[]byte("EnergyBid1"), // energyBidID
			[]byte("BidMatch1"),  // bidMatchID
			[]byte("10.5"),       // initialBidUnits
			[]byte("8.7"),        // acceptedBidUnits
			[]byte("5.0"),        // buyerMeterUnit
			[]byte("7.2"),        // sellerMeterUnit
			[]byte("3.0"),        // buyerBroughtUnitFromSeller
			[]byte("6.5"),        // sellerSoldUnitToBuyer
			[]byte("4.8"),        // sellerSoldUnitToGrid
			[]byte("9.2"),        // buyerSoldUnitToGrid
			[]byte("2.1"),        // buyerBroughtUnitFromGrid
			[]byte("Reason1"),    // reason
		})

		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		energyBidAsBytes, err := stub.GetState("EnergyBid_EnergyBid1")
		assert.NoError(t, err, "Error getting EnergyBid from ledger")

		var energyBid EnergyBid
		err = json.Unmarshal(energyBidAsBytes, &energyBid)
		assert.NoError(t, err, "Error unmarshalling EnergyBid")

		assert.Equal(t, "EnergyBid1", energyBid.ID, "EnergyBid ID mismatch")
		assert.Equal(t, 10.5, energyBid.InitialBidUnits, "InitialBidUnits mismatch")
		// Add assertions for other fields
	})

	// Test Case 2: Provide incorrect number of arguments
	t.Run("Incorrect Number of Arguments", func(t *testing.T) {
		response := stub.MockInvoke("2", [][]byte{
			[]byte("ProcessEnergyBid"),
			[]byte("1"),
		})

		assert.Equal(t, int32(shim.ERROR), response.GetStatus(), "Function unexpectedly succeeded")
		assert.Contains(t, response.GetMessage(), "Incorrect number of arguments")
	})
}

func TestReadOrder(t *testing.T) {
	// Mock stub creation
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	// Registering a new order
	response := stub.MockInvoke("1", [][]byte{
		[]byte("RegisterOrder"),
		[]byte("1"),          // bidMatchID
		[]byte("BidCreated"), // bidStatus
		[]byte("4"),          // orderID
		[]byte("0"),          // onMarketPrice
		[]byte("200"),        // orderCost
		[]byte("payment5"),   // paymentID
		[]byte("slot1237"),   // slotID
		[]byte("300"),        // totalQuantity
		[]byte("3.5"),        // unitCost
		[]byte("6"),          // userID
		[]byte("50"),         // slotExecDate
		[]byte("Sell"),       // action
	})
	assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

	// Test Case: Successfully read an order
	t.Run("Successfully Read an Order", func(t *testing.T) {
		response := stub.MockInvoke("1", [][]byte{
			[]byte("ReadOrder"),
			[]byte("4"), // orderID
		})

		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		var order Order
		err := json.Unmarshal(response.GetPayload(), &order)
		assert.NoError(t, err, "Error unmarshalling order")

		assert.Equal(t, string("4"), order.ID, "Order ID mismatch")
		assert.Equal(t, string("1"), order.BidMatchID, "BidMatchID mismatch")
	})

	// Test Case: Try to read an order that doesn't exist
	t.Run("Try to Read Nonexistent Order", func(t *testing.T) {
		response := stub.MockInvoke("2", [][]byte{
			[]byte("ReadOrder"),
			[]byte("99"), // orderID
		})

		assert.Equal(t, int32(shim.ERROR), response.GetStatus(), "Function unexpectedly succeeded")
		assert.Contains(t, response.GetMessage(), "not found")
	})
}

func TestReadBidMatch(t *testing.T) {
	// Mock stub creation
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	// Registering a new BidMatch
	response := stub.MockInvoke("1", [][]byte{
		[]byte("ProcessBidMatch"),
		[]byte("120"),        // bidMatchTms
		[]byte("Slot2"),      // slotID
		[]byte("BidCreated"), // bidStatus
		[]byte("100"),        // bidUnitPrice
		[]byte("Buyer4"),     // buyerUserID
		[]byte("2.5"),        // deliveredBidUnits
		[]byte("BidMatch1"),  // bidMatchID
		[]byte("3.5"),        // originalBidUnits
		[]byte("Seller5"),    // sellerUserID
		[]byte("Buy6"),       // transactionBuyID
		[]byte("Sell7"),      // transactionSellID
	})
	assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

	// Test Case 1: Successfully read a BidMatch
	t.Run("Successfully Read a BidMatch", func(t *testing.T) {
		response := stub.MockInvoke("1", [][]byte{
			[]byte("ReadBidMatch"),
			[]byte("BidMatch1"), // bidMatchID
		})

		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		var bidMatch BidMatch
		err := json.Unmarshal(response.GetPayload(), &bidMatch)
		assert.NoError(t, err, "Error unmarshalling BidMatch")

		assert.Equal(t, string("BidMatch1"), bidMatch.ID, "BidMatch ID mismatch")
		assert.Equal(t, string("Buyer4"), bidMatch.BuyerUserId, "BuyerUserId mismatch")
	})

	// Test Case 2: Try to read a BidMatch that doesn't exist
	t.Run("Try to Read Nonexistent BidMatch", func(t *testing.T) {
		response := stub.MockInvoke("2", [][]byte{
			[]byte("ReadBidMatch"),
			[]byte("99"), // bidMatchID
		})

		assert.Equal(t, int32(shim.ERROR), response.GetStatus(), "Function unexpectedly succeeded")
		assert.Contains(t, response.GetMessage(), "not found")
	})
}

func TestReadEnergyBid(t *testing.T) {
	// Mock stub creation
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	// Registering a new BidMatch
	response := stub.MockInvoke("1", [][]byte{
		[]byte("ProcessEnergyBid"),
		[]byte("EnergyBid1"), // energyBidID
		[]byte("BidMatch1"),  // bidMatchID
		[]byte("10.5"),       // initialBidUnits
		[]byte("8.7"),        // acceptedBidUnits
		[]byte("5.0"),        // buyerMeterUnit
		[]byte("7.2"),        // sellerMeterUnit
		[]byte("3.0"),        // buyerBroughtUnitFromSeller
		[]byte("6.5"),        // sellerSoldUnitToBuyer
		[]byte("4.8"),        // sellerSoldUnitToGrid
		[]byte("9.2"),        // buyerSoldUnitToGrid
		[]byte("2.1"),        // buyerBroughtUnitFromGrid
		[]byte("Reason1"),    // reason
	})

	assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

	// Test Case 1: Successfully read a BidMatch
	t.Run("Successfully Read a BidMatch", func(t *testing.T) {
		response := stub.MockInvoke("1", [][]byte{
			[]byte("ReadEnergyBid"),
			[]byte("EnergyBid1"), // bidMatchID
		})

		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		var energyBid EnergyBid
		err := json.Unmarshal(response.GetPayload(), &energyBid)
		assert.NoError(t, err, "Error unmarshalling EnergyBid")

		assert.Equal(t, string("EnergyBid1"), energyBid.ID, "EnergyBid ID mismatch")
		assert.Equal(t, string("BidMatch1"), energyBid.BidMatchID, "BuyerUserId mismatch")
	})

	// Test Case 2: Try to read a EnergyBid that doesn't exist
	t.Run("Try to Read Nonexistent EnergyBid", func(t *testing.T) {
		response := stub.MockInvoke("2", [][]byte{
			[]byte("ReadEnergyBid"),
			[]byte("EnergyBid2"), // bidMatchID
		})

		assert.Equal(t, int32(shim.ERROR), response.GetStatus(), "Function unexpectedly succeeded")
		assert.Contains(t, response.GetMessage(), "not found")
	})
}
