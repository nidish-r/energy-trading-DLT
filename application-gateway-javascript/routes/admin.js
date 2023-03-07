const express = require('express');
const { SimulationMap, progressHelper } = require('../util/simulationProgress')
const { parameterMap } = require('../util/methodParams')
const router = express.Router();
const { InitializeFabric, contractWrite, getSimulationOutput, cleanJSONResponse } = require('../util/fabricHelper')

let contractInstance; 
let gateway;

const products = [];

let simulationPhase = 0;

//  /admin/add-product => GET 
router.get('/', (req, res, next) => {
    const labels = [
        SimulationMap[0],
        SimulationMap[1],
        SimulationMap[2],
        SimulationMap[3],
        SimulationMap[4],
        SimulationMap[5],
        SimulationMap[6],
        SimulationMap[7],
        SimulationMap[8],
        SimulationMap[9],
        SimulationMap[10],
        SimulationMap[11],
        SimulationMap[12]
    ];

    // Render the run-simulation page 
    res.render('overview', {labels, progress: 1});
});

//  /admin/add-product => POST
router.post('/', async (req, res, next) => {
    const labels = [
        SimulationMap[0],
        SimulationMap[1],
        SimulationMap[2],
        SimulationMap[3],
        SimulationMap[4],
        SimulationMap[5],
        SimulationMap[6],
        SimulationMap[7],
        SimulationMap[8],
        SimulationMap[9],
        SimulationMap[10],
        SimulationMap[11],
        SimulationMap[12]
    ];

    /* code block for running each phase on Hyperledger Fabric */
    if(contractInstance == undefined) {
        const response = await InitializeFabric();
        contractInstance = response.contractInstance;
        gateway = response.gateway;
    }

    const phaseArray = progressHelper(simulationPhase + 1);
    const paramArray = parameterMap[simulationPhase + 1];

    let phaseCruncher = 0;
    if (simulationPhase + 1 <= 3) {
        phaseCruncher = 2
    } else {
        phaseCruncher = simulationPhase
    }

    const functionCalls = [];
    const explorerLinks = [];

    if(paramArray != undefined) {
        for(let i = 0; i < paramArray.length; i++) {
        const response = await contractWrite(contractInstance, paramArray[i].methodName, paramArray[i].methodParams);
        const parsendResponse = cleanJSONResponse(new TextDecoder().decode(response));
        functionCalls.push(paramArray[i].methodName + ": " + paramArray[i].methodParams[0]);
        explorerLinks.push(`http://localhost:8080/?tab=transactions&transId=${parsendResponse.result}#/transactions`);
    }}
    
    const simulationObject = await getSimulationOutput(contractInstance, simulationPhase + 1);
    res.render('overview', { simulationPhase, simulationObject, phaseArray, labels, progress: phaseCruncher, functionCalls, explorerLinks });
    simulationPhase++;
});

exports.rouutes = router;
exports.products = products;
