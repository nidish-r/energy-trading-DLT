/*
 SPDX-License-Identifier: Apache-2.0
*/

/*
====CHAINCODE EXECUTION SAMPLES (CLI) ==================

==== Invoke batteries ====
peer chaincode invoke -C myc1 -n battery_transfer -c '{"Args":["RecordBattery","battery1","blue","5","tom","35"]}'
peer chaincode invoke -C myc1 -n battery_transfer -c '{"Args":["RecordBattery","battery2","red","4","tom","50"]}'
peer chaincode invoke -C myc1 -n battery_transfer -c '{"Args":["RecordBattery","battery3","blue","6","tom","70"]}'
peer chaincode invoke -C myc1 -n battery_transfer -c '{"Args":["TransferBattery","battery2","jerry"]}'
peer chaincode invoke -C myc1 -n battery_transfer -c '{"Args":["TransferBatteryByColor","blue","jerry"]}'
peer chaincode invoke -C myc1 -n battery_transfer -c '{"Args":["DeleteBattery","battery1"]}'

==== Query batteries ====
peer chaincode query -C myc1 -n battery_transfer -c '{"Args":["ReadBattery","battery1"]}'
peer chaincode query -C myc1 -n battery_transfer -c '{"Args":["GetBatteriesByRange","battery1","battery3"]}'
peer chaincode query -C myc1 -n battery_transfer -c '{"Args":["GetBatteryHistory","battery1"]}'

Rich Query (Only supported if CouchDB is used as state database):
peer chaincode query -C myc1 -n battery_transfer -c '{"Args":["QueryBatteriesByOwner","tom"]}'
peer chaincode query -C myc1 -n battery_transfer -c '{"Args":["QueryBatteries","{\"selector\":{\"owner\":\"tom\"}}"]}'

Rich Query with Pagination (Only supported if CouchDB is used as state database):
peer chaincode query -C myc1 -n battery_transfer -c '{"Args":["QueryBatteriesWithPagination","{\"selector\":{\"owner\":\"tom\"}}","3",""]}'

INDEXES TO SUPPORT COUCHDB RICH QUERIES

Indexes in CouchDB are required in order to make JSON queries efficient and are required for
any JSON query with a sort. Indexes may be packaged alongside
chaincode in a META-INF/statedb/couchdb/indexes directory. Each index must be defined in its own
text file with extension *.json with the index definition formatted in JSON following the
CouchDB index JSON syntax as documented at:
http://docs.couchdb.org/en/2.3.1/api/database/find.html#db-index

This battery transfer ledger example chaincode demonstrates a packaged
index which you can find in META-INF/statedb/couchdb/indexes/indexOwner.json.

If you have access to the your peer's CouchDB state database in a development environment,
you may want to iteratively test various indexes in support of your chaincode queries.  You
can use the CouchDB Fauxton interface or a command line curl utility to create and update
indexes. Then once you finalize an index, include the index definition alongside your
chaincode in the META-INF/statedb/couchdb/indexes directory, for packaging and deployment
to managed environments.

In the examples below you can find index definitions that support battery transfer ledger
chaincode queries, along with the syntax that you can use in development environments
to create the indexes in the CouchDB Fauxton interface or a curl command line utility.


Index for docType, owner.

Example curl command line to define index in the CouchDB channel_chaincode database
curl -i -X POST -H "Content-Type: application/json" -d "{\"index\":{\"fields\":[\"docType\",\"owner\"]},\"name\":\"indexOwner\",\"ddoc\":\"indexOwnerDoc\",\"type\":\"json\"}" http://hostname:port/myc1_batteries/_index


Index for docType, owner, size (descending order).

Example curl command line to define index in the CouchDB channel_chaincode database:
curl -i -X POST -H "Content-Type: application/json" -d "{\"index\":{\"fields\":[{\"size\":\"desc\"},{\"docType\":\"desc\"},{\"owner\":\"desc\"}]},\"ddoc\":\"indexSizeSortDoc\", \"name\":\"indexSizeSortDesc\",\"type\":\"json\"}" http://hostname:port/myc1_batteries/_index

Rich Query with index design doc and index name specified (Only supported if CouchDB is used as state database):
peer chaincode query -C myc1 -n battery_transfer -c '{"Args":["QueryBatteries","{\"selector\":{\"docType\":\"battery\",\"owner\":\"tom\"}, \"use_index\":[\"_design/indexOwnerDoc\", \"indexOwner\"]}"]}'

Rich Query with index design doc specified only (Only supported if CouchDB is used as state database):
peer chaincode query -C myc1 -n battery_transfer -c '{"Args":["QueryBatteries","{\"selector\":{\"docType\":{\"$eq\":\"battery\"},\"owner\":{\"$eq\":\"tom\"},\"size\":{\"$gt\":0}},\"fields\":[\"docType\",\"owner\",\"size\"],\"sort\":[{\"size\":\"desc\"}],\"use_index\":\"_design/indexSizeSortDoc\"}"]}'
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const index = "model~number"

// SimpleChaincode implements the fabric-contract-api-go programming model
type SimpleChaincode struct {
	contractapi.Contract
}

type Battery struct {
	DocType         string  `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Cdc             uint16  `json:"cdc"`
	EnergyContent   float32 `json:"energyContent"`
	Id              string  `json:"id"` //the field tags are needed to keep case from bouncing around
	ManufacturerId  string  `json:"manufacturerId"`
	ManufactureDate string  `json:"manufactureDate"`
	ModelNumber     string  `json:"modelNumber"`
	Owner           string  `json:"owner"`
	SoC             uint8   `json:"soC"`
	SoH             uint8   `json:"soH"`
	Status          string  `json:"status"`
	User            string  `json:"user"`
}

//
// Structure for `SwappingStation`
//
type SwappingStation struct {
	Address             string `json:"address"`
	Company             string `json:"company"`
	ContactNumber       string `json:"contactNumber"`
	EmailId             string `json:"emailId"`
	GeoCoordinates      string `json:"geoCoordinates"`
	Id                  string `json:"id"`
	LicenseNumber       string `json:"licenseNumber"`
	Password            string `json:"password"`
	SwappingStationName string `json:"swappingStationName"`
	Wallet              int64  `json:"wallet"`
	DocType             string `json:"docType"`
}

// HistoryQueryResult structure used for returning result of history query
type HistoryQueryResult struct {
	Record    *Battery  `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

// PaginatedQueryResult structure used for returning paginated query results and metadata
type PaginatedQueryResult struct {
	Records             []*Battery `json:"records"`
	FetchedRecordsCount int32      `json:"fetchedRecordsCount"`
	Bookmark            string     `json:"bookmark"`
}

// RecordBattery initializes a new battery in the ledger
func (t *SimpleChaincode) RecordBattery(ctx contractapi.TransactionContextInterface,
	docType string,
	cdc uint16,
	energyContent float32,
	id string,
	manufacturerId string,
	manufactureDate string,
	modelNumber string,
	owner string,
	soC uint8,
	soH uint8,
	status string,
	user string) error {
	exists, err := t.BatteryExists(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get battery: %v", err)
	}
	if exists {
		return fmt.Errorf("battery already exists: %s", id)
	}

	// Check battery record authorization - this sample assumes Org1 is TruePower with privilege to add new batteries
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSPID: %v", err)
	}
	if clientMSPID != "Org1MSP" {
		return fmt.Errorf("client is not authorized to record new batteries")
	}

	battery := &Battery{
		DocType:         "battery",
		Cdc:             cdc,
		EnergyContent:   energyContent,
		Id:              id, //the field tags are needed to keep case from bouncing around
		ManufacturerId:  manufacturerId,
		ManufactureDate: manufactureDate,
		ModelNumber:     modelNumber,
		Owner:           owner,
		SoC:             soC,
		SoH:             soH,
		Status:          status,
		User:            user,
	}

	batteryBytes, err := json.Marshal(battery)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, batteryBytes)
	if err != nil {
		return err
	}

	//  Create an index to enable mode number-based range queries, e.g. return all batteries with Model Number EV-12V-27AH .
	//  An 'index' is a normal key-value entry in the ledger.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~model~number.
	//  This will enable very efficient state range queries based on composite keys matching indexName~model~*
	colorNameIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{battery.ModelNumber, battery.Id})
	if err != nil {
		return err
	}
	//  Save index entry to world state. Only the key name is needed, no need to store a duplicate copy of the battery.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	return ctx.GetStub().PutState(colorNameIndexKey, value)
}

// ReadBattery retrieves an battery from the ledger
func (t *SimpleChaincode) ReadBattery(ctx contractapi.TransactionContextInterface, batteryID string) (*Battery, error) {
	batteryBytes, err := ctx.GetStub().GetState(batteryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get battery %s: %v", batteryID, err)
	}
	if batteryBytes == nil {
		return nil, fmt.Errorf("battery %s does not exist", batteryID)
	}

	var battery Battery
	err = json.Unmarshal(batteryBytes, &battery)
	if err != nil {
		return nil, err
	}

	return &battery, nil
}

// DeleteBattery removes an battery key-value pair from the ledger
func (t *SimpleChaincode) DeleteBattery(ctx contractapi.TransactionContextInterface, batteryID string) error {
	battery, err := t.ReadBattery(ctx, batteryID)
	if err != nil {
		return err
	}

	err = ctx.GetStub().DelState(batteryID)
	if err != nil {
		return fmt.Errorf("failed to delete battery %s: %v", batteryID, err)
	}

	colorNameIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{battery.ModelNumber, battery.Id})
	if err != nil {
		return err
	}

	// Delete index entry
	return ctx.GetStub().DelState(colorNameIndexKey)
}

// TransferBattery transfers an battery by setting a new owner name on the battery
func (t *SimpleChaincode) TransferBattery(ctx contractapi.TransactionContextInterface, batteryID, newOwner string) error {
	battery, err := t.ReadBattery(ctx, batteryID)
	if err != nil {
		return err
	}

	battery.Owner = newOwner
	batteryBytes, err := json.Marshal(battery)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(batteryID, batteryBytes)
}

// constructQueryResponseFromIterator constructs a slice of batteries from the resultsIterator
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Battery, error) {
	var batteries []*Battery
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var battery Battery
		err = json.Unmarshal(queryResult.Value, &battery)
		if err != nil {
			return nil, err
		}
		batteries = append(batteries, &battery)
	}

	return batteries, nil
}

// GetBatteriesByRange performs a range query based on the start and end keys provided.
// Read-only function results are not typically submitted to ordering. If the read-only
// results are submitted to ordering, or if the query is used in an update transaction
// and submitted to ordering, then the committing peers will re-execute to guarantee that
// result sets are stable between endorsement time and commit time. The transaction is
// invalidated by the committing peers if the result set has changed between endorsement
// time and commit time.
// Therefore, range queries are a safe option for performing update transactions based on query results.
func (t *SimpleChaincode) GetBatteriesByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*Battery, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

// TransferBatteryBy  will transfer batteries of a given color to a certain new owner.
// Uses GetStateByPartialCompositeKey (range query) against color~name 'index'.
// Committing peers will re-execute range queries to guarantee that result sets are stable
// between endorsement time and commit time. The transaction is invalidated by the
// committing peers if the result set has changed between endorsement time and commit time.
// Therefore, range queries are a safe option for performing update transactions based on query results.
// Example: GetStateByPartialCompositeKey/RangeQuery
func (t *SimpleChaincode) TransferBatteryByModelNumber(ctx contractapi.TransactionContextInterface, modelNumber, newOwner string) error {
	// Execute a key range query on all keys starting with 'color'
	BatteryModelResultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(index, []string{modelNumber})
	if err != nil {
		return err
	}
	defer BatteryModelResultsIterator.Close()

	for BatteryModelResultsIterator.HasNext() {
		responseRange, err := BatteryModelResultsIterator.Next()
		if err != nil {
			return err
		}

		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return err
		}

		if len(compositeKeyParts) > 1 {
			returnedBatteryID := compositeKeyParts[1]
			battery, err := t.ReadBattery(ctx, returnedBatteryID)
			if err != nil {
				return err
			}
			battery.Owner = newOwner
			batteryBytes, err := json.Marshal(battery)
			if err != nil {
				return err
			}
			err = ctx.GetStub().PutState(returnedBatteryID, batteryBytes)
			if err != nil {
				return fmt.Errorf("transfer failed for battery %s: %v", returnedBatteryID, err)
			}
		}
	}

	return nil
}

// QueryBatteriesByOwner queries for batteries based on the owners name.
// This is an example of a parameterized query where the query logic is baked into the chaincode,
// and accepting a single query parameter (owner).
// Only available on state databases that support rich query (e.g. CouchDB)
// Example: Parameterized rich query
func (t *SimpleChaincode) QueryBatteriesByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]*Battery, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"battery","owner":"%s"}}`, owner)
	return getQueryResultForQueryString(ctx, queryString)
}

// QueryBatteriesByUser queries for batteries based on the users name.
// This is an example of a parameterized query where the query logic is baked into the chaincode,
// and accepting a single query parameter (user).
// Only available on state databases that support rich query (e.g. CouchDB)
// Example: Parameterized rich query
func (t *SimpleChaincode) QueryBatteriesByUser(ctx contractapi.TransactionContextInterface, user string) ([]*Battery, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"battery","user":"%s"}}`, user)
	return getQueryResultForQueryString(ctx, queryString)
}

// QueryBatteries uses a query string to perform a query for batteries.
// Query string matching state database syntax is passed in and executed as is.
// Supports ad hoc queries that can be defined at runtime by the client.
// If this is not desired, follow the QueryBatteriesForOwner example for parameterized queries.
// Only available on state databases that support rich query (e.g. CouchDB)
// Example: Ad hoc rich query
func (t *SimpleChaincode) QueryBatteries(ctx contractapi.TransactionContextInterface, queryString string) ([]*Battery, error) {
	return getQueryResultForQueryString(ctx, queryString)
}

// getQueryResultForQueryString executes the passed in query string.
// The result set is built and returned as a byte array containing the JSON results.
func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Battery, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

// GetBatteriesByRangeWithPagination performs a range query based on the start and end key,
// page size and a bookmark.
// The number of fetched records will be equal to or lesser than the page size.
// Paginated range queries are only valid for read only transactions.
// Example: Pagination with Range Query
func (t *SimpleChaincode) GetBatteriesByRangeWithPagination(ctx contractapi.TransactionContextInterface, startKey string, endKey string, pageSize int, bookmark string) (*PaginatedQueryResult, error) {

	resultsIterator, responseMetadata, err := ctx.GetStub().GetStateByRangeWithPagination(startKey, endKey, int32(pageSize), bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	batteries, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	return &PaginatedQueryResult{
		Records:             batteries,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

// QueryBatteriesWithPagination uses a query string, page size and a bookmark to perform a query
// for batteries. Query string matching state database syntax is passed in and executed as is.
// The number of fetched records would be equal to or lesser than the specified page size.
// Supports ad hoc queries that can be defined at runtime by the client.
// If this is not desired, follow the QueryBatteriesForOwner example for parameterized queries.
// Only available on state databases that support rich query (e.g. CouchDB)
// Paginated queries are only valid for read only transactions.
// Example: Pagination with Ad hoc Rich Query
func (t *SimpleChaincode) QueryBatteriesWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int, bookmark string) (*PaginatedQueryResult, error) {

	return getQueryResultForQueryStringWithPagination(ctx, queryString, int32(pageSize), bookmark)
}

// getQueryResultForQueryStringWithPagination executes the passed in query string with
// pagination info. The result set is built and returned as a byte array containing the JSON results.
func getQueryResultForQueryStringWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	batteries, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	return &PaginatedQueryResult{
		Records:             batteries,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

// GetBatteryHistory returns the chain of custody for an battery since issuance.
func (t *SimpleChaincode) GetBatteryHistory(ctx contractapi.TransactionContextInterface, batteryID string) ([]HistoryQueryResult, error) {
	log.Printf("GetBatteryHistory: ID %v", batteryID)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(batteryID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var battery Battery
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &battery)
			if err != nil {
				return nil, err
			}
		} else {
			battery = Battery{
				Id: batteryID,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &battery,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}

// BatteryExists returns true when battery with given ID exists in the ledger.
func (t *SimpleChaincode) BatteryExists(ctx contractapi.TransactionContextInterface, batteryID string) (bool, error) {
	batteryBytes, err := ctx.GetStub().GetState(batteryID)
	if err != nil {
		return false, fmt.Errorf("failed to read battery %s from world state. %v", batteryID, err)
	}

	return batteryBytes != nil, nil
}

// // InitLedger creates the initial set of batteries in the ledger.
// func (t *SimpleChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
// 	batteries := []Battery{
// 		{DocType: "battery", ID: "battery1", Color: "blue", Size: 5, Owner: "Tomoko", AppraisedValue: 300},
// 		{DocType: "battery", ID: "battery2", Color: "red", Size: 5, Owner: "Brad", AppraisedValue: 400},
// 		{DocType: "battery", ID: "battery3", Color: "green", Size: 10, Owner: "Jin Soo", AppraisedValue: 500},
// 		{DocType: "battery", ID: "battery4", Color: "yellow", Size: 10, Owner: "Max", AppraisedValue: 600},
// 		{DocType: "battery", ID: "battery5", Color: "black", Size: 15, Owner: "Adriana", AppraisedValue: 700},
// 		{DocType: "battery", ID: "battery6", Color: "white", Size: 15, Owner: "Michel", AppraisedValue: 800},
// 	}

// 	for _, battery := range batteries {
// 		err := t.RecordBattery(ctx, battery.ID, battery.Color, battery.Size, battery.Owner, battery.AppraisedValue)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func main() {
	chaincode, err := contractapi.NewChaincode(&SimpleChaincode{})
	if err != nil {
		log.Panicf("Error creating battery chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting battery chaincode: %v", err)
	}
}
