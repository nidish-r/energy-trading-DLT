const express = require('express');
const { SimulationMap, progressHelper } = require('../util/simulationProgress')
const { parameterMap } = require('../util/methodParams')
const router = express.Router();
const { InitializeFabric, contractWrite, getSimulationOutput } = require('../util/fabricHelper')

let contractInstance; 
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
        contractInstance = await InitializeFabric();
    }
    const { phaseArray, steps } = progressHelper(simulationPhase + 1);
    const paramArray = parameterMap[simulationPhase + 1];

    const functionCalls = []
    if(paramArray != undefined) {
        for(let i = 0; i < paramArray.length; i++) {
        await contractWrite(contractInstance, paramArray[i].methodName, paramArray[i].methodParams)
        functionCalls.push(paramArray[i].methodName + ": " + paramArray[i].methodParams[0]);
    }}
    
    const simulationObject = await getSimulationOutput(contractInstance, simulationPhase + 1);
    res.render('overview', { simulationPhase, simulationObject, phaseArray, labels, progress: simulationPhase + 2, functionCalls });
    simulationPhase++;
});

function timeout(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

exports.rouutes = router;
exports.products = products;
