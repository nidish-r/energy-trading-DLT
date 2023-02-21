/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

const grpc = require('@grpc/grpc-js');
const { connect, Contract, Identity, Signer, signers } = require('@hyperledger/fabric-gateway');
const crypto = require('crypto');
const fs = require('fs');
const path = require('path');
const { TextDecoder } = require('util');
const {
    readSSNetwork1,
    readSSNetwork2,
    readSwappingStation1,
    readSwappingStation2,
    readUser,
    readBattery1,
    readBattery2} = require('./methodParams');

const channelName = envOrDefault('CHANNEL_NAME', 'mychannel');
const chaincodeName = envOrDefault('CHAINCODE_NAME', 'basic');
const mspId = envOrDefault('MSP_ID', 'Org1MSP');

// Path to crypto materials.
const cryptoPath = envOrDefault('CRYPTO_PATH', path.resolve(__dirname, '..', '..', '..', 'test-network', 'organizations', 'peerOrganizations', 'org1.example.com'));

// Path to user private key directory.
const keyDirectoryPath = envOrDefault('KEY_DIRECTORY_PATH', path.resolve(cryptoPath, 'users', 'User1@org1.example.com', 'msp', 'keystore'));

// Path to user certificate.
const certPath = envOrDefault('CERT_PATH', path.resolve(cryptoPath, 'users', 'User1@org1.example.com', 'msp', 'signcerts', 'User1@org1.example.com-cert.pem'));

// Path to peer tls certificate.
const tlsCertPath = envOrDefault('TLS_CERT_PATH', path.resolve(cryptoPath, 'peers', 'peer0.org1.example.com', 'tls', 'ca.crt'));

// Gateway peer endpoint.
const peerEndpoint = envOrDefault('PEER_ENDPOINT', 'localhost:7051');

// Gateway peer SSL host name override.
const peerHostAlias = envOrDefault('PEER_HOST_ALIAS', 'peer0.org1.example.com');

const utf8Decoder = new TextDecoder();

async function InitializeFabric() {

    await displayInputParameters();

    // The gRPC client connection should be shared by all Gateway connections to this endpoint.
    const client = await newGrpcConnection();

    const gateway = connect({
        client,
        identity: await newIdentity(),
        signer: await newSigner(),
        // Default timeouts for different gRPC calls
        evaluateOptions: () => {
            return { deadline: Date.now() + 5000 }; // 5 seconds
        },
        endorseOptions: () => {
            return { deadline: Date.now() + 15000 }; // 15 seconds
        },
        submitOptions: () => {
            return { deadline: Date.now() + 5000 }; // 5 seconds
        },
        commitStatusOptions: () => {
            return { deadline: Date.now() + 60000 }; // 1 minute
        },
    });

        // Get a network instance representing the channel where the smart contract is deployed.
        const network = gateway.getNetwork(channelName);

        // Get the smart contract from the network.
        return network.getContract(chaincodeName);
}

async function newGrpcConnection() {
    const tlsRootCert = await fs.promises.readFile(tlsCertPath);
    const tlsCredentials = grpc.credentials.createSsl(tlsRootCert);
    return new grpc.Client(peerEndpoint, tlsCredentials, {
        'grpc.ssl_target_name_override': peerHostAlias,
    });
}

async function newIdentity() {
    const credentials = await fs.promises.readFile(certPath);
    return { mspId, credentials };
}

async function newSigner() {
    const files = await fs.promises.readdir(keyDirectoryPath);
    const keyPath = path.resolve(keyDirectoryPath, files[0]);
    const privateKeyPem = await fs.promises.readFile(keyPath);
    const privateKey = crypto.createPrivateKey(privateKeyPem);
    return signers.newPrivateKeySigner(privateKey);
}

/**
 * Submit a transaction synchronously, blocking until it has been committed to the ledger.
 */
async function contractWrite(contract, methodType, args) {
    console.log(`\n--> Submit Transaction: ${methodType}, with arguments ${args}`);

    await contract.submitTransaction(
        methodType,
        ...args
    );

    console.log('*** Transaction committed successfully');
}

/**
 * Submit a transaction synchronously, blocking until it has been committed to the ledger.
 */
async function contractRead(contract, methodType, args) {
    console.log(`\n--> Evaluate Transaction: ${methodType}, with arguments ${args}`);

    const resultBytes =  await contract.evaluateTransaction(
        methodType,
        ...args
    );
    
    const resultJson = utf8Decoder.decode(resultBytes);
    console.log('*** Read Transaction Successful');
    return resultJson;
}

