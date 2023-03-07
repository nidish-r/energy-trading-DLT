const SimulationMap = {
  0 : "BEGIN SIMULATION",
  1: "Init Network Components",
  2: "Init Fleet",
  3: "Users Init",
  4: "Users Recharge",
  5: "Allocate Batteries to Fleet",
  6: "Dock Batteries",
  7: "Return From Service",
  8: "Transfer Batteries from SS to Users",
  9: "Transfer Batteries for Users to SS",
  10: "Transfer Batteries from SS to Users",
  11: "Transfer Batteries for Users to SS",
  12: "SIMULATION FINISHED"
}

function progressHelper(simulationPhase) {
  const phaseArray = []
  for(let i = 0; i <= simulationPhase + 1; i++) { 
    if(i < simulationPhase + 1) {
      phaseArray.push(SimulationMap[i]);
    } else {
      if(SimulationMap[i] != undefined ) {
        phaseArray.push("NEXT: " + SimulationMap[i]);
      }
    }
  }

  return phaseArray;
}

module.exports = { SimulationMap, progressHelper }