/**
 * Submit a transaction synchronously, blocking until it has been committed to the ledger.
 */
 async function getSimulationOutput(contract, simulationPhase) {
    let simulationObject = {
        swappingNetworks: [],
        batteries: [],
        swappingStations: [],
        user: [],
    }

    let readSwappingStation1Response, readSwappingStation2Response, readBattery1Response, readBattery2Response, readUserResponse;

    const readSSNetwork1Response = cleanJSONResponse(await contractRead(contract, readSSNetwork1.methodName, readSSNetwork1.methodParams));
    const readSSNetwork2Response = cleanJSONResponse(await contractRead(contract, readSSNetwork2.methodName, readSSNetwork2.methodParams));
    
    const readSSNetworkArray = [
    [readSSNetwork1Response.result.name_network, 
        readSSNetwork1Response.result.Id_Network,
        readSSNetwork1Response.result.totalBatteries,
        readSSNetwork1Response.result.status,
        readSSNetwork1Response.result.wallet], 
    [readSSNetwork2Response.result.name_network, 
        readSSNetwork2Response.result.Id_Network,
        readSSNetwork2Response.result.totalBatteries,
        readSSNetwork2Response.result.status,
        readSSNetwork2Response.result.wallet]]

    simulationObject.swappingNetworks = readSSNetworkArray;
    
    if(simulationPhase > 1) {
        readSwappingStation1Response = cleanJSONResponse(await contractRead(contract, readSwappingStation1.methodName, readSwappingStation1.methodParams));
        readSwappingStation2Response = cleanJSONResponse(await contractRead(contract, readSwappingStation2.methodName, readSwappingStation2.methodParams));
        
        const readSwappingStationArray = [ 
    [readSwappingStation1Response.result.id_swappingStation, 
        readSwappingStation1Response.result.id_Network,
        readSwappingStation1Response.result.totalBatteries,
        readSwappingStation1Response.result.activeBatteries,
        readSwappingStation1Response.result.dischargedBatteries], 
    [readSwappingStation2Response.result.id_swappingStation, 
        readSwappingStation2Response.result.id_Network,
        readSwappingStation2Response.result.totalBatteries,
        readSwappingStation2Response.result.activeBatteries,
        readSwappingStation2Response.result.dischargedBatteries]]

        simulationObject.swappingStations = readSwappingStationArray;
    }

    if(simulationPhase > 2) {
        readBattery1Response = cleanJSONResponse(await contractRead(contract, readBattery1.methodName, readBattery1.methodParams));
        readBattery2Response = cleanJSONResponse(await contractRead(contract, readBattery2.methodName, readBattery2.methodParams));
        
        const readSwappingStationArray = [ 
    [readBattery1Response.result.id_battery, 
        readBattery1Response.result.id_Network,
        readBattery1Response.result.user,
        readBattery1Response.result.dockedStation,
        readBattery1Response.result.soC,
        readBattery1Response.result.soH,
        readBattery1Response.result.energyContent,
        readBattery1Response.result.escrowedAmount],
    [readBattery2Response.result.id_battery, 
        readBattery2Response.result.id_Network,
        readBattery2Response.result.user,
        readBattery2Response.result.dockedStation,
        readBattery2Response.result.soC,
        readBattery2Response.result.soH,
        readBattery2Response.result.energyContent,
        readBattery2Response.result.escrowedAmount]]
    

        simulationObject.batteries = readSwappingStationArray;
    }

    if(simulationPhase > 3) {
        readUserResponse = cleanJSONResponse(await contractRead(contract, readUser.methodName, readUser.methodParams));
        const readSwappingStationArray = [ 
    [readUserResponse.result.userName, 
        readUserResponse.result.id_user,
        readUserResponse.result.rentedBattery,
        readUserResponse.result.wallet]]

        simulationObject.user = readSwappingStationArray;
    }

    console.log(simulationObject);
    return simulationObject;
}


function cleanJSONResponse(response) {
    const jsonObject = JSON.parse(response.replace('RESULT-->', '').replace('<--RESULT', ''))
    // let returnValue = {result : {}};
    // for (const key in jsonObject.result) {
    //     returnValue.result[key] = !jsonObject.result[key] ? '-' : jsonObject.result[key];
    // }

    // return returnValue;
    return jsonObject;
}

/**
 * envOrDefault() will return the value of an environment variable, or a default value if the variable is undefined.
 */
function envOrDefault(key, defaultValue) {
    return process.env[key] || defaultValue;
}

/**
 * displayInputParameters() will print the global scope parameters used by the main driver routine.
 */
async function displayInputParameters() {
    console.log(`channelName:       ${channelName}`);
    console.log(`chaincodeName:     ${chaincodeName}`);
    console.log(`mspId:             ${mspId}`);
    console.log(`cryptoPath:        ${cryptoPath}`);
    console.log(`keyDirectoryPath:  ${keyDirectoryPath}`);
    console.log(`certPath:          ${certPath}`);
    console.log(`tlsCertPath:       ${tlsCertPath}`);
    console.log(`peerEndpoint:      ${peerEndpoint}`);
    console.log(`peerHostAlias:     ${peerHostAlias}`);
}

module.exports = { InitializeFabric, contractWrite, getSimulationOutput